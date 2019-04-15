package nasync

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type TaskQueue struct {
	taskChan       chan *task // buffered queue used in non-runtime  tasks
	taskChanSize   int64      // queue buffer size
	chanAutoResize bool

	goroutineCount      int  //task worker goroutine number
	goroutineAutoResize bool //协程个数自动扩容
	goroutineContext    context.Context
	closeFunc           context.CancelFunc

	lostMessageCount int //丢失的消息个数

	sendTaskBlock   bool          //task send 是否阻塞
	sendTaskTimeout time.Duration //阻塞send 超时的时间

	syncMutex *sync.Mutex

	logFunc LogFunc
	closed  bool
}

func NewTaskQueue(queueSize int64, goroutineCount int) *TaskQueue {

	if queueSize < 1 {
		panic("queueSize must greater than 0")
	}

	if goroutineCount < 1 {
		panic("goroutineCount must greater than 0")
	}

	taskQueue := &TaskQueue{}

	taskQueue.taskChanSize = queueSize
	taskQueue.taskChan = make(chan *task, queueSize)

	taskQueue.goroutineCount = goroutineCount
	taskQueue.goroutineContext, taskQueue.closeFunc = context.WithCancel(context.Background())

	taskQueue.sendTaskBlock = false //默认非阻塞
	taskQueue.sendTaskTimeout = 0   //默认阻塞情况下不超时

	taskQueue.syncMutex = &sync.Mutex{}

	//启动协程 异步处理task
	taskQueue.initHandler()

	return taskQueue
}

//非阻塞发送
func (tq *TaskQueue) SetSendUnBlock() {
	tq.sendTaskBlock = false
}

//阻塞发送  无限超时
func (tq *TaskQueue) SetSendBlock() {
	tq.sendTaskBlock = true
	tq.sendTaskTimeout = 0
}

//阻塞发送 超时限制
func (tq *TaskQueue) SetSendBlockWithTimeout(timeout time.Duration) {
	if timeout < 0 {
		timeout = 0

	}

	tq.sendTaskBlock = true
	tq.sendTaskTimeout = timeout
}

//send task
func (tq *TaskQueue) Send(handler interface{}, params ...interface{}) (resultStatus bool, resultError *Error) {

	defer func() {
		if err := recover(); err != nil {
			resultStatus = false
			msg := fmt.Sprintf("%s", err)
			resultError = &Error{Code: ERROR_UNKNOWN, Msg: msg}
		}

	}()

	if tq.closed {
		return false, &Error{Code: ERROR_QUEUE_CLOSED, Msg: "queue is closed"}
	}

	task := newTask(handler, params...)

	if tq.sendTaskBlock {
		if tq.sendTaskTimeout > 0 { //超时控制
			select {
			case tq.taskChan <- task:
				tq.log(LEVEL_DEBUG, fmt.Sprintf("send task %v", task))
			case <-time.After(tq.sendTaskTimeout):
				tq.log(LEVEL_INFO, fmt.Sprintf("task abandoned for timeout %d, %v", tq.sendTaskTimeout, task))
				return false, &Error{Code: ERROR_TIMEOUT, Msg: "task abandoned for timeout"}
			}

		} else {
			tq.taskChan <- task //一直阻塞直到成功
			tq.log(LEVEL_DEBUG, fmt.Sprintf("send task %v", task))
		}

	} else { //非阻塞send task
		select {
		case tq.taskChan <- task:
			tq.log(LEVEL_DEBUG, fmt.Sprintf("send task %v", task))
		default:
			tq.log(LEVEL_INFO, fmt.Sprintf("task abandoned for queue full %d, %v", tq.taskChanSize, task))
			return false, &Error{Code: ERROR_QUEUE_FULL, Msg: "task abandoned for queue full"}
		}
	}

	return true, nil
}

//处理task
func (tq *TaskQueue) initHandler() {

	for i := 0; i < tq.goroutineCount; i++ {
		tq.AddHandler()
	}
}

//增加协程
func (tq *TaskQueue) AddHandler() {
	tq.syncMutex.Lock()

	defer tq.syncMutex.Unlock()

	go func() {
		gid := getGID()
		tq.log(LEVEL_INFO, fmt.Sprintf("add handler goroutine %d", gid))

		for {

			select {
			case task := <-tq.taskChan:
				task.Do()
				tq.log(LEVEL_DEBUG, fmt.Sprintf("gid:%d, handle task %v", gid, task))
			case <-tq.goroutineContext.Done():
				tq.log(LEVEL_INFO, fmt.Sprintf("close handler goroutine %d", gid))
				return
			}
		}

	}()
}

//关闭队列
func (tq *TaskQueue) Close() {
	tq.syncMutex.Lock()

	defer tq.syncMutex.Unlock()
	if !tq.closed {
		close(tq.taskChan)

		tq.closeFunc()

		tq.closed = true
	}
}

func (tq *TaskQueue) SetLogFunc(logFunc LogFunc) {

	tq.logFunc = logFunc

}

func (tq *TaskQueue) log(level int, logStr string) {
	if tq.logFunc != nil {
		tq.logFunc(level, logStr)
	}
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

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

	goroutineCount   int //task worker goroutine number
	goroutineContext context.Context
	closeFunc        context.CancelFunc

	sendTaskBlock   bool          //
	sendTaskTimeout time.Duration //block send timeout

	syncMutex *sync.Mutex

	logFunc LogFunc
	closed  bool
}

func NewUnBlockQueue(queueSize int64, handlerNum int) *TaskQueue {

	queue := newTaskQueue(queueSize, handlerNum)

	return queue
}

func NewBlockQueue(queueSize int64, handlerNum int) *TaskQueue {

	queue := newTaskQueue(queueSize, handlerNum)

	queue.setSendBlock()

	return queue
}

func NewBlockTimeoutQueue(queueSize int64, handlerNum int, timeout time.Duration) *TaskQueue {

	queue := newTaskQueue(queueSize, handlerNum)

	queue.setSendBlockWithTimeout(timeout)

	return queue
}

func newTaskQueue(queueSize int64, handlerNum int) *TaskQueue {

	if queueSize < 1 {
		panic("queueSize must greater than 0")
	}

	if handlerNum < 1 {
		panic("handlerNum must greater than 0")
	}

	taskQueue := &TaskQueue{}

	taskQueue.taskChanSize = queueSize
	taskQueue.taskChan = make(chan *task, queueSize)

	taskQueue.goroutineContext, taskQueue.closeFunc = context.WithCancel(context.Background())

	taskQueue.sendTaskBlock = false //default not block
	taskQueue.sendTaskTimeout = 0

	taskQueue.syncMutex = &sync.Mutex{}

	//init handler goroutine
	taskQueue.initHandler(handlerNum)

	return taskQueue
}

//unblock
func (tq *TaskQueue) setSendUnBlock() {
	tq.sendTaskBlock = false
}

//block send until succeed
func (tq *TaskQueue) setSendBlock() {
	tq.sendTaskBlock = true
	tq.sendTaskTimeout = 0
}

//block send until timeout
func (tq *TaskQueue) setSendBlockWithTimeout(timeout time.Duration) {
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
		if tq.sendTaskTimeout > 0 { //timeout
			select {
			case tq.taskChan <- task:
				tq.log(LEVEL_DEBUG, fmt.Sprintf("send task %v", task))
			case <-time.After(tq.sendTaskTimeout):
				tq.log(LEVEL_WARN, fmt.Sprintf("task abandoned for timeout %d", tq.sendTaskTimeout))
				return false, &Error{Code: ERROR_TIMEOUT, Msg: "task abandoned for timeout"}
			}

		} else {
			tq.taskChan <- task //block until success
			tq.log(LEVEL_DEBUG, fmt.Sprintf("send task %v", task))
		}

	} else { //unblock send task
		select {
		case tq.taskChan <- task:
			tq.log(LEVEL_DEBUG, fmt.Sprintf("send task %v", task))
		default:
			tq.log(LEVEL_WARN, fmt.Sprintf("task abandoned for queue full, chanSize:%d", tq.taskChanSize))
			return false, &Error{Code: ERROR_QUEUE_FULL, Msg: "task abandoned for queue full"}
		}
	}

	return true, nil
}

func (tq *TaskQueue) initHandler(handlerNum int) {

	for i := 0; i < handlerNum; i++ {
		tq.AddHandler()
	}
}

//add handler goroutine
func (tq *TaskQueue) AddHandler() bool {
	if tq.closed {
		return false
	}

	tq.syncMutex.Lock()

	defer tq.syncMutex.Unlock()

	go func() {
		gid := getGID()
		tq.log(LEVEL_INFO, fmt.Sprintf("add handler goroutine %d", gid))

		for {

			select {
			case task := <-tq.taskChan:
				if task != nil {
					func() {
						defer func() {
							if err := recover(); err != nil {
								tq.log(LEVEL_WARN, fmt.Sprintf("gid:%d, handle task error %v", gid, err))
							}

						}()
						task.Do()
						tq.log(LEVEL_DEBUG, fmt.Sprintf("gid:%d, handle task %v", gid, task))

					}()
				}

			case <-tq.goroutineContext.Done():
				tq.log(LEVEL_INFO, fmt.Sprintf("close handler goroutine %d", gid))
				return
			}
		}

	}()

	tq.goroutineCount++

	return true
}

//关闭队列
func (tq *TaskQueue) Close() {
	tq.syncMutex.Lock()

	defer func() {
		if err := recover(); err != nil {
			tq.log(LEVEL_WARN, fmt.Sprintf("close error %+v", err))
		}

		tq.closeFunc()

		tq.closed = true

		tq.syncMutex.Unlock()
	}()

	if !tq.closed {
		close(tq.taskChan)
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

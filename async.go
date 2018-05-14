package nasync

import "sync"

const (
	DEFAULT_REQSIZE = 1000
	DEFAULT_BUFSIZE = 1000
)

var DefaultAsync *Async

func Do(handler interface{}, params ...interface{}) {
	if DefaultAsync == nil {
		DefaultAsync = New(DEFAULT_REQSIZE, DEFAULT_BUFSIZE)
	}
	DefaultAsync.Do(handler, params...)
}

type Async struct {
	quit     chan bool  // quit signal for the watcher to quit
	taskChan chan *task // queue used in non-runtime  tasks
	bufSize  int
	wait     *sync.WaitGroup
}

func New(ReqSize int, BufSzie int) *Async {
	as := Async{
		quit:     make(chan bool),
		taskChan: make(chan *task, ReqSize),
		bufSize:  BufSzie,
		wait:     &sync.WaitGroup{},
	}

	go as.watcher()
	return &as
}

// Destroy sends quit signal to watcher and releases all the resources.
func (this *Async) Do(handler interface{}, params ...interface{}) {
	t := newTask(handler, params...)
	this.taskChan <- t
}

// Destroy sends quit signal to watcher and releases all the resources.
// Wait for all tasks complete to close
func (this *Async) Close() {
	this.quit <- true
	// wait for watcher quit
	<-this.quit
}

package nasync

import (
	"time"
)

// watcher watches the async.queue channel, and writes the logs to output
func (a *Async) watcher() {
	var buf buffer
	for {
		if a.waiting {
			time.Sleep(time.Second)
			if len(a.taskChan) == 0 {
				a.done <- true
				return
			}
		}
		timeout := time.After(time.Second / 10)
		for i := 0; i < a.bufSize; i++ {
			select {
			case req := <-a.taskChan:
				a.flushReq(&buf, req)
			case <-timeout:
				i = a.bufSize
			case <-a.quit:
				// If quit signal received, cleans the channel
				for {
					select {
					case req := <-a.taskChan:
						a.flushReq(&buf, req)
					default:
						a.flushBuf(&buf)
						a.quit <- true
						return
					}
				}
			}
		}
		a.flushBuf(&buf)
	}
}

// flushReq handles the request and writes the result to writer
func (a *Async) flushReq(b *buffer, t *task) {
	//do print for this
	b.Append(t)
}

// flushBuf flushes the content of buffer to out and reset the buffer
func (a *Async) flushBuf(b *buffer) {
	tasks := b.Tasks()
	if len(tasks) > 0 {
		for _, t := range tasks {
			a.wait.Add(1)
			go func(t *task) {
				t.Do()
				a.wait.Done()
			}(t)
		}
		a.wait.Wait()
		b.Reset()
	}
}

package nasync

import (
	"time"
)

// watcher watches the async.queue channel, and writes the logs to output
func (this *Async) watcher() {
	var buf buffer
	for {
		timeout := time.After(time.Second / 10)
		for i := 0; i < this.bufSize; i++ {
			select {
			case req := <-this.taskChan:
				this.flushReq(&buf, req)
			case <-timeout:
				i = this.bufSize
			case <-this.quit:
			// If quit signal received, cleans the channel
				for {
					select {
					case req := <-this.taskChan:
						this.flushReq(&buf, req)
					default:
						this.flushBuf(&buf)
						this.quit <- true
						return
					}
				}
			}
		}
		this.flushBuf(&buf)
	}
}


// flushReq handles the request and writes the result to writer
func (this *Async) flushReq(b *buffer, t *task) {
	//do print for this
	b.Append(t)
}


// flushBuf flushes the content of buffer to out and reset the buffer
func (this *Async) flushBuf(b *buffer) {
	tasks := b.Tasks()
	if len(tasks) > 0 {
		for _, t := range tasks {
			this.wait.Add(1)
			go func(t *task) {
				t.Do()
				this.wait.Done()
			}(t)
		}
		this.wait.Wait()
		b.Reset()
	}
}



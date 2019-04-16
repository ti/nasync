package nasync

import (
	"fmt"
	"testing"
	"time"
)

func queueLog(level int, logStr string) {
	if level > LEVEL_INFO {
		//fmt.Println(logStr)
	}
}

var handler = func(msg string) string {
	time.Sleep(time.Second * 2)
	fmt.Println(msg)
	panic("xxxxxxxxxxxxxx")
	return "pong"
}

func TestNewUnBlockQueue(t *testing.T) {
	chanSize := 2
	goroutineCount := 2
	taskQueue := NewUnBlockQueue(int64(chanSize), goroutineCount)
	taskQueue.SetLogFunc(queueLog)

	logPrefix := "TestNewUnBlockQueue"

	fmt.Printf("%s queue size chanSize:%d, goroutineCount:%d\n", logPrefix, chanSize, goroutineCount)

	for i := 0; i <= 1; i++ {
		min := i*10 + 1
		max := min + 6
		for j := min; j < max; j++ {

			_, err := taskQueue.Send(handler, fmt.Sprintf("handle, task TestNewUnBlockQueue :%d", j))

			if err != nil {
				fmt.Printf("%s send task fail: %v\n", logPrefix, err)
			} else {
				fmt.Printf("%s send task success: %d\n", logPrefix, j)
			}
		}

		time.Sleep(time.Second * 3)
	}

}

func TestNewBlockQueue(t *testing.T) {
	chanSize := 2
	goroutineCount := 2
	taskQueue := NewBlockQueue(int64(chanSize), goroutineCount)
	taskQueue.SetLogFunc(queueLog)

	logPrefix := "TestNewBlockQueue"

	fmt.Printf("%s queue size chanSize:%d, goroutineCount:%d\n", logPrefix, chanSize, goroutineCount)

	for i := 0; i <= 0; i++ {
		min := i*10 + 1
		max := min + 6
		for j := min; j < max; j++ {

			_, err := taskQueue.Send(handler, fmt.Sprintf("handle, task NewBlockQueue :%d", j))

			if err != nil {
				fmt.Printf("%s send task fail: %v\n", logPrefix, err)
			} else {
				fmt.Printf("%s send task success: %d\n", logPrefix, j)
			}
		}

		time.Sleep(time.Second * 3)
	}
}

func TestNewBlockTimeoutQueue(t *testing.T) {
	chanSize := 2
	goroutineCount := 2
	taskQueue := NewBlockTimeoutQueue(int64(chanSize), goroutineCount, time.Second*1)
	taskQueue.SetLogFunc(queueLog)

	logPrefix := "NewBlockTimeoutQueue"

	fmt.Printf("%s queue size chanSize:%d, goroutineCount:%d\n", logPrefix, chanSize, goroutineCount)

	for i := 0; i <= 0; i++ {
		min := i*10 + 1
		max := min + 6
		for j := min; j < max; j++ {

			_, err := taskQueue.Send(handler, fmt.Sprintf("handle, task NewBlockTimeoutQueue :%d", j))

			if err != nil {
				fmt.Printf("%s send task fail: %v\n", logPrefix, err)
			} else {
				fmt.Printf("%s send task success: %d\n", logPrefix, j)
			}
		}

		time.Sleep(time.Second * 3)
	}
}

func TestClose(t *testing.T) {
	chanSize := 2
	goroutineCount := 2
	taskQueue := NewBlockQueue(int64(chanSize), goroutineCount)
	taskQueue.SetLogFunc(queueLog)

	logPrefix := "TestClose"

	fmt.Printf("%s queue size chanSize:%d, goroutineCount:%d\n", logPrefix, chanSize, goroutineCount)

	go func() {
		time.Sleep(time.Second * 3)
		taskQueue.Close()
	}()

	for i := 0; i <= 0; i++ {
		min := i*10 + 1
		max := min + 10
		for j := min; j < max; j++ {

			_, err := taskQueue.Send(handler, fmt.Sprintf("handle, task TestClose :%d", j))

			if err != nil {
				fmt.Printf("%s send task fail: %v\n", logPrefix, err)
			} else {
				fmt.Printf("%s send task success: %d\n", logPrefix, j)
			}
		}

		time.Sleep(time.Second * 3)
	}
}

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cclehui/nasync"
)

func queueLog(level int, logStr string) {
	if level > nasync.LEVEL_DEBUG {
		fmt.Println(logStr)
	}
}

func main() {
	handler := func(msg string) string {
		time.Sleep(time.Second * 2)
		fmt.Println(msg)
		return "pong"
	}

	taskQueue := nasync.NewTaskQueue(3, 3)
	taskQueue.SetLogFunc(queueLog)

	//taskQueue.SetSendBlock()
	taskQueue.SetSendBlockWithTimeout(time.Second * 1)

	go func() {
		time.Sleep(time.Second * 2)

		taskQueue.Close()

	}()

	for j := 1; j <= 15; j++ {

		_, err := taskQueue.Send(handler, fmt.Sprintf("ttttttttt, send j:%d", j))

		fmt.Printf("send response %v\n", err)

	}
	//taskQueue.Close()

	time.Sleep(time.Second * 5000)

	for j := 1; j < 10; j++ {

		taskQueue.Send(handler, fmt.Sprintf("xxxxxxxxxxx, send j:%d", j))

	}

	time.Sleep(time.Second * 500000)

	os.Exit(0)

	//do async max 1000 tasks in max 10 go goroutine
	//as := New(1000, 100)
	as := nasync.New(2, 2)
	defer as.Close()

	i := 1

	for {
		for i = 1; i <= 10; i++ {

			as.Do(handler, fmt.Sprintf("recieve, %d", i))

			fmt.Printf("send task , %s, %d\n", time.Now().Format("2006-01-02 15:04:05"), i)

		}
		time.Sleep(time.Second * 10)
	}

}

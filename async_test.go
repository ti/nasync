package nasync

import (
	"fmt"
	"testing"
)

func TestAsyncTask(t *testing.T) {
	Do(ping)
	//time.Sleep(5 * time.Second)
}

func TestAsyncAdvanced(t *testing.T) {
	//do async max 1000 tasks in max 10 go goroutine
	as := New(1000, 100)
	defer as.Close()

	handler := func(msg string) string {
		fmt.Print(msg)
		return "pong"
	}
	as.Do(handler, "ping")
	//time.Sleep(5 * time.Second)
}

func ping() string {
	fmt.Println("ping")
	return "pong"
}

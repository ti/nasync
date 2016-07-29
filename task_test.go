package nasync

import (
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	handler := func(msg string) string {
		if msg == "ping" {
			return "pong"
		}
		return "error";
	}
	tsk := newTask(handler,"ping")
	tsk.Do()
	time.Sleep(2*time.Second)
}


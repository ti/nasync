package nasync

import (
	"testing"
	"fmt"
	"time"

)

func TestAsync(t *testing.T) {

	//do async max 1000 tasks in max 10 go goroutine
	as := New(1000,100)

	handler := func(msg string) string {
		fmt.Print(msg)
		time.Sleep(time.Second/2)
		return "pong";
	}


	for i :=0;i<1000;i++ {
		as.Do(handler,"ping " + fmt.Sprint(i))
	}
	fmt.Println("------------doing the task in async---------------")

	time.Sleep(10*time.Second)

	as.Close()
}

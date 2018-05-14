package nasync

import (
	"reflect"
)

// model function into task
type task struct {
	handler reflect.Value
	params  []reflect.Value
}

func newTask(handler interface{}, params ...interface{}) *task {

	handlerValue := reflect.ValueOf(handler)

	if handlerValue.Kind() == reflect.Func {
		task := task{
			handler: handlerValue,
			params:  make([]reflect.Value, 0),
		}
		if paramNum := len(params); paramNum > 0 {
			task.params = make([]reflect.Value, paramNum)
			for index, v := range params {
				task.params[index] = reflect.ValueOf(v)
			}
		}
		return &task
	}
	panic("handler not func")
}

//Do single task functions do
func (t *task) Do() []reflect.Value {
	return t.handler.Call(t.params)
}

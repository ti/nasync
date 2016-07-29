package nasync

import (
	"reflect"
)


//一个任务
type task struct {
	handler reflect.Value
	params []reflect.Value
}

//新建一个任务
func newTask(handler interface{},params ...interface{}) *task {

	handlerValue := reflect.ValueOf(handler);

	if(handlerValue.Kind() == reflect.Func){
		task := task{
			handler : handlerValue ,
			params : make([]reflect.Value,0),
		}
		if paramNum := len(params);paramNum > 0{
			task.params = make([]reflect.Value,paramNum);
			for index, v := range params {
				task.params[index] = reflect.ValueOf(v);
			}
		}
		return &task;
	}
	panic("handler not func");
}

//启动异步任务
func (this *task) Do() []reflect.Value {
	return this.handler.Call(this.params)
}

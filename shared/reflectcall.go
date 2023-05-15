package shared

import (
	"fmt"
	"reflect"
)

func Call(method reflect.Method, args []reflect.Value) (rets []reflect.Value, err error) {
	RecoverInfo(fmt.Errorf("methodName:%s|args:%v", method.Name, args))
	r := method.Func.Call(args)
	return r, err
}

func SafeCall(f func()) {
	if f == nil {
		return
	}
	Recover()
	f()
}

package shared

import (
	"fmt"
	"reflect"
)

func Call(method reflect.Method, args []reflect.Value) (rets []reflect.Value, err error) {
	defer RecoverInfo(fmt.Errorf("methodName:%s|args:%v", method.Name, args))
	r := method.Func.Call(args)
	return r, err
}

func SafeCall(f func()) {
	if f == nil {
		return
	}
	defer Recover()
	f()
}

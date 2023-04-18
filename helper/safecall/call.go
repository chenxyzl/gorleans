package safecall

import (
	"errors"
	"fmt"
	"github.com/chenxyzl/gorleans/logger"
	"reflect"
	"runtime/debug"
	"strconv"
)

func Call(method reflect.Method, args []reflect.Value) (rets []reflect.Value, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			// Try to use logger from context here to help trace error cause
			stackTrace := debug.Stack()
			stackTraceAsRawStringLiteral := strconv.Quote(string(stackTrace))
			logger.Errorf("methodName=%s panicData=%v stackTrace=%s", method.Name, rec, stackTraceAsRawStringLiteral)

			if s, ok := rec.(string); ok {
				err = errors.New(s)
			} else {
				err = fmt.Errorf("rpc call internal error - %s: %v", method.Name, rec)
			}
		}
	}()

	r := method.Func.Call(args)
	return r, err
}

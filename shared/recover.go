package shared

import (
	"github.com/chenxyzl/gorleans/glog"
	"runtime/debug"
	"strconv"
)

func Recover() {
	err := recover()
	if err != nil {
		stackTrace := debug.Stack()
		stackTraceAsRawStringLiteral := strconv.Quote(string(stackTrace))
		glog.Errorf("err:%v|stackTrace:%v", err, stackTraceAsRawStringLiteral)
	}
}

func RecoverInfo(info error) {
	if info == nil {
		Recover()
	} else {
		err := recover()
		if err != nil {
			stackTrace := debug.Stack()
			stackTraceAsRawStringLiteral := strconv.Quote(string(stackTrace))
			glog.Errorf("%v|err:%v|stackTrace:%v", info, err, stackTraceAsRawStringLiteral)
		}
	}
}

func RecoverFunc(info error, pc func(err any)) {
	if pc == nil {
		RecoverInfo(info)
	} else {
		err := recover()
		if err != nil {
			stackTrace := debug.Stack()
			stackTraceAsRawStringLiteral := strconv.Quote(string(stackTrace))
			glog.Errorf("%v|err:%v|stackTrace:%v", info, err, stackTraceAsRawStringLiteral)
			pc(err)
		}
	}
}

package shared

import (
	"github.com/chenxyzl/gorleans/logger"
	"runtime/debug"
	"strconv"
)

func Recover() {
	err := recover()
	if err != nil {
		stackTrace := debug.Stack()
		stackTraceAsRawStringLiteral := strconv.Quote(string(stackTrace))
		logger.Errorf("err:%v|stackTrace:%v", err, stackTraceAsRawStringLiteral)
	}
}
func RecoverFunc(pc func(err any)) {
	if pc == nil {
		Recover()
	} else {
		err := recover()
		if err != nil {
			stackTrace := debug.Stack()
			stackTraceAsRawStringLiteral := strconv.Quote(string(stackTrace))
			logger.Errorf("err:%v|stackTrace:%v", err, stackTraceAsRawStringLiteral)
			pc(err)
		}
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
			logger.Errorf("%v|err:%v|stackTrace:%v", info, err, stackTraceAsRawStringLiteral)
		}
	}
}

package shared

import (
	"github.com/chenxyzl/gorleans/glog"
	"go.uber.org/zap"
	"runtime/debug"
	"strconv"
)

func Recover(logger ...*zap.SugaredLogger) {
	err := recover()
	if err != nil {
		stackTrace := debug.Stack()
		stackTraceAsRawStringLiteral := strconv.Quote(string(stackTrace))
		if len(logger) > 0 && logger[0] != nil {
			logger[0].Errorf("err:%v|stackTrace:%v", err, stackTraceAsRawStringLiteral)
		} else {
			glog.Errorf("err:%v|stackTrace:%v", err, stackTraceAsRawStringLiteral)
		}
	}
}

func RecoverInfo(info error, logger ...*zap.SugaredLogger) {
	if info == nil {
		Recover()
	} else {
		err := recover()
		if err != nil {
			stackTrace := debug.Stack()
			stackTraceAsRawStringLiteral := strconv.Quote(string(stackTrace))
			if len(logger) > 0 && logger[0] != nil {
				logger[0].Errorf("%v|err:%v|stackTrace:%v", info, err, stackTraceAsRawStringLiteral)
			} else {
				glog.Errorf("%v|err:%v|stackTrace:%v", info, err, stackTraceAsRawStringLiteral)
			}
		}
	}
}

func RecoverFunc(info error, pc func(err any), logger ...*zap.SugaredLogger) {
	if pc == nil {
		RecoverInfo(info)
	} else {
		err := recover()
		if err != nil {
			stackTrace := debug.Stack()
			stackTraceAsRawStringLiteral := strconv.Quote(string(stackTrace))
			if len(logger) > 0 && logger[0] != nil {
				logger[0].Errorf("%v|err:%v|stackTrace:%v", info, err, stackTraceAsRawStringLiteral)
			} else {
				glog.Errorf("%v|err:%v|stackTrace:%v", info, err, stackTraceAsRawStringLiteral)
			}
			pc(err)
		}
	}
}

package shared

import (
	"github.com/chenxyzl/gorleans/logger"
	"runtime"
	"strings"
)

func Recover(pc func()) {
	if err := recover(); err != nil {
		logger.Errorf("frame Err: %s", err)
		for i := 0; i < 10; i++ {
			pc, file, line, ok := runtime.Caller(i)
			if ok {
				p := strings.Index(file, "/src/")
				if p != -1 {
					file = file[p+len("/src/"):]
				}
				logger.Errorf("frame %d:[func:%s,file:%s,line:%d]",
					i, runtime.FuncForPC(pc).Name(), file, line)
			} else {
				break
			}
		}
		if pc != nil {
			pc()
		}
	}
}

func SafeCall(f func()) {
	if f == nil {
		return
	}
	Recover(nil)
	f()
}

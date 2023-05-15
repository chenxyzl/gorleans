package system

import (
	"github.com/chenxyzl/gorleans/logger"
	"os"
	"os/signal"
	"syscall"
)

func waitStopSignal() {
	// signal.Notify的ch信道是阻塞的(signal.Notify不会阻塞发送信号), 需要设置缓冲
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			callFuncSlice(s.reloadFs)
		default:
			//go time.AfterFunc(conf.failFastTimeout, func() {
			//	// log.Warn("app exit now by force...")
			//	// os.Exit(1)
			//	logger.Errorf("app exit now by force...")
			//})
			//fmt.Println("app exit now...")
			logger.Infof("app exit now...")
			return
		}
	}
}

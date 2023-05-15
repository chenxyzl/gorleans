package system

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/chenxyzl/gorleans/logger"
	"github.com/chenxyzl/gorleans/shared"
	"os"
	"os/exec"
	"path/filepath"
)

type state int

const (
	stateInit    state = 0
	stateRunning state = 1
	stateStop    state = 2
)

var s *S

type S struct {
	system   *actor.ActorSystem
	cluster  *cluster.Cluster
	schedule *scheduler.TimerScheduler

	status state

	reloadFs []func()

	execDir  string
	execFile string
	pid      string
}

func init() {
	s = &S{}
	s.status = stateInit
	arg0, err := exec.LookPath(os.Args[0])
	if err != nil {
		logger.Panic(err)
	}
	absExecFile, err := filepath.Abs(arg0)
	if err != nil {
		logger.Panic(err)
	}
	s.execDir, s.execFile = filepath.Split(absExecFile)
}

func callFuncSlice(fs []func()) {
	defer shared.Recover()
	for _, f := range fs {
		if f != nil {
			f()
		}
	}
}

func (a *S) createPid() {
	// create pid file
	s.pid = s.execDir + s.execFile + ".pid"
	err := os.WriteFile(s.pid, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	if err != nil {
		logger.Panic(err)
	}
	logger.Infof("create pid, pid:%v", s.pid)
}

func (a *S) removePid() {
	if a.pid != "" {
		err := os.Remove(a.pid)
		if err != nil {
			logger.Error(err)
		}
		logger.Infof("remove pid, pid:%v", s.pid)
	}
}

package system

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/chenxyzl/gorleans/logger"
)

func CreateLocalActor(producer actor.Producer, opts ...actor.PropsOption) *actor.PID {
	if s.status != stateRunning {
		logger.Panicf("status error, status:%v", s.status)
	}

	props := actor.PropsFromProducer(producer, opts...)

	pid := s.system.Root.Spawn(props)
	return pid
}

func RootCtx() *actor.RootContext {
	if s.status != stateRunning {
		logger.Panicf("status error, status:%v", s.status)
	}
	return s.system.Root
}

func Cluster() *cluster.Cluster {
	if s.status != stateRunning {
		logger.Panicf("status error, status:%v", s.status)
	}
	return s.cluster
}

func GetSchedule() *scheduler.TimerScheduler {
	if s.status != stateRunning {
		logger.Panicf("status error, status:%v", s.status)
	}
	return s.schedule
}

func RegisterReloadFunc(f func()) {
	s.reloadFs = append(s.reloadFs, f)
}

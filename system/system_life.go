package system

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/etcd"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/chenxyzl/gorleans/logger"
	"github.com/robfig/cron"
	_ "net/http/pprof"
	"time"
)

func Init() {
	s.createPid()
}
func Clean() {
	s.remotePid()
}

func Run(clusterName string, remoteUrl string, etcdBaseKey string, etcdUrl []string, options ...cluster.ConfigOption) {
	if s.status != stateInit {
		logger.Panicf("status error, status:%v", s.status)
	}
	s.status = stateRunning
	s.system = actor.NewActorSystem()
	//provider, err := etcd.NewWithConfig(etcdBaseKey, clientv3.Config{
	//	Endpoints:   etcdUrl,
	//	DialTimeout: time.Second * 5,
	//})
	//if err != nil {
	//	logger.Panic(err)
	//}
	provider, _ := etcd.New()
	config := remote.Configure("localhost", 0)
	lookup := disthash.New()
	clusterConfig := cluster.Configure(clusterName, provider, lookup, config, options...)
	s.cluster = cluster.New(s.system, clusterConfig)

	s.cluster.StartMember()
}

func Tick(f func(timestamp int64)) *cron.Cron {
	tick := func() {
		now := time.Now().Unix()
		f(now)
	}
	cron2 := cron.New() //创建一个cron实例
	//执行定时任务（每5秒执行一次）
	err := cron2.AddFunc("*/1 * * * * *", tick)
	if err != nil {
		panic(err)
	}
	//启动/关闭
	cron2.Start()
	return cron2
}

func WaitStop(beforeQuitFunc func()) {
	if s.status != stateRunning {
		logger.Panicf("status error, status:%v", s.status)
	}
	//等待退出
	waitStopSignal()
	if beforeQuitFunc != nil {
		beforeQuitFunc()
	}
	s.cluster.Shutdown(true)
	s.status = stateStop
}

package system

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/etcd"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/chenxyzl/gorleans/glog"
	"github.com/robfig/cron"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	_ "net/http/pprof"
	"time"
)

func Init() {
	s.createPid()
}
func Clean() {
	s.removePid()
}

func Run(clusterName string, remoteUrl string, etcdBaseKey string, etcdUrl []string, options ...cluster.ConfigOption) {
	if s.status != stateInit {
		glog.Panicf("status error, status:%v", s.status)
	}
	s.status = stateRunning
	s.system = actor.NewActorSystem()
	s.schedule = scheduler.NewTimerScheduler(s.system.Root)
	provider, err := etcd.NewWithConfig(etcdBaseKey, clientv3.Config{
		Endpoints:   etcdUrl,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		glog.Panic(err)
	}
	//provider, _ := etcd.New()
	config := remote.Configure(remoteUrl,
		0,
		remote.WithDialOptions(grpc.WithKeepaliveParams(keepalive.ClientParameters{PermitWithoutStream: true}),
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
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

func WaitStop(needRunning bool, beforeQuitFunc func()) {
	if needRunning && s.status != stateRunning {
		glog.Panicf("status error, status:%v", s.status)
	}
	//等待退出
	waitStopSignal()
	if beforeQuitFunc != nil {
		beforeQuitFunc()
	}
	s.cluster.Shutdown(true)
	s.status = stateStop
}

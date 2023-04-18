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
	"os"
	"os/signal"
	"syscall"
	"time"
)

var s *S

type S struct {
	system  *actor.ActorSystem
	cluster *cluster.Cluster
}

func Run(clusterName string, remoteUrl string, etcdBaseKey string, etcdUrl []string, options ...cluster.ConfigOption) {
	if s != nil {
		logger.Panicf("repeat start")
	}
	s = &S{}
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
		// reload()
		default:
			//go time.AfterFunc(conf.failFastTimeout, func() {
			//	// log.Warn("app exit now by force...")
			//	// os.Exit(1)
			//	logger.Errorf("app exit now by force...")
			//})
			// fmt.Println("app exit now...")
			logger.Infof("app exit now...")
			return
		}
	}
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
	//等待退出
	waitStopSignal()
	if beforeQuitFunc != nil {
		beforeQuitFunc()
	}
	s.cluster.Shutdown(true)
}

func CreateLocalActor(producer actor.Producer, opts ...actor.PropsOption) *actor.PID {
	if s == nil {
		logger.Panicf("must call system.Run first")
	}

	props := actor.PropsFromProducer(producer, opts...)

	pid := s.system.Root.Spawn(props)
	return pid
}

func RootCtx() *actor.RootContext {
	if s == nil {
		logger.Panicf("must call system.Run first")
	}
	return s.system.Root
}

func Cluster() *cluster.Cluster {
	if s == nil || s.cluster == nil {
		logger.Panicf("s.cluster must call system.Run first")
	}
	return s.cluster
}

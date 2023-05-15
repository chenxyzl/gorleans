package grain

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/chenxyzl/gorleans/logger"
	pb "github.com/chenxyzl/gorleans/proto"
	"github.com/chenxyzl/gorleans/shared"
	"github.com/chenxyzl/gorleans/system"
	"time"
)

type GrainActor struct {
	ctx           cluster.GrainContext
	xGrainFactory func(INext) IGrainActor
	inner         IGrainActor
	Timeout       time.Duration
	nextFunc      []func(context actor.Context)
}

func NewClusterKind(factory func(INext) IGrainActor, kindStr string, timeout time.Duration, opts ...actor.PropsOption) *cluster.Kind {
	props := actor.PropsFromProducer(func() actor.Actor {
		return &GrainActor{Timeout: timeout, xGrainFactory: factory, nextFunc: make([]func(context actor.Context), 0)}
	}, opts...)
	kind := cluster.NewKind(kindStr, props)
	return kind
}

// Receive ensures the lifecycle of the actor for the received message
func (a *GrainActor) Receive(ctx actor.Context) {
	//
	shared.RecoverInfo(fmt.Errorf("msg:%v", ctx.Message()))
	//
	switch msg := ctx.Message().(type) {
	case *actor.Started: // pass
	case *cluster.ClusterInit:
		logger.Debugf("GrainActor started:%v ", ctx.Self())
		a.ctx = cluster.NewGrainContext(ctx, msg.Identity, msg.Cluster)
		a.inner = a.xGrainFactory(a)
		a.inner.Init(a.ctx)

		if a.Timeout > 0 {
			ctx.SetReceiveTimeout(a.Timeout)
		}
	case *actor.ReceiveTimeout:
		ctx.Poison(ctx.Self())
	case *actor.PoisonPill:
		ctx.Stop(ctx.Self())
	case *actor.Stopped:
		logger.Debugf("GrainActor stopped:%v ", ctx.Self())
		a.inner.Terminate(a.ctx)
	case actor.AutoReceiveMessage: // pass
	case actor.SystemMessage: // pass

	case *cluster.GrainRequest:
		a.inner.ReceiveDefault(a.ctx)
	case *pb.NextStep:
		a.handleNextStep(ctx)
	default:
		a.inner.ReceiveDefault(a.ctx)
	}
}

func (a *GrainActor) Next(f func(ctx actor.Context)) {
	a.nextFunc = append(a.nextFunc, f)
	system.GetSchedule().SendOnce(0, a.ctx.Self(), &pb.NextStep{})
}
func (a *GrainActor) handleNextStep(ctx actor.Context) {
	var list []func(context actor.Context)
	if len(a.nextFunc) > 0 {
		list = a.nextFunc[:]
		a.nextFunc = nil
	}
	for _, f := range list {
		f(ctx)
	}
}

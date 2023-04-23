package shared

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/chenxyzl/gorleans/logger"
	pb "github.com/chenxyzl/gorleans/proto"
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
	switch msg := ctx.Message().(type) {
	case *actor.Started: // pass
	case *cluster.ClusterInit:
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
		a.inner.Terminate(a.ctx)
	case actor.AutoReceiveMessage: // pass
	case actor.SystemMessage: // pass

	case *cluster.GrainRequest:
		a.inner.ReceiveDefault(a.ctx)
	case *pb.NextStep:
		logger.Debugf("actor:next", ctx.Self())
		a.handleNextStep(ctx)
	default:
		a.inner.ReceiveDefault(a.ctx)
	}
}

func (state *GrainActor) Next(f func(ctx actor.Context)) {
	state.nextFunc = append(state.nextFunc, f)
	system.GetSchedule().SendOnce(0, state.ctx.Self(), &pb.NextStep{})
}
func (state *GrainActor) handleNextStep(ctx actor.Context) {
	var list []func(context actor.Context)
	if len(state.nextFunc) > 0 {
		list = state.nextFunc[:]
		state.nextFunc = nil
	}
	for _, f := range list {
		f(ctx)
	}
}

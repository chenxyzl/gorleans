package shared

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/gorleans/logger"
	pb "github.com/chenxyzl/gorleans/proto"
	"github.com/chenxyzl/gorleans/system"
)

type LocalActor struct {
	ctx           actor.Context
	xLocalFactory func(INext) ILocalActor
	inner         ILocalActor
	nextFunc      []func(context actor.Context)
}

func NewLocalKind(factory func(INext) ILocalActor, opts ...actor.PropsOption) *actor.PID {
	props := actor.PropsFromProducer(func() actor.Actor {
		a := &LocalActor{xLocalFactory: factory, nextFunc: make([]func(context actor.Context), 0)}
		inner := factory(a)
		a.inner = inner
		return a
	}, opts...)
	return system.RootCtx().Spawn(props)
}

func (a *LocalActor) Receive(ctx actor.Context) {
	a.ctx = ctx
	switch ctx.Message().(type) {
	case *actor.Started:
		logger.Debugf("actor:%v started", ctx.Self())
		a.inner.Init(ctx)
	case *actor.Stopping:
		logger.Debugf("actor:%v stopping", ctx.Self())
		a.inner.Terminate(ctx)
	case *actor.Stopped:
		logger.Debugf("actor:%v stopped", ctx.Self())
	case *actor.Restarting:
		logger.Debugf("actor:%v restarting", ctx.Self())
	case *pb.NextStep:
		logger.Debugf("actor:next", ctx.Self())
		a.handleNextStep(ctx)
	default:
		a.inner.ReceiveDefault(ctx)
	}
}

func (a *LocalActor) Next(f func(ctx actor.Context)) {
	a.nextFunc = append(a.nextFunc, f)
	system.GetSchedule().SendOnce(0, a.ctx.Self(), &pb.NextStep{})
}
func (a *LocalActor) handleNextStep(ctx actor.Context) {
	var list []func(context actor.Context)
	if len(a.nextFunc) > 0 {
		list = a.nextFunc[:]
		a.nextFunc = nil
	}
	for _, f := range list {
		f(ctx)
	}
}

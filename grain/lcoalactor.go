package shared

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/gorleans/logger"
)

type LocalActor struct {
	inner ILocalActor
}

func NewLocalActor(inner ILocalActor) *LocalActor {
	return &LocalActor{inner: inner}
}

func (state *LocalActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		logger.Debugf("actor:%v started", context.Self())
		state.inner.Init(context)
	case *actor.Stopping:
		logger.Debugf("actor:%v stopping", context.Self())
		state.inner.Terminate(context)
	case *actor.Stopped:
		logger.Debugf("actor:%v stopped", context.Self())
	case *actor.Restarting:
		logger.Debugf("actor:%v restarting", context.Self())
	default:
		state.inner.ReceiveDefault(context)
	}
}

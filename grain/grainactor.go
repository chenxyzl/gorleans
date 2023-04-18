package shared

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"time"
)

type GrainActor struct {
	ctx     cluster.GrainContext
	inner   IGrainActor
	Timeout time.Duration
}

// Receive ensures the lifecycle of the actor for the received message
func (a *GrainActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started: // pass
	case *cluster.ClusterInit:
		a.ctx = cluster.NewGrainContext(ctx, msg.Identity, msg.Cluster)
		a.inner = xGrainFactory()
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
	default:
		a.inner.ReceiveDefault(a.ctx)
	}
}

func NewClusterKind(factory func() IGrainActor, kindStr string, timeout time.Duration, opts ...actor.PropsOption) *cluster.Kind {
	xGrainFactory = factory
	props := actor.PropsFromProducer(func() actor.Actor {
		return &GrainActor{
			Timeout: timeout,
		}
	}, opts...)
	kind := cluster.NewKind(kindStr, props)
	return kind
}

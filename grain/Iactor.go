package shared

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
)

var xGrainFactory func() IGrainActor

type IGrainActor interface {
	Init(ctx cluster.GrainContext)
	Terminate(ctx cluster.GrainContext)
	ReceiveDefault(ctx cluster.GrainContext)
}

type ILocalActor interface {
	Init(ctx actor.Context)
	Terminate(ctx actor.Context)
	ReceiveDefault(ctx actor.Context)
}

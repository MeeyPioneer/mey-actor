package persistence

import (
	"log"
	"reflect"

	"github.com/meeypioneer/mey-actor/actor"
)

func Using(provider Provider) func(next actor.ActorFunc) actor.ActorFunc {
	return func(next actor.ActorFunc) actor.ActorFunc {
		fn := func(ctx actor.Context) {
			switch ctx.Message().(type) {
			case *actor.Started:
				next(ctx)
				if p, ok := ctx.Actor().(persistent); ok {
					p.init(provider, ctx)
				} else {
					log.Fatalf("Actor type %v is not persistent", reflect.TypeOf(ctx.Actor()))
				}
			default:
				next(ctx)
			}
		}
		return fn
	}
}
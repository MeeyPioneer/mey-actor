package router_test

import (
	"strconv"
	"sync"
	"testing"
	"time"

	actor "github.com/meeypioneer/mey-actor/actor"
	"github.com/meeypioneer/mey-actor/router"
)

type myMessage struct {
	i   int
	pid *actor.PID
}
type getRoutees struct {
	pid *actor.PID
}

func (m *myMessage) Hash() string {
	return strconv.Itoa(m.i)
}

var wait sync.WaitGroup

type routerActor struct{}
type tellerActor struct{}
type managerActor struct {
	set  []*actor.PID
	rpid *actor.PID
}

func (state *routerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *myMessage:
		//log.Printf("%v got message %d", context.Self(), msg.i)
		msg.i++
		wait.Done()
	}
}
func (state *tellerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *myMessage:
		for i := 0; i < 100; i++ {
			msg.pid.Tell(msg)
			time.Sleep(10 * time.Millisecond)
		}

	}
}

func (state *managerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *router.Routees:
		state.set = msg.PIDs
		for i, v := range state.set {
			if i%2 == 0 {
				state.rpid.Tell(&router.RemoveRoutee{PID: v})
				//log.Println(v)

			} else {
				props := actor.FromProducer(func() actor.Actor { return &routerActor{} })
				pid := actor.Spawn(props)
				state.rpid.Tell(&router.AddRoutee{PID: pid})
				//log.Println(v)
			}
		}
		context.Self().Tell(&getRoutees{state.rpid})
	case *getRoutees:
		state.rpid = msg.pid
		msg.pid.Request(&router.GetRoutees{}, context.Self())
	}
}

func TestConcurrency(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	wait.Add(100 * 10000)
	rpid := actor.Spawn(router.NewConsistentHashPool(100).WithProducer(func() actor.Actor { return &routerActor{} }))

	props := actor.FromProducer(func() actor.Actor { return &tellerActor{} })
	for i := 0; i < 10000; i++ {
		pid := actor.Spawn(props)
		pid.Tell(&myMessage{i, rpid})
	}

	props = actor.FromProducer(func() actor.Actor { return &managerActor{} })
	pid := actor.Spawn(props)
	pid.Tell(&getRoutees{rpid})
	wait.Wait()
}

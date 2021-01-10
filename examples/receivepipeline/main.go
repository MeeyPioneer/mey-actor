package main

import (
	"fmt"

	"github.com/AsynkronIT/goconsole"
	"github.com/meeypioneer/mey-actor/actor"
	"github.com/meeypioneer/mey-actor/actor/middleware"
)

type hello struct{ Who string }

func receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		fmt.Printf("Hello %v\n", msg.Who)
	}
}

func main() {
	props := actor.FromFunc(receive).WithMiddleware(middleware.Logger)
	pid := actor.Spawn(props)
	pid.Tell(&hello{Who: "Roger"})
	console.ReadLine()
}

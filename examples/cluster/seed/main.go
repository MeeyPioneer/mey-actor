package main

import (
	"log"

	console "github.com/AsynkronIT/goconsole"
	"github.com/meeypioneer/mey-actor/actor"
	"github.com/meeypioneer/mey-actor/cluster"
	"github.com/meeypioneer/mey-actor/cluster/consul"
	"github.com/meeypioneer/mey-actor/examples/cluster/shared"
	"github.com/meeypioneer/mey-actor/remote"
)

func main() {

	//this node knows about Hello kind
	remote.Register("Hello", actor.FromProducer(func() actor.Actor {
		return &shared.HelloActor{}
	}))

	cp, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}
	cluster.Start("mycluster", "127.0.0.1:8080", cp)

	hello := shared.GetHelloGrain("MyGrain")

	res, err := hello.SayHello(&shared.HelloRequest{Name: "Roger"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Message from grain: %v", res.Message)
	console.ReadLine()

	cluster.Shutdown(true)
}

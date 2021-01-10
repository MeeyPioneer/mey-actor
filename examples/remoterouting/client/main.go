package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/meeypioneer/mey-actor/actor"
	"github.com/meeypioneer/mey-actor/examples/remoterouting/messages"
	"github.com/meeypioneer/mey-actor/mailbox"
	"github.com/meeypioneer/mey-actor/remote"
	"github.com/meeypioneer/mey-actor/router"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GC()

	remote.Start("127.0.0.1:8100")

	p1 := actor.NewPID("127.0.0.1:8101", "remote")
	p2 := actor.NewPID("127.0.0.1:8102", "remote")
	remotePID := actor.Spawn(router.NewConsistentHashGroup(p1, p2))

	messageCount := 1000000

	var wgStop sync.WaitGroup

	props := actor.
		FromProducer(newLocalActor(&wgStop, messageCount)).
		WithMailbox(mailbox.Bounded(10000))

	pid := actor.Spawn(props)

	log.Println("Starting to send")

	t := time.Now()

	for i := 0; i < messageCount; i++ {
		message := &messages.Ping{User: fmt.Sprintf("User_%d", i)}
		remotePID.Request(message, pid)
	}

	wgStop.Wait()

	pid.Stop()

	fmt.Printf("elapsed: %v\n", time.Since(t))

	console.ReadLine()
}

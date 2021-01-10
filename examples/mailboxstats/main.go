package main

import (
	"log"

	console "github.com/AsynkronIT/goconsole"
	"github.com/meeypioneer/mey-actor/actor"
	"github.com/meeypioneer/mey-actor/mailbox"
)

type mailboxLogger struct{}

func (m *mailboxLogger) MailboxStarted() {
	log.Print("Mailbox started")
}
func (m *mailboxLogger) MessagePosted(msg interface{}) {
	log.Printf("Message posted %v", msg)
}
func (m *mailboxLogger) MessageReceived(msg interface{}) {
	log.Printf("Message received %v", msg)
}
func (m *mailboxLogger) MailboxEmpty() {
	log.Print("No more messages")
}

func main() {
	props := actor.FromFunc(func(ctx actor.Context) {

	}).WithMailbox(mailbox.Unbounded(&mailboxLogger{}))
	pid := actor.Spawn(props)
	pid.Tell("Hello")
	console.ReadLine()
}

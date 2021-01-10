package actor

import "github.com/meeypioneer/mey-actor/mailbox"

var (
	defaultDispatcher = mailbox.NewDefaultDispatcher(300)
)

var defaultMailboxProducer = mailbox.Unbounded()

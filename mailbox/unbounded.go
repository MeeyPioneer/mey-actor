package mailbox

import (
	"github.com/meeypioneer/mey-actor/internal/queue/goring"
	"github.com/meeypioneer/mey-actor/internal/queue/mpsc"
)

type unboundedMailboxQueue struct {
	userMailbox *goring.Queue
}

func (q *unboundedMailboxQueue) Push(m interface{}) {
	q.userMailbox.Push(m)
}

func (q *unboundedMailboxQueue) Pop() interface{} {
	m, o := q.userMailbox.Pop()
	if o {
		return m
	}
	return nil
}

// Unbounded returns a producer which creates an unbounded mailbox
func Unbounded(mailboxStats ...Statistics) Producer {
	return func(invoker MessageInvoker, dispatcher Dispatcher) Inbound {
		q := &unboundedMailboxQueue{
			userMailbox: goring.New(10),
		}
		return &defaultMailbox{
			systemMailbox: mpsc.New(),
			userMailbox:   q,
			invoker:       invoker,
			mailboxStats:  mailboxStats,
			dispatcher:    dispatcher,
		}
	}
}

package remote

import "github.com/meeypioneer/mey-actor/actor"

func remoteHandler(pid *actor.PID) (actor.Process, bool) {
	ref := newProcess(pid)
	return ref, true
}

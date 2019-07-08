// Package ferry provides a Ferry synchronization. This acts similarly to a
// sync.WaitGroup in that many gorotuines can Wait for a single Done, but it
// is less efficient and less picky about how it is used.
package ferry

import (
	"sync"
)

type empty struct{}

var e = empty{}

// A Ferry is a synchronization mechanism similar to a WaitGroup.
//
// All calls to Wait are blocking, and become unblocked the next time Done is
// called. A single call to Done will unblock all goroutines Waiting on the
// Ferry.
//
// Calling Done will unblock all currently Waiting goroutines.
//
// If a Wait call occurs while a Ferry is still handling a Done call, then it
// will wait for the Done to complete before it starts listening for a new
// Done call.
//
// In other words: Done always performs a finite amount of work, and will
// always eventually return, even with Wait being called frequently.
type Ferry struct {
	mtx sync.Mutex
	chs []chan empty
}

// Wait causes the goroutine to block until the next Done call to the ferry
func (b *Ferry) Wait() {
	b.mtx.Lock()
	ch := make(chan empty)
	b.chs = append(b.chs, ch)
	b.mtx.Unlock()
	<-ch
}

// Done causes all Waiting goroutines to be unblocked
//
// It is possible to call Done with no goroutines Waiting on the Ferry. This
// is supported behaviour, and is not a runtime error.
func (b *Ferry) Done() {
	b.mtx.Lock()
	sl := b.chs
	b.chs = nil
	b.mtx.Unlock()

	for _, ch := range sl {
		ch <- empty{}
	}
}

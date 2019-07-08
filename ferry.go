package ferry

import (
	"sync"
)

// A Ferry is a synchronization mechanism
// * All calls to Lock are blocking
// * Calling Unlock will unblock all currently waiting calls to Lock
// 
// If there are very many calls to Lock() and the ferry is spending a
// significant amount of time in the Unlock() call, then subsequent
// Lock() calls will wait for the current Unlock() to finish before listening
// for an Unlock().
// In other words: Unlock() always performs a finite amount of work, and will
// always eventually return.
type Ferry interface {
	// Lock causes the goroutine to block until an Unlock call to the ferry
	Lock()

	// Unlock causes all waiting goroutines to be unblocked
	// It is expected that some calls to Unlock() will occur with no Locks
	// acquired. This is expected behaviour, and is not a runtime error.
	Unlock()
}

// NewFerry returns a new Ferry
func NewFerry() Ferry {
	return &block{
		chs: make([]chan empty, 0),
	}
}

type empty struct{}

var e = empty{}

type block struct {
	mtx sync.Mutex
	chs []chan empty
}

func (b *block) Lock() {
	b.mtx.Lock()
	ch := make(chan empty)
	b.chs = append(b.chs, ch)
	b.mtx.Unlock()
	<-ch
}

func (b *block) Unlock() {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	for _, ch := range b.chs {
		ch <- empty{}
	}

	b.chs = make([]chan empty, 0)
}

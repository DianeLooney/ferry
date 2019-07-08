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
	Unlock()
}

// NewFerry returns a new Ferry
func NewFerry() Ferry {
	return &block{
		ch: make(chan empty),
	}
}

type empty struct{}

var e = empty{}

type block struct {
	wg sync.WaitGroup
	ch  chan empty
}

func (b *block) Lock() {
	b.wg.Wait()

	<-b.ch
}

func (b *block) Unlock() {
	b.wg.Wait()

	b.wg.Add(1)
	defer b.wg.Done()

	for {
		select {
		case b.ch <- empty{}:
		default:
			return
		}
	}
}

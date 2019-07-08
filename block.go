package ferry

import (
	"sync"
)

// A Ferry is a synchronization mechanism
// Conceptually, there is a Ferry that loads up many Waits, and then
// a single Go causes all Waits to unblock.
type Ferry interface {
	// Wait blocks the goroutine
	Wait()

	// Go causes all waiting goroutines to be unblocked
	Go()
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
	mtx sync.Mutex
	ch  chan empty
}

func (b *block) unlock() { b.mtx.Unlock() }
func (b *block) lock()   { b.mtx.Lock() }

func (b *block) Wait() {
	b.lock()
	b.unlock()

	<-b.ch
}

func (b *block) Go() {
	b.lock()
	defer b.unlock()

	for {
		select {
		case b.ch <- empty{}:
			continue
		default:
			return
		}
	}
}

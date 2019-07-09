// Package ferry provides a Ferry synchronization. This acts similarly to a
// sync.WaitGroup in that many gorotuines can Wait for a single Done, but it
// is less efficient and less picky about how it is used.
package ferry

type empty struct{}

var e = empty{}

// A Ferry is a synchronization mechanism similar to a WaitGroup.
//
// Calling Wait blocks the current goroutine.
//
// Calling Done will unblocks all goroutines currently Waiting on the Ferry.
//
// A Ferry should only ever be used from ferry.New
//
// If a Wait call occurs while a Ferry is still handling a Done call, then it
// will wait for the Done to complete before it starts listening for a new
// Done call.
//
// In other words: Done always performs a finite amount of work, and will
// always eventually return, even with Wait being called frequently.
type Ferry struct {
	ch chan empty
}

func New() Ferry {
	return Ferry{
		ch: make(chan empty),
	}
}

// Wait causes the goroutine to block until the next Done call to the ferry
func (b *Ferry) Wait() {
	<-b.ch
}

// Done causes all Waiting goroutines to be unblocked
//
// It is possible to call Done with no goroutines Waiting on the Ferry. This
// is supported behaviour, and is not a runtime error.
func (b *Ferry) Done() {
	ch := b.ch
	b.ch = make(chan empty)

	close(ch)
}

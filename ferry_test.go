package ferry_test

import (
	"sync"
	"testing"
	"time"

	. "github.com/dianelooney/ferry"
)

func TestItBlocks(t *testing.T) {
	var order []int
	b := NewFerry()
	go func() {
		time.Sleep(250 * time.Millisecond)
		order = append(order, 1)
		b.Go()
	}()
	b.Wait()
	order = append(order, 2)
	if order[0] != 1 || order[1] != 2 {
		t.Error("It did not block in the proper order")
	}
}
func TestItBlocksMany(t *testing.T) {
	const n = 1000
	done := 0
	mtx := sync.Mutex{}

	b := NewFerry()
	for i := 0; i < n; i++ {
		go func(i int) {
			b.Wait()

			mtx.Lock()
			defer mtx.Unlock()
			done++
		}(i)
	}
	time.Sleep(250 * time.Millisecond)
	if done != 0 {
		t.Errorf("It did not block all goroutines")
	}
	b.Go()
	time.Sleep(250 * time.Millisecond)
	if done != n {
		t.Errorf("It did not unblock all goroutines")
	}
}

func TestItBlocksWhileUnblocking(t *testing.T) {
	const n = 10000
	fail := false
	b := NewFerry()
	for i := 0; i < n; i++ {
		go b.Wait()
	}
	go func() {
		time.Sleep(1 * time.Millisecond)
		b.Wait()
		fail = true
	}()
	b.Go()
	time.Sleep(250 * time.Millisecond)
	if fail {
		t.Errorf("It did not block after an immediate unblock call")
	}
}

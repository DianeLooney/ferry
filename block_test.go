package ferry_test

import (
	"sync"
	"testing"
	"time"

	. "github.com/dianelooney/block"
)

func TestItBlocks(t *testing.T) {
	var order []int
	b := NewBlocker()
	go func() {
		time.Sleep(250 * time.Millisecond)
		order = append(order, 1)
		b.Unblock()
	}()
	b.Block()
	order = append(order, 2)
	if order[0] != 1 || order[1] != 2 {
		t.Error("It did not block in the proper order")
	}
}
func TestItBlocksMany(t *testing.T) {
	const n = 1000
	done := 0
	mtx := sync.Mutex{}

	b := NewBlocker()
	for i := 0; i < n; i++ {
		go func(i int) {
			b.Block()

			mtx.Lock()
			defer mtx.Unlock()
			done++
		}(i)
	}
	time.Sleep(250 * time.Millisecond)
	if done != 0 {
		t.Errorf("It did not block all goroutines")
	}
	b.UnblockSingle()
	time.Sleep(250 * time.Millisecond)
	if done != 1 {
		t.Errorf("It did not unblock all goroutines")
	}
	b.Unblock()
	time.Sleep(250 * time.Millisecond)
	if done != n {
		t.Errorf("It did not unblock all goroutines")
	}
}

func TestItBlocksWhileUnblocking(t *testing.T) {
	const n = 10000
	fail := false
	b := NewBlocker()
	for i := 0; i < n; i++ {
		go func(i int) {
			b.Block()
		}(i)
	}
	b.Unblock()
	go func() {
		b.Block()
		fail = true
	}()
	time.Sleep(250 * time.Millisecond)
	if fail {
		t.Errorf("It did not block after an immediate unblock call")
	}
}

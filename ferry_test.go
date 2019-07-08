package ferry_test

import (
	"sync/atomic"
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
		b.Unlock()
	}()
	b.Lock()
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
			b.Lock()

			mtx.Lock()
			defer mtx.Unlock()
			done++
		}(i)
	}
	time.Sleep(250 * time.Millisecond)
	if done != 0 {
		t.Errorf("It did not block all goroutines")
	}
	b.Unlock()
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
		go func() {
			for {
				b.Lock()
				b.Lock()
				fail = true
			}
		}()
	}

	done := false
	go func() {
		b.Unlock()
		done = true
	}()
	time.Sleep(250 * time.Millisecond)
	
	if !done {
		t.Errorf("It did not return")
	}
	if fail {
		t.Errorf("It unlocked the same goroutine more than once")
	}
}

func TestMultipleUnlockCalls(t *testing.T) {
	const n = 10000
	b := NewFerry()
	var count int32
	for i := 0; i < n; i++ {
		go func() {
				b.Lock()
				b.Lock()
				atomic.AddInt32(&count, 1)
		}()
	}

	time.Sleep(250 * time.Millisecond)
	b.Unlock()
	time.Sleep(250 * time.Millisecond)
	b.Unlock()
	time.Sleep(250 * time.Millisecond)
	
	if n != atomic.LoadInt32(&count) {
		t.Errorf("It did not unblock all goroutines twice, expected '%v' got '%v'", n, count)
	}
}
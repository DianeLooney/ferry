package ferry_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/dianelooney/ferry"
)

func TestItBlocks(t *testing.T) {
	var order []int
	b := New()
	go func() {
		time.Sleep(250 * time.Millisecond)
		order = append(order, 1)
		b.Done()
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

	b := New()
	for i := 0; i < n; i++ {
		go func() {
			b.Wait()

			mtx.Lock()
			defer mtx.Unlock()
			done++
		}()
	}
	time.Sleep(250 * time.Millisecond)
	if done != 0 {
		t.Errorf("It did not block all goroutines")
	}
	b.Done()
	time.Sleep(250 * time.Millisecond)

	mtx.Lock()
	defer mtx.Unlock()

	if done != n {
		t.Errorf("It did not unblock all goroutines")
	}
}

func TestItBlocksWhileUnblocking(t *testing.T) {
	const n = 1000
	fail := false
	b := New()
	for i := 0; i < n; i++ {
		go func() {
			for {
				b.Wait()
				b.Wait()
				fail = true
			}
		}()
	}

	done := false
	mtx := sync.Mutex{}
	go func() {
		b.Done()
		mtx.Lock()
		defer mtx.Unlock()
		done = true
	}()
	time.Sleep(250 * time.Millisecond)

	mtx.Lock()
	defer mtx.Unlock()

	if !done {
		t.Errorf("It did not return")
	}

	if fail {
		t.Errorf("It unlocked the same goroutine more than once")
	}
}

func TestMultipleUnlockCalls(t *testing.T) {
	const n = 1000
	b := New()
	var count int32
	for i := 0; i < n; i++ {
		go func() {
			b.Wait()
			b.Wait()
			atomic.AddInt32(&count, 1)
		}()
	}

	time.Sleep(250 * time.Millisecond)
	b.Done()
	time.Sleep(250 * time.Millisecond)
	b.Done()
	time.Sleep(250 * time.Millisecond)

	if n != atomic.LoadInt32(&count) {
		t.Errorf("It did not unblock all goroutines twice, expected '%v' got '%v'", n, count)
	}
}

func BenchmarkDone(b *testing.B) {
	f := New()
	for i := 0; i < b.N; i++ {
		f.Done()
	}
}

func benchmarkDone_NWaiters(b *testing.B, n int) {
	arr := make([]Ferry, b.N)
	for i := 0; i < b.N; i++ {
		f := New()
		arr[i] = f
		for j := 0; j < n; j++ {
			go f.Wait()
		}
	}
	for i := 0; i < b.N; i++ {
		arr[i].Done()
	}
}
func BenchmarkDone_0Waiters(b *testing.B) {
	benchmarkDone_NWaiters(b, 0)
}
func BenchmarkDone_1Waiter(b *testing.B) {
	benchmarkDone_NWaiters(b, 1)
}
func BenchmarkDone_2Waiter(b *testing.B) {
	benchmarkDone_NWaiters(b, 2)
}
func BenchmarkDone_5Waiter(b *testing.B) {
	benchmarkDone_NWaiters(b, 5)
}
func BenchmarkDone_10Waiter(b *testing.B) {
	benchmarkDone_NWaiters(b, 10)
}
func BenchmarkDone_20Waiter(b *testing.B) {
	benchmarkDone_NWaiters(b, 20)
}
func BenchmarkDone_50Waiter(b *testing.B) {
	benchmarkDone_NWaiters(b, 50)
}
func BenchmarkDone_100Waiter(b *testing.B) {
	benchmarkDone_NWaiters(b, 100)
}

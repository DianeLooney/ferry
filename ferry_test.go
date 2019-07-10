package ferry_test

import (
	"sync"
	"testing"
	"time"

	. "github.com/dianelooney/ferry"
)

func TestItBlocks(t *testing.T) {
	b := New()
	go func() {
		b.Wait()
		t.Error("It did not block")
	}()
	time.Sleep(10 * time.Millisecond)
}

func TestItBlocksMany(t *testing.T) {
	const n = 1000
	b := New()
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			b.Wait()
			wg.Done()
		}()
	}
	time.Sleep(10 * time.Millisecond)
	b.Done()
	wg.Wait()
}

func TestItUnblocksGoroutinesOnce(t *testing.T) {
	const n = 1000
	b := New()
	for i := 0; i < n; i++ {
		go func() {
			b.Wait()
			b.Wait()
			t.Errorf("It unlocked the same goroutine more than once")
		}()
	}
	b.Done()
	time.Sleep(10 * time.Millisecond)
}

func BenchmarkDone(b *testing.B) {
	f := New()
	for i := 0; i < b.N; i++ {
		f.Done()
	}
}

func benchmarkDone_NWaiters(b *testing.B, n int) {
	arr := make([]*Ferry, b.N)
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

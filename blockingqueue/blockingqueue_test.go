package blockingqueue

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBlockingQueue_Bounded(t *testing.T) {
	bq := New[int](2)

	if bq.Len() != 0 {
		t.Errorf("Expected len 0, got %d", bq.Len())
	}

	bq.Push(1)
	bq.Push(2)

	if !bq.TryPush(3) {
		// Expected
	} else {
		t.Error("Expected TryPush to fail on full queue")
	}

	go func() {
		time.Sleep(10 * time.Millisecond)
		bq.Pop()
	}()

	// Should block until the pop happens
	bq.Push(3)

	time.Sleep(5 * time.Millisecond)
	if bq.Len() != 2 {
		t.Errorf("Expected len 2, got %d", bq.Len())
	}
}

func TestContext(t *testing.T) {
	t.Run("PushCtxCancels", func(t *testing.T) {
		bq := New[int](0)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		err := bq.PushCtx(ctx, 1)
		if err == nil {
			t.Error("Expected error on context cancel")
		}
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})

	t.Run("PopCtxCancels", func(t *testing.T) {
		bq := New[int](0)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		_, err := bq.PopCtx(ctx)
		if err == nil {
			t.Error("Expected error on context cancel")
		}
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})
}

func TestBlockingQueue_Concurrency(t *testing.T) {
	bq := New[int](10) // Small bound to force contention
	const count = 1000
	const goroutines = 10

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// Producers
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				bq.Push(j)
			}
		}()
	}

	// Consumers
	var totalPopped int
	var mu sync.Mutex

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				_ = bq.Pop()
				mu.Lock()
				totalPopped++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if totalPopped != count*goroutines {
		t.Errorf("Expected %d items, got %d", count*goroutines, totalPopped)
	}
	if bq.Len() != 0 {
		t.Errorf("Expected empty queue, got len %d", bq.Len())
	}
}

package blockingdeque

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestBlocking verifies that operations block when expected.
func TestBlocking(t *testing.T) {
	t.Run("PushBlocksWhenFull", func(t *testing.T) {
		bd := New[int](1)
		bd.PushBack(1)

		start := time.Now()
		done := make(chan bool)
		go func() {
			time.Sleep(50 * time.Millisecond)
			bd.PopFront()
			close(done)
		}()

		// Should block until PopFront runs
		bd.PushBack(2)
		elapsed := time.Since(start)
		<-done

		if elapsed < 50*time.Millisecond {
			t.Error("PushBack should have blocked waiting for space")
		}
	})

	t.Run("PopBlocksWhenEmpty", func(t *testing.T) {
		bd := New[int](1)

		start := time.Now()
		done := make(chan bool)
		go func() {
			time.Sleep(50 * time.Millisecond)
			bd.PushBack(1)
			close(done)
		}()

		// Should block until PushBack runs
		val := bd.PopFront()
		elapsed := time.Since(start)
		<-done

		if elapsed < 50*time.Millisecond {
			t.Error("PopFront should have blocked waiting for data")
		}
		if val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}
	})
}

func TestContext(t *testing.T) {
	t.Run("PushCtxCancels", func(t *testing.T) {
		bd := New[int](1)
		bd.PushBack(1) // Full

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		err := bd.PushBackCtx(ctx, 2)
		if err == nil {
			t.Error("Expected error on context cancel")
		}
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})

	t.Run("PopCtxCancels", func(t *testing.T) {
		bd := New[int](1) // Empty

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		_, err := bd.PopFrontCtx(ctx)
		if err == nil {
			t.Error("Expected error on context cancel")
		}
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})
}

// TestConcurrency verifies thread safety under heavy load.
func TestConcurrency(t *testing.T) {
	bd := New[int](10) // Small buffer to force contention
	const (
		count   = 1000
		workers = 10
	)

	var wg sync.WaitGroup
	wg.Add(workers * 4) // 2 sets of producers, 2 sets of consumers

	// Producers Front
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				bd.PushFront(j)
			}
		}()
	}

	// Producers Back
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				bd.PushBack(j)
			}
		}()
	}

	// Consumers Front & Back
	var mu sync.Mutex
	var totalPopped int

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				bd.PopFront()
				mu.Lock()
				totalPopped++
				mu.Unlock()
			}
		}()
	}

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				bd.PopBack()
				mu.Lock()
				totalPopped++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	expected := workers * 2 * count
	if totalPopped != expected {
		t.Errorf("Expected %d items, got %d", expected, totalPopped)
	}
	if bd.Len() != 0 {
		t.Errorf("Queue should be empty, got len %d", bd.Len())
	}
}

func BenchmarkPushPop(b *testing.B) {
	bd := New[int](1024)
	go func() {
		for i := 0; ; i++ {
			bd.PushBack(i)
		}
	}()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bd.PopFront()
	}
}

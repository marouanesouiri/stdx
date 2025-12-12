package scheduler

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSchedulerBasic(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	executed := atomic.Bool{}
	s.Schedule(50*time.Millisecond, func() {
		executed.Store(true)
	})

	time.Sleep(100 * time.Millisecond)
	if !executed.Load() {
		t.Error("task was not executed")
	}
}

func TestSchedulerOrder(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	var mu sync.Mutex
	order := []int{}

	s.Schedule(100*time.Millisecond, func() {
		mu.Lock()
		order = append(order, 3)
		mu.Unlock()
	})

	s.Schedule(30*time.Millisecond, func() {
		mu.Lock()
		order = append(order, 1)
		mu.Unlock()
	})

	s.Schedule(60*time.Millisecond, func() {
		mu.Lock()
		order = append(order, 2)
		mu.Unlock()
	})

	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(order) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(order))
	}
	if order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Errorf("tasks executed in wrong order: %v", order)
	}
}

func TestSchedulerCancel(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	executed := atomic.Bool{}
	id := s.Schedule(50*time.Millisecond, func() {
		executed.Store(true)
	})

	cancelled := s.Cancel(id)
	if !cancelled {
		t.Error("cancel returned false")
	}

	time.Sleep(100 * time.Millisecond)
	if executed.Load() {
		t.Error("cancelled task was executed")
	}
}

func TestSchedulerConcurrent(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	count := atomic.Int32{}
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Schedule(10*time.Millisecond, func() {
				count.Add(1)
			})
		}()
	}

	wg.Wait()
	time.Sleep(50 * time.Millisecond)

	if count.Load() != 100 {
		t.Errorf("expected 100 tasks, executed %d", count.Load())
	}
}

func TestSchedulerPending(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	s.Schedule(100*time.Millisecond, func() {})
	s.Schedule(200*time.Millisecond, func() {})
	s.Schedule(300*time.Millisecond, func() {})

	pending := s.Pending()
	if pending != 3 {
		t.Errorf("expected 3 pending tasks, got %d", pending)
	}

	time.Sleep(350 * time.Millisecond)

	pending = s.Pending()
	if pending != 0 {
		t.Errorf("expected 0 pending tasks after execution, got %d", pending)
	}
}

func TestSchedulerTimingAccuracy(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	start := time.Now()
	executed := make(chan time.Time, 1)

	s.Schedule(100*time.Millisecond, func() {
		executed <- time.Now()
	})

	execTime := <-executed
	elapsed := execTime.Sub(start)

	tolerance := 20 * time.Millisecond
	if elapsed < 100*time.Millisecond || elapsed > 100*time.Millisecond+tolerance {
		t.Errorf("timing inaccurate: expected ~100ms, got %v", elapsed)
	}
}

func TestScheduleAtPanicOnPastTime(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	defer func() {
		if r := recover(); r == nil {
			t.Error("ScheduleAt did not panic when scheduling in the past")
		}
	}()

	pastTime := time.Now().Add(-1 * time.Hour)
	s.ScheduleAt(pastTime, func() {})
}

func TestSchedulerClear(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	count := atomic.Int32{}

	// Schedule 5 tasks
	for i := 0; i < 5; i++ {
		s.Schedule(100*time.Millisecond, func() {
			count.Add(1)
		})
	}

	// Verify tasks are scheduled
	if pending := s.Pending(); pending != 5 {
		t.Errorf("expected 5 pending tasks, got %d", pending)
	}

	// Clear all tasks
	s.Clear()

	// Verify queue is empty
	if pending := s.Pending(); pending != 0 {
		t.Errorf("expected 0 pending tasks after clear, got %d", pending)
	}

	// Wait and verify no tasks executed
	time.Sleep(150 * time.Millisecond)
	if count.Load() != 0 {
		t.Errorf("expected 0 tasks to execute after clear, got %d", count.Load())
	}
}

func BenchmarkSchedule(b *testing.B) {
	s := New()
	s.Start()
	defer s.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Schedule(time.Hour, func() {})
	}
}

func BenchmarkScheduleAndExecute(b *testing.B) {
	s := New()
	s.Start()
	defer s.Stop()

	done := make(chan struct{})
	count := atomic.Int32{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Schedule(time.Millisecond, func() {
			if count.Add(1) == int32(b.N) {
				close(done)
			}
		})
	}

	<-done
}

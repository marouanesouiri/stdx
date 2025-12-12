package scheduler

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"
)

// Scheduler manages scheduled tasks using a single goroutine.
// It efficiently schedules many tasks with minimal resource overhead.
type Scheduler struct {
	tasks   taskHeap
	mu      sync.Mutex
	wakeup  chan struct{}
	stopCh  chan struct{}
	running atomic.Bool
	nextID  atomic.Uint64
}

// New creates a new Scheduler.
// Call Start() to begin processing scheduled tasks.
func New() *Scheduler {
	s := &Scheduler{
		tasks:  make(taskHeap, 0),
		wakeup: make(chan struct{}, 1),
		stopCh: make(chan struct{}),
	}
	heap.Init(&s.tasks)
	return s
}

// Start launches the scheduler's goroutine.
// This must be called before scheduling tasks.
// Calling Start() on an already running scheduler has no effect.
func (s *Scheduler) Start() {
	if s.running.Swap(true) {
		return
	}
	go s.run()
}

// Stop gracefully stops the scheduler.
// Pending tasks will not be executed.
// Calling Stop() on an already stopped scheduler has no effect.
func (s *Scheduler) Stop() {
	if !s.running.Swap(false) {
		return
	}
	close(s.stopCh)
}

// Schedule schedules a function to execute after the specified delay.
// Returns a TaskID that can be used to cancel the task.
//
// WARNING: Task execution time must be less than the gap to the next scheduled task.
// For example, if tasks are scheduled 10s apart, each can take up to ~10s.
// If tasks are 100ms apart, each must complete in < 100ms to avoid delays.
// For long-running work, spawn a goroutine inside the task function.
func (s *Scheduler) Schedule(delay time.Duration, fn func()) TaskID {
	return s.ScheduleAt(time.Now().Add(delay), fn)
}

// ScheduleAt schedules a function to execute at the specified time.
// Returns a TaskID that can be used to cancel the task.
//
// Panics if the specified time is before the current time.
//
// WARNING: Task execution time must be less than the gap to the next scheduled task.
// For example, if tasks are scheduled 10s apart, each can take up to ~10s.
// If tasks are 100ms apart, each must complete in < 100ms to avoid delays.
// For long-running work, spawn a goroutine inside the task function.
func (s *Scheduler) ScheduleAt(at time.Time, fn func()) TaskID {
	if at.Before(time.Now()) {
		panic("scheduler: cannot schedule task in the past")
	}

	id := TaskID(s.nextID.Add(1))
	task := newTask(id, at, fn)

	s.mu.Lock()
	wasEmpty := s.tasks.Len() == 0
	earliestBefore := !wasEmpty && s.tasks.peek() != nil

	s.tasks.push(task)

	isEarliest := s.tasks.peek().ID() == id
	s.mu.Unlock()

	if wasEmpty || (earliestBefore && isEarliest) {
		select {
		case s.wakeup <- struct{}{}:
		default:
		}
	}

	return id
}

// Cancel cancels a scheduled task by its ID.
// Returns true if a task with the given ID was found (may already be cancelled).
// The task will be skipped when its execution time arrives.
func (s *Scheduler) Cancel(id TaskID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks.Len() {
		if s.tasks[i].ID() == id {
			s.tasks[i].Cancel()
			return true
		}
	}
	return false
}

// Pending returns the number of tasks currently scheduled (including cancelled).
// Cancelled tasks are lazily removed from the queue.
func (s *Scheduler) Pending() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tasks.Len()
}

// Clear removes all scheduled tasks from the scheduler.
// Tasks are removed immediately and will not be executed.
// This operation is thread-safe and signals the scheduler to wake up.
func (s *Scheduler) Clear() {
	s.mu.Lock()
	s.tasks = make(taskHeap, 0)
	heap.Init(&s.tasks)
	s.mu.Unlock()

	select {
	case s.wakeup <- struct{}{}:
	default:
	}
}

// run is the main scheduler loop that executes in a single goroutine.
func (s *Scheduler) run() {
	var timer *time.Timer

	for {
		s.mu.Lock()

		for s.tasks.Len() > 0 && s.tasks.peek().IsCancelled() {
			s.tasks.pop()
		}

		if s.tasks.Len() == 0 {
			s.mu.Unlock()
			if timer != nil {
				timer.Stop()
			}

			select {
			case <-s.wakeup:
				continue
			case <-s.stopCh:
				return
			}
		}

		nextTask := s.tasks.peek()
		waitDuration := time.Until(nextTask.RunAt())
		s.mu.Unlock()

		if waitDuration <= 0 {
			s.mu.Lock()
			task := s.tasks.pop()
			s.mu.Unlock()

			if task != nil && !task.IsCancelled() {
				start := time.Now()
				task.Execute()
				executionTime := time.Since(start)

				_ = executionTime
			}
			continue
		}

		if timer == nil {
			timer = time.NewTimer(waitDuration)
		} else {
			timer.Reset(waitDuration)
		}

		select {
		case <-timer.C:
		case <-s.wakeup:
			if !timer.Stop() {
				<-timer.C
			}
		case <-s.stopCh:
			if timer != nil {
				timer.Stop()
			}
			return
		}
	}
}

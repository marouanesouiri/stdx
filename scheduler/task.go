package scheduler

import (
	"sync/atomic"
	"time"
)

// TaskID uniquely identifies a scheduled task.
// Can be used to cancel the task before it executes.
type TaskID uint64

// Task represents a scheduled function to be executed at a specific time.
type Task struct {
	id        TaskID
	runAt     time.Time
	fn        func()
	cancelled atomic.Bool
}

// newTask creates a new task with the given ID, execution time, and function.
func newTask(id TaskID, runAt time.Time, fn func()) *Task {
	return &Task{
		id:    id,
		runAt: runAt,
		fn:    fn,
	}
}

// ID returns the task's unique identifier.
func (t *Task) ID() TaskID {
	return t.id
}

// RunAt returns the scheduled execution time.
func (t *Task) RunAt() time.Time {
	return t.runAt
}

// Cancel marks the task as cancelled.
// The task will be skipped when its execution time arrives.
func (t *Task) Cancel() {
	t.cancelled.Store(true)
}

// IsCancelled returns true if the task has been cancelled.
func (t *Task) IsCancelled() bool {
	return t.cancelled.Load()
}

// Execute runs the task function if it hasn't been cancelled.
// Returns true if the task was executed, false if it was cancelled.
func (t *Task) Execute() bool {
	if t.IsCancelled() {
		return false
	}
	t.fn()
	return true
}

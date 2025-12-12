// Package scheduler provides an efficient task scheduler that uses a single goroutine
// to execute many scheduled tasks with minimal overhead.
//
// The scheduler uses a min-heap data structure for O(log n) insertions and O(1) access
// to the next task, making it highly efficient even with thousands of scheduled tasks.
//
// # Basic Usage
//
//	s := scheduler.New()
//	s.Start()
//	defer s.Stop()
//
//	// Schedule a task 5 seconds from now
//	id := s.Schedule(5*time.Second, func() {
//	    fmt.Println("Task executed!")
//	})
//
//	// Cancel a task before it executes
//	s.Cancel(id)
//
// # Performance Characteristics
//
//   - Memory: O(n) where n = number of scheduled tasks
//   - Schedule: O(log n) insertion into min-heap
//   - Next Task: O(1) access to heap root
//   - Cancel: O(1) atomic flag set (lazy removal)
//   - Goroutines: Exactly 1, regardless of task count
//
// # Important Warnings
//
// Task execution time must be less than the gap to the next scheduled task.
//
// Examples:
//   - Tasks scheduled 10s apart: each can take up to ~10s
//   - Tasks scheduled 100ms apart: each must complete in < 100ms
//   - Single task: can take as long as needed (no following task to delay)
//
// If a task runs longer than the gap, subsequent tasks will be delayed.
// For truly long-running work, spawn a goroutine inside the task:
//
//	// ❌ BAD - blocks scheduler, delays next task
//	s.Schedule(time.Second, func() {
//	    time.Sleep(10 * time.Second) // This delays all subsequent tasks!
//	    processData()
//	})
//
//	// ✅ GOOD - doesn't block scheduler
//	s.Schedule(time.Second, func() {
//	    go processData() // Returns immediately
//	})
//
//	// ✅ ALSO GOOD - if tasks are far apart
//	s.Schedule(1 * time.Hour, func() {
//	    processData() // Can take up to ~1 hour if next task is 1 hour away
//	})
//
// # Thread Safety
//
// The scheduler is safe for concurrent use. Multiple goroutines can schedule
// and cancel tasks simultaneously.
//
// # Smart Delay Calculation
//
// The scheduler intelligently adjusts sleep durations by subtracting:
//   - Time already elapsed waiting for previous tasks
//   - Execution time of previous tasks
//
// This minimizes delay between sequential tasks and ensures accurate timing.
package scheduler

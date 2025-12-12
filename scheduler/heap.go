package scheduler

import (
	"container/heap"
)

// taskHeap implements heap.Interface for tasks ordered by execution time.
// The task with the earliest runAt time is at the root (index 0).
type taskHeap []*Task

// Len returns the number of tasks in the heap.
func (h taskHeap) Len() int {
	return len(h)
}

// Less reports whether the task at index i should execute before the task at index j.
// Tasks are ordered by their runAt time (earliest first).
func (h taskHeap) Less(i, j int) bool {
	return h[i].runAt.Before(h[j].runAt)
}

// Swap exchanges the tasks at indices i and j.
func (h taskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Push adds a task to the heap.
// This method is called by heap.Push, not directly.
func (h *taskHeap) Push(x interface{}) {
	*h = append(*h, x.(*Task))
}

// Pop removes and returns the task with the earliest execution time.
// This method is called by heap.Pop, not directly.
func (h *taskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	task := old[n-1]
	*h = old[0 : n-1]
	return task
}

// peek returns the next task to execute without removing it.
// Returns nil if the heap is empty.
func (h *taskHeap) peek() *Task {
	if len(*h) == 0 {
		return nil
	}
	return (*h)[0]
}

// push adds a task to the heap and maintains heap invariant.
func (h *taskHeap) push(task *Task) {
	heap.Push(h, task)
}

// pop removes and returns the next task to execute.
// Returns nil if the heap is empty.
func (h *taskHeap) pop() *Task {
	if len(*h) == 0 {
		return nil
	}
	return heap.Pop(h).(*Task)
}

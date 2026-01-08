/*
Package blockingqueue implements a thread-safe queue backed by a Go channel.

The queue is strictly FIFO (First-In-First-Out). For a double-ended blocking
queue, see the blockingdeque package.

Example usage:

	bq := blockingqueue.New[string](10) // Bounded queue wrapper around make(chan string, 10)

	// Producer
	go func() {
		bq.Push("hello")
	}()

	// Consumer
	msg := bq.Pop() // Blocks on channel receive
*/
package blockingqueue

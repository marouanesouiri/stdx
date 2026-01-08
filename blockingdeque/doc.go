/*
Package blockingdeque implements a thread-safe double-ended queue (Deque).

It functions exactly like a buffered Go channel, but with the ability to push
and pop elements from both the front and the back. It supports blocking operations,
non-blocking "Try" operations, and integration with `context.Context` for cancellation.
*/
package blockingdeque

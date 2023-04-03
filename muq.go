package ring

import "sync"

// MuQ is a thread-safe FIFO queue backed by a ring buffer.
// The zero value for MuQ is an empty queue ready to use.
//
// MuQ is safe for concurrent use.
// If you need to use it from a single goroutine, use [Q] instead.
type MuQ[T any] struct {
	mu sync.RWMutex
	q  Q[T]
}

// The API for MuQ differs from Q somewhat:
//
// - There's no Do method because we don't want to hold onto a lock
//   while we wait for a user-specified callback
// - There's no panicking Pop or Peek method because
//   if the queue is read from concurrently,
//   verifying that it's non-empty and removing the entry
//   must be a single atomic operation,
//   which means users will only need Try variants,
//   and having the panicking versions will just cause bugs.

// NewMuQ returns a new thread-safe queue with the given capacity.
func NewMuQ[T any](capacity int) *MuQ[T] {
	var m MuQ[T]
	m.q.init(capacity)
	return &m
}

// Empty returns true if the queue is empty.
//
// This is an O(1) operation and does not allocate.
func (q *MuQ[T]) Empty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.q.Empty()
}

// Len returns the number of items in the queue.
//
// This is an O(1) operation and does not allocate.
func (q *MuQ[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.q.Len()
}

// Clear removes all items from the queue.
// It does not adjust its internal capacity.
//
// This is an O(1) operation and does not allocate.
func (q *MuQ[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.q.Clear()
}

// Push adds x to the back of the queue.
//
// This operation is O(n) in the worst case if the queue needs to grow.
// However, for target use cases, it's amortized O(1).
// See package documentation for details.
// If your usage pattern is bursts of pushes followed by bursts of pops,
func (q *MuQ[T]) Push(x T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.q.Push(x)
}

// TryPop removes and returns the item at the front of the queue.
// It returns false if the queue is empty.
// Otherwise, it returns true and the item.
//
// This is an O(1) operation and does not allocate.
func (q *MuQ[T]) TryPop() (x T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.q.TryPop()
}

// TryPeek returns the item at the front of the queue.
// It returns false if the queue is empty.
// Otherwise, it returns true and the item.
//
// This is an O(1) operation and does not allocate.
func (q *MuQ[T]) TryPeek() (x T, ok bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.q.TryPeek()
}

// Snapshot appends the contents of the queue to dst and returns the result.
//
// Use dst to avoid allocations when you know the capacity of the queue
// or pass nil to let the function allocate a new slice.
//
// The returned slice is a copy of the internal buffer and is safe to modify.
func (q *MuQ[T]) Snapshot(dst []T) []T {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.q.Snapshot(dst)
}

// Package ring implements a FIFO queue backed by a ring buffer.
//
// Its main advantage over a queue backed by a slice or container/list
// is that it requires significantly fewer allocations
// for specific usage patterns.
//
//   - consistent rate of pushes and pops
//   - bursts of pushes followed by bursts of pops
//
// # Implementation
//
// It's valuable to understand how the queue works under the hood
// so that you can reason about its performance characteristics
// and where it might be useful.
//
// The queue is backed by a ring buffer: a slice that wraps around.
// When the queue is full, the ring buffer grows by doubling its capacity.
//
// A ring buffer is a slice with two pointers: head and tail.
// The contents of the queue are usually [head:tail].
// The queue is empty if head == tail.
//
//	+-------------------------------+
//	| ... | head | ... | tail | ... |
//	+-------------------------------+
//	      '............'
//
// Every pop moves the head forward.
// Every push moves the tail forward.
//
// When the tail reaches the end of the current buffer,
// it wraps around to the beginning if there's room.
// At this point, the contents of the queue are [head:] + [:tail].
//
//	+-------------------------------+
//	| ... | tail | ... | head | ... |
//	+-------------------------------+
//	......'            '.............
//
// If there's no room for a push,
// the ring buffer grows by doubling its capacity.
//
// When the head reaches the end of the current buffer,
// it wraps around to the beginning too, continuing to chase the tail.
// If it catches up with the tail, the queue is full.
package ring

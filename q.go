package ring

const _defaultCapacity = 16

// Q is a FIFO queue backed by a ring buffer.
// The zero value for Q is an empty queue ready to use.
//
// Q is not safe for concurrent use.
// If you need to use it from multiple goroutines, use [MuQ] instead.
type Q[T any] struct {
	// buff is the ring buffer.
	//
	// The first item in the queue is at buff[head].
	// The last item in the queue is at buff[tail-1].
	// The queue is empty if head == tail.
	buff []T

	// head is the index of the first item in the queue.
	head int // inv: 0 <= head < len(buff)

	// tail is the index of the next empty slot in buff.
	tail int // inv: 0 <= tail < len(buff)
}

// NewQ returns a new queue with the given capacity.
// If capacity is zero, the queue is initialized with a default capacity.
//
// The capacity defines the leeway for bursts of pushes
// before the queue needs to grow.
func NewQ[T any](capacity int) *Q[T] {
	var q Q[T]
	q.init(capacity)
	return &q
}

func (q *Q[T]) init(capacity int) {
	if capacity == 0 {
		capacity = _defaultCapacity
	}
	// Allocate requested capacity plus one slot
	// so that filling the queue to exactly the requested capacity
	// doesn't require resizing.
	q.buff = make([]T, capacity+1)
	q.head = 0
	q.tail = 0
}

// Empty returns true if the queue is empty.
//
// This is an O(1) operation and does not allocate.
func (q *Q[T]) Empty() bool {
	return q.head == q.tail
}

// Len returns the number of items in the queue.
//
// This is an O(1) operation and does not allocate.
func (q *Q[T]) Len() int {
	if q.head <= q.tail {
		return q.tail - q.head
	}
	return len(q.buff) - q.head + q.tail
}

// Clear removes all items from the queue.
// It does not adjust its internal capacity.
//
// This is an O(1) operation and does not allocate.
func (q *Q[T]) Clear() {
	q.head = 0
	q.tail = 0
}

// Push adds x to the back of the queue.
//
// This operation is O(n) in the worst case if the queue needs to grow.
// However, for target use cases, it's amortized O(1).
// See package documentation for details.
func (q *Q[T]) Push(x T) {
	if len(q.buff) == 0 {
		q.buff = make([]T, _defaultCapacity)
	}

	q.buff[q.tail] = x
	q.tail++

	if q.tail == len(q.buff) {
		// Wrap around.
		q.tail = 0
	}

	// We'll hit this only if the tail has wrapped around
	// and has caught up with the head (the queue is full).
	// In that case, we need to grow the queue
	// copying buff[head:] and buff[:tail] to the new buffer.
	if q.head == q.tail {
		// The queue is full. Make room.
		buff := make([]T, 2*len(q.buff))
		n := copy(buff, q.buff[q.head:])
		n += copy(buff[n:], q.buff[:q.tail])
		q.head = 0
		q.tail = n
		q.buff = buff
	}
}

// Pop removes and returns the item at the front of the queue.
// It panics if the queue is empty.
//
// This is an O(1) operation and does not allocate.
func (q *Q[T]) Pop() T {
	x, ok := q.TryPop()
	if !ok {
		panic("empty queue")
	}
	return x
}

// TryPop removes and returns the item at the front of the queue.
// It returns false if the queue is empty.
// Otherwise, it returns true and the item.
//
// This is an O(1) operation and does not allocate.
func (q *Q[T]) TryPop() (x T, ok bool) {
	if q.head == q.tail {
		return x, false
	}

	x = q.buff[q.head]
	q.head++
	if q.head == len(q.buff) {
		// Wrap around.
		//
		// If tail has wrapped around too,
		// the next Pop will catch it when head == tail.
		q.head = 0
	}
	return x, true
}

// Peek returns the item at the front of the queue without removing it.
// It panics if the queue is empty.
//
// This is an O(1) operation and does not allocate.
func (q *Q[T]) Peek() T {
	x, ok := q.TryPeek()
	if !ok {
		panic("empty queue")
	}
	return x
}

// TryPeek returns the item at the front of the queue.
// It returns false if the queue is empty.
// Otherwise, it returns true and the item.
//
// This is an O(1) operation and does not allocate.
func (q *Q[T]) TryPeek() (x T, ok bool) {
	if q.head == q.tail {
		return x, false
	}
	return q.buff[q.head], true
}

// Snapshot appends the contents of the queue to dst and returns the result.
// Use dst to avoid allocations when you know the capacity of the queue
//
//	dst := make([]T, 0, q.Len())
//	dst = q.Snapshot(dst)
//
// Pass nil to let the function allocate a new slice.
//
//	q.Snapshot(nil) // allocates a new slice
//
// The returned slice is a copy of the internal buffer and is safe to modify.
func (q *Q[T]) Snapshot(dst []T) []T {
	if q.head <= q.tail {
		return append(dst, q.buff[q.head:q.tail]...)
	}

	dst = append(dst, q.buff[q.head:]...)
	return append(dst, q.buff[:q.tail]...)
}

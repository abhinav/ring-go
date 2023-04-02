package ring

// Exports some internals for testing.

// Items returns the items in the queue as a slice.
//
// Do not modify the returned slice
// as it may be shared with the queue's internal buffer.
func (q *Q[T]) Items() []T {
	if q.head <= q.tail {
		return q.buff[q.head:q.tail]
	}
	return append(q.buff[q.head:], q.buff[:q.tail]...)
}

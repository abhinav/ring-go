package ring_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.abhg.dev/container/ring"
)

func TestQ_empty(t *testing.T) {
	t.Parallel()

	var q ring.Q[int]
	assert.True(t, q.Empty(), "empty")
	assert.Zero(t, q.Len(), "len")
	assert.Panics(t, func() { q.Peek() }, "peek")
	assert.Panics(t, func() { q.Pop() }, "pop")
	assert.Empty(t, q.Snapshot(nil), "snapshot")
}

func TestQ_PeekPop(t *testing.T) {
	t.Parallel()

	var q ring.Q[int]
	q.Push(42)
	assert.Equal(t, 42, q.Peek(), "peek")
	assert.Equal(t, 42, q.Pop(), "pop")
	assert.True(t, q.Empty(), "empty")
}

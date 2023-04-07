package ring_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.abhg.dev/container/ring"
)

func TestQ_empty(t *testing.T) {
	t.Parallel()

	var q ring.Q[int]
	assert.Panics(t, func() { q.Peek() }, "peek")
	assert.Panics(t, func() { q.Pop() }, "pop")

	q.Do(func(i int) bool {
		t.Errorf("unexpected item: %v", i)
		return true
	})
}

func TestQ_PeekPop(t *testing.T) {
	t.Parallel()

	var q ring.Q[int]
	q.Push(42)
	assert.Equal(t, 42, q.Peek(), "peek")
	assert.Equal(t, 42, q.Pop(), "pop")
	assert.True(t, q.Empty(), "empty")
}

func TestQ_Do(t *testing.T) {
	t.Parallel()

	const NumItems = 100
	var q ring.Q[int]

	want := make([]int, 0, NumItems)
	for i := 0; i < NumItems; i++ {
		q.Push(i)
		want = append(want, i)
	}

	got := make([]int, 0, NumItems)
	q.Do(func(i int) bool {
		got = append(got, i)
		return true
	})

	assert.Equal(t, want, got, "do did not iterate fully")
}

func TestQ_Do_returnEarly(t *testing.T) {
	t.Parallel()

	const NumItems = 100
	var q ring.Q[int]

	stopAt := NumItems / 2
	want := make([]int, 0, NumItems)
	for i := 0; i < NumItems; i++ {
		q.Push(i)
		if i < stopAt {
			want = append(want, i)
		}
	}

	got := make([]int, 0, stopAt)
	q.Do(func(i int) bool {
		if i >= stopAt {
			return false
		}

		got = append(got, i)
		return true
	})

	assert.Equal(t, want, got, "iterated over unexpected items")
}

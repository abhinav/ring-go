package ring_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.abhg.dev/container/ring"
)

func TestQ(t *testing.T) {
	t.Parallel()

	testQueueSuite(t, func(capacity int) *ring.Q[int] {
		return ring.NewQ[int](capacity)
	})
}

func testQueueSuite(t *testing.T, newWithCap func(capacity int) *ring.Q[int]) {
	capacities := []int{
		-1, // special case: use zero value
		0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024,
	}
	sizes := []int{1, 10, 100, 1000, 10000}

	for _, capacity := range capacities {
		for _, size := range sizes {
			require.Greater(t, size, 0,
				"invalid test: sizes must be greater than 0")

			capacity, size := capacity, size
			name := fmt.Sprintf("Capacity=%d/Size=%d", capacity, size)
			newEmpty := func() *ring.Q[int] {
				if capacity < 0 {
					return new(ring.Q[int])
				}
				return newWithCap(capacity)
			}

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				suite := &queueSuite{
					NewEmpty: newEmpty,
					NumItems: size,
				}

				t.Run("Empty", suite.TestEmpty)
				t.Run("PushPop", suite.TestPushPop)
				t.Run("PushPopInterleaved", suite.TestPushPopInterleaved)
				t.Run("PushPopWraparound", suite.TestPushPopWraparound)
				t.Run("Snapshot", suite.TestSnapshot)
				t.Run("SnapshotReuse", suite.TestSnapshotReuse)
			})
		}
	}
}

type queueSuite struct {
	NewEmpty func() *ring.Q[int]
	NumItems int
}

func (s *queueSuite) TestEmpty(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	assert.True(t, q.Empty(), "empty")
	assert.Zero(t, q.Len(), "length")

	assert.Panics(t, func() {
		q.Pop()
	}, "pop")
	assert.Panics(t, func() {
		q.Peek()
	}, "pop")

	q.Do(func(item int) bool {
		t.Errorf("unexpected item: %v", item)
		return true
	})
}

func (s *queueSuite) TestPushPop(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
	}
	assert.False(t, q.Empty(), "empty")
	assert.Equal(t, s.NumItems, q.Len(), "length")

	t.Run("Do", func(t *testing.T) {
		want := 0
		q.Do(func(item int) bool {
			assert.Equal(t, want, item, "item")
			want++
			return true
		})
	})

	for i := 0; i < s.NumItems; i++ {
		assert.Equal(t, i, q.Peek(), "peek")
		assert.Equal(t, i, q.Pop(), "pop")
	}

	assert.True(t, q.Empty(), "empty")
	assert.Zero(t, q.Len(), "length")
}

func (s *queueSuite) TestPushPopInterleaved(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
		q.Do(func(item int) bool {
			assert.Equal(t, i, item, "item")
			return true
		})
		assert.Equal(t, i, q.Peek(), "peek")
		assert.Equal(t, i, q.Pop(), "pop")
	}

	assert.True(t, q.Empty(), "empty")
	assert.Zero(t, q.Len(), "length")
}

func (s *queueSuite) TestPushPopWraparound(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
		q.Push(q.Pop())
	}

	got := make([]int, 0, q.Len())
	for !q.Empty() {
		got = append(got, q.Pop())
	}
	sort.Ints(got)

	want := make([]int, 0, s.NumItems)
	for i := 0; i < s.NumItems; i++ {
		want = append(want, i)
	}

	assert.Equal(t, want, got, "items")
}

func (s *queueSuite) TestSnapshot(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
	}

	snap := q.Snapshot(nil /* dst */)
	assert.Len(t, snap, q.Len(), "length")
	for _, item := range snap {
		assert.Equal(t, item, q.Pop(), "item")
	}
}

func (s *queueSuite) TestSnapshotReuse(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
	}

	snap := []int{42}
	snap = q.Snapshot(snap)
	assert.Len(t, snap, q.Len()+1, "length")

	assert.Equal(t, 42, snap[0], "item")
	for _, item := range snap[1:] {
		assert.Equal(t, item, q.Pop(), "item")
	}
}

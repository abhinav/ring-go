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

	capacities := []int{
		-1, // special case: use zero value
		0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024,
	}
	sizes := []int{1, 10, 100, 1000, 10000}

	for _, capacity := range capacities {
		for _, size := range sizes {
			capacity, size := capacity, size
			name := fmt.Sprintf("capacity=%d/items=%d", capacity, size)
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				require.Greater(t, size, 0,
					"invalid test: sizes must be greater than 0")

				testQCases(t, func() *ring.Q[int] {
					if capacity < 0 {
						return new(ring.Q[int])
					}
					return ring.NewQ[int](capacity)
				}, size)
			})
		}
	}
}

func testQCases(t *testing.T, newEmpty func() *ring.Q[int], NumItems int) {
	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		q := newEmpty()
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
	})

	t.Run("push all then pop", func(t *testing.T) {
		t.Parallel()

		q := newEmpty()
		for i := 0; i < NumItems; i++ {
			q.Push(i)
		}
		assert.False(t, q.Empty(), "empty")
		assert.Equal(t, NumItems, q.Len(), "length")

		t.Run("Do", func(t *testing.T) {
			want := 0
			q.Do(func(item int) bool {
				assert.Equal(t, want, item, "item")
				want++
				return true
			})
		})

		for i := 0; i < NumItems; i++ {
			assert.Equal(t, i, q.Peek(), "peek")
			assert.Equal(t, i, q.Pop(), "pop")
		}

		assert.True(t, q.Empty(), "empty")
		assert.Zero(t, q.Len(), "length")
	})

	t.Run("push and pop interleaved", func(t *testing.T) {
		t.Parallel()

		q := newEmpty()
		for i := 0; i < NumItems; i++ {
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
	})

	t.Run("push and pop with wraparound", func(t *testing.T) {
		t.Parallel()

		q := newEmpty()
		for i := 0; i < NumItems; i++ {
			q.Push(i)
			q.Push(q.Pop())
		}

		got := make([]int, 0, q.Len())
		for !q.Empty() {
			got = append(got, q.Pop())
		}
		sort.Ints(got)

		want := make([]int, 0, NumItems)
		for i := 0; i < NumItems; i++ {
			want = append(want, i)
		}

		assert.Equal(t, want, got, "items")
	})
}

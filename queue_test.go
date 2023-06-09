package ring_test

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.abhg.dev/container/ring"
)

func TestQ(t *testing.T) {
	t.Parallel()

	testQueueSuite(t, func(capacity int) queue[int] {
		return ring.NewQ[int](capacity)
	})
}

func TestMuQ(t *testing.T) {
	t.Parallel()

	testQueueSuite(t, func(capacity int) queue[int] {
		return ring.NewMuQ[int](capacity)
	})
}

type queue[T any] interface {
	Empty() bool
	Len() int
	Clear()
	Push(x T)
	TryPop() (T, bool)
	TryPeek() (T, bool)
	Snapshot([]T) []T
}

var (
	_ queue[int] = (*ring.Q[int])(nil)
	_ queue[int] = (*ring.MuQ[int])(nil)
)

func testQueueSuite(t *testing.T, newWithCap func(capacity int) queue[int]) {
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
			newEmpty := func() queue[int] {
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

				suitev := reflect.ValueOf(suite)
				suitet := suitev.Type()
				for i := 0; i < suitet.NumMethod(); i++ {
					name, ok := cutPrefix(suitet.Method(i).Name, "Test")
					if !ok {
						continue
					}

					testfn, ok := suitev.Method(i).Interface().(func(*testing.T))
					if !ok {
						continue
					}

					t.Run(name, testfn)
				}
			})
		}
	}
}

type queueSuite struct {
	NewEmpty func() queue[int]
	NumItems int
}

func (s *queueSuite) TestEmpty(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	assert.True(t, q.Empty(), "empty")
	assert.Zero(t, q.Len(), "length")

	t.Run("TryPeekPop", func(t *testing.T) {
		_, ok := q.TryPeek()
		assert.False(t, ok, "peek should fail")

		_, ok = q.TryPop()
		assert.False(t, ok, "pop should fail")
	})

	assert.Empty(t, q.Snapshot(nil), "snapshot")
}

func (s *queueSuite) TestPushPop(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
	}
	assert.False(t, q.Empty(), "empty")
	assert.Equal(t, s.NumItems, q.Len(), "length")

	for i := 0; i < s.NumItems; i++ {
		assert.Equal(t, i, requirePeek(t, q), "peek")
		assert.Equal(t, i, requirePop(t, q), "pop")
	}

	assert.True(t, q.Empty(), "empty")
	assert.Zero(t, q.Len(), "length")
}

func (s *queueSuite) TestPushPopInterleaved(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
		assert.Equal(t, i, requirePeek(t, q), "peek")
		assert.Equal(t, i, requirePop(t, q), "pop")
	}

	assert.True(t, q.Empty(), "empty")
	assert.Zero(t, q.Len(), "length")
}

func (s *queueSuite) TestPushPopWraparound(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
		q.Push(requirePop(t, q))
	}

	got := make([]int, 0, q.Len())
	for !q.Empty() {
		got = append(got, requirePop(t, q))
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
		assert.Equal(t, item, requirePop(t, q), "item")
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
		assert.Equal(t, item, requirePop(t, q), "item")
	}
}

func requirePeek[T any](t require.TestingT, q queue[T]) T {
	v, ok := q.TryPeek()
	require.True(t, ok, "peek")
	return v
}

func requirePop[T any](t require.TestingT, q queue[T]) T {
	v, ok := q.TryPop()
	require.True(t, ok, "pop")
	return v
}

// Copy of strings.CutPrefix for Go 1.19.
// Delete once Go 1.20 is minimum supported version.
func cutPrefix(s, prefix string) (after string, found bool) {
	if !strings.HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
}

package ring_test

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.abhg.dev/container/ring"
	"pgregory.net/rapid"
)

func TestQ_rapid(t *testing.T) {
	t.Parallel()

	rapid.Check(t, rapid.Run[*qMachine]())
}

type qMachine struct {
	q *ring.Q[int]

	golden *list.List
}

var _ rapid.StateMachine = (*qMachine)(nil)

func (m *qMachine) Init(t *rapid.T) {
	capacity := rapid.IntRange(1, 100).Draw(t, "capacity")
	m.q = ring.NewQ[int](capacity)
	m.golden = list.New()
}

func (m *qMachine) Check(t *rapid.T) {
	assert.Equal(t, m.q.Len(), m.golden.Len())

	got := make([]int, 0, m.q.Len())
	got = m.q.Snapshot(got)

	for i, e := 0, m.golden.Front(); e != nil; i, e = i+1, e.Next() {
		assert.Equal(t, e.Value, got[i])
	}
}

func (m *qMachine) Push(t *rapid.T) {
	x := rapid.Int().Draw(t, "x")
	m.q.Push(x)
	m.golden.PushBack(x)
}

func (m *qMachine) Pop(t *rapid.T) {
	if m.q.Empty() {
		t.Skip()
	}

	m.q.Pop()
	m.golden.Remove(m.golden.Front())
}

func (m *qMachine) Peek(t *rapid.T) {
	if m.q.Empty() {
		t.Skip()
	}

	got := m.q.Peek()
	assert.Equal(t, m.golden.Front().Value, got)
}

func (m *qMachine) Clear(t *rapid.T) {
	m.q.Clear()
	m.golden.Init()
}

func (m *qMachine) Empty(t *rapid.T) {
	assert.Equal(t, m.golden.Len() == 0, m.q.Empty())
}

func (m *qMachine) Len(t *rapid.T) {
	assert.Equal(t, m.golden.Len(), m.q.Len())
}

func (m *qMachine) Snapshot(t *rapid.T) {
	got := make([]int, 0, m.q.Len())
	got = m.q.Snapshot(got)

	for i, e := 0, m.golden.Front(); e != nil; i, e = i+1, e.Next() {
		assert.Equal(t, e.Value, got[i])
	}
}

func (m *qMachine) Do(t *rapid.T) {
	if m.q.Empty() {
		t.Skip()
	}

	el := m.golden.Front()
	m.q.Do(func(x int) bool {
		assert.Equal(t, el.Value, x)
		if rapid.Bool().Draw(t, "proceed") {
			el = el.Next()
			return true
		}
		return false
	})
}

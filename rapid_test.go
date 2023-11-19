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

	rapid.Check(t, func(t *rapid.T) {
		t.Repeat(rapid.StateMachineActions(newQMachine[*ring.Q[int]](t)))
	})
}

func TestMuQ_rapid(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		t.Repeat(rapid.StateMachineActions(newQMachine[*ring.MuQ[int]](t)))
	})
}

type qMachine[QT queue[int]] struct {
	q QT

	golden *list.List
}

var _ rapid.StateMachine = (*qMachine[queue[int]])(nil)

func newQMachine[QT queue[int]](t *rapid.T) *qMachine[QT] {
	capacity := rapid.IntRange(1, 100).Draw(t, "capacity")

	var q queue[int]
	switch queue[int](*new(QT)).(type) {
	case *ring.Q[int]:
		q = ring.NewQ[int](capacity)
	case *ring.MuQ[int]:
		q = ring.NewMuQ[int](capacity)
	default:
		t.Fatalf("cannot instantiate queue type: %T", *new(QT))
	}

	return &qMachine[QT]{
		q:      q.(QT),
		golden: list.New(),
	}
}

func (m *qMachine[QT]) Check(t *rapid.T) {
	assert.Equal(t, m.q.Len(), m.golden.Len())

	got := make([]int, 0, m.q.Len())
	got = m.q.Snapshot(got)

	for i, e := 0, m.golden.Front(); e != nil; i, e = i+1, e.Next() {
		assert.Equal(t, e.Value, got[i])
	}
}

func (m *qMachine[QT]) Push(t *rapid.T) {
	x := rapid.Int().Draw(t, "x")
	m.q.Push(x)
	m.golden.PushBack(x)
}

func (m *qMachine[QT]) Pop(t *rapid.T) {
	if m.golden.Len() == 0 {
		t.Skip()
	}

	want := m.golden.Remove(m.golden.Front())
	got := requirePop(t, queue[int](m.q))
	assert.Equal(t, want, got)
}

func (m *qMachine[QT]) TryPop(t *rapid.T) {
	got, ok := m.q.TryPop()

	front := m.golden.Front()
	if front == nil {
		assert.False(t, ok)
	} else {
		assert.Equal(t, m.golden.Remove(front), got)
	}
}

func (m *qMachine[QT]) Peek(t *rapid.T) {
	if m.golden.Len() == 0 {
		t.Skip()
	}

	got := requirePeek(t, queue[int](m.q))
	assert.Equal(t, m.golden.Front().Value, got)
}

func (m *qMachine[QT]) TryPeek(t *rapid.T) {
	got, ok := m.q.TryPeek()

	front := m.golden.Front()
	if front == nil {
		assert.False(t, ok)
		return
	}

	assert.Equal(t, front.Value, got)
}

func (m *qMachine[QT]) Clear(_ *rapid.T) {
	m.q.Clear()
	m.golden.Init()
}

func (m *qMachine[QT]) Empty(t *rapid.T) {
	assert.Equal(t, m.golden.Len() == 0, m.q.Empty())
}

func (m *qMachine[QT]) Len(t *rapid.T) {
	assert.Equal(t, m.golden.Len(), m.q.Len())
}

func (m *qMachine[QT]) Snapshot(t *rapid.T) {
	got := make([]int, 0, m.q.Len())
	got = m.q.Snapshot(got)

	for i, e := 0, m.golden.Front(); e != nil; i, e = i+1, e.Next() {
		assert.Equal(t, e.Value, got[i])
	}
}

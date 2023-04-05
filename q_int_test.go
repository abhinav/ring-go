package ring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Verifies that a queue filled exactly to capacity does not resize.
func TestQ_fillNoResize(t *testing.T) {
	t.Parallel()

	q := NewQ[int](3)
	initCap := cap(q.buff)
	q.Push(1)
	q.Push(2)
	q.Push(3)
	assert.Equal(t, initCap, cap(q.buff), "capacity")
}

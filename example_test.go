package ring_test

import (
	"fmt"

	"go.abhg.dev/container/ring"
)

func ExampleQ_Pop_loop() {
	var q ring.Q[int]
	for i := 0; i < 3; i++ {
		q.Push(i)
	}

	for !q.Empty() {
		fmt.Println(q.Pop())
	}

	// Output:
	// 0
	// 1
	// 2
}

func ExampleQ_TryPop_loop() {
	var q ring.Q[int]
	for i := 0; i < 3; i++ {
		q.Push(i)
	}

	for v, ok := q.TryPop(); ok; v, ok = q.TryPop() {
		fmt.Println(v)
	}

	// Output:
	// 0
	// 1
	// 2
}

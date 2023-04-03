package ring_test

import (
	"fmt"
	"testing"

	"go.abhg.dev/container/ring"
)

func BenchmarkPushPop_sameRate(b *testing.B) {
	var q ring.Q[int]
	for i := 0; i < b.N; i++ {
		q.Push(i)
		q.Pop()
	}
}

func BenchmarkPush_burst(b *testing.B) {
	bursts := []int{1, 10, 100}

	for _, burst := range bursts {
		name := fmt.Sprintf("burst=%d", burst)
		b.Run(name, func(b *testing.B) {
			var q ring.Q[int]
			for i := 0; i < b.N; i++ {
				for j := 0; j < burst; j++ {
					q.Push(i)
				}
				for j := 0; j < burst; j++ {
					q.Pop()
				}
			}
		})
	}
}

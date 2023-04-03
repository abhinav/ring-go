package ring_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.abhg.dev/container/ring"
)

func BenchmarkQ_pushPop_sameRate(b *testing.B) {
	var q ring.Q[int]
	for i := 0; i < b.N; i++ {
		q.Push(i)
		q.Pop()
	}
}

func BenchmarkMuQ_pushPop_sameRate(b *testing.B) {
	var q ring.MuQ[int]
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			q.Push(i)
			_, ok := q.TryPop()
			require.True(b, ok)
		}
	})
}

func BenchmarkQ_push_burst(b *testing.B) {
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

func BenchmarkMuQ_push_burst(b *testing.B) {
	bursts := []int{1, 10, 100}

	for _, burst := range bursts {
		name := fmt.Sprintf("burst=%d", burst)
		b.Run(name, func(b *testing.B) {
			var q ring.MuQ[int]
			b.RunParallel(func(pb *testing.PB) {
				for i := 0; pb.Next(); i++ {
					for j := 0; j < burst; j++ {
						q.Push(i)
					}
					for j := 0; j < burst; j++ {
						_, ok := q.TryPop()
						require.True(b, ok)
					}
				}
			})
		})
	}
}

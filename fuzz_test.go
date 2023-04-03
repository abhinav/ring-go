package ring_test

import (
	"testing"

	"go.abhg.dev/container/ring"
	"pgregory.net/rapid"
)

func FuzzQ_rapid(f *testing.F) {
	f.Fuzz(rapid.MakeFuzz(rapid.Run[*qMachine[*ring.Q[int]]]()))
}

func FuzzMuQ_rapid(f *testing.F) {
	f.Fuzz(rapid.MakeFuzz(rapid.Run[*qMachine[*ring.MuQ[int]]]()))
}

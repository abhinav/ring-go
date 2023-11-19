package ring_test

import (
	"testing"

	"go.abhg.dev/container/ring"
	"pgregory.net/rapid"
)

func FuzzQ_rapid(f *testing.F) {
	f.Fuzz(rapid.MakeFuzz(func(t *rapid.T) {
		t.Repeat(rapid.StateMachineActions(newQMachine[*ring.Q[int]](t)))
	}))
}

func FuzzMuQ_rapid(f *testing.F) {
	f.Fuzz(rapid.MakeFuzz(func(t *rapid.T) {
		t.Repeat(rapid.StateMachineActions(newQMachine[*ring.MuQ[int]](t)))
	}))
}

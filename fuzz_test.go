package ring_test

import (
	"testing"

	"pgregory.net/rapid"
)

func FuzzQ_rapid(f *testing.F) {
	f.Fuzz(rapid.MakeFuzz(rapid.Run[*qMachine]()))
}

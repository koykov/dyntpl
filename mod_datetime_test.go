package dyntpl

import "testing"

func TestModDatetime(t *testing.T) {
	t.Run("modNow", testMod)
}

func BenchmarkModDatetime(b *testing.B) {
	b.Run("modNow", benchMod)
}

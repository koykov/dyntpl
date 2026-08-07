package dyntpl

import (
	"math"
	"testing"
)

func TestModEncode(t *testing.T) {
	t.Run("hexStr", func(t *testing.T) { testModWA(t, modArgs{"value": "lorem ipsum dolor sit..."}) })
	t.Run("hexInt", func(t *testing.T) { testModWA(t, modArgs{"value": -int(math.MaxInt - 5e12)}) })
	t.Run("hexUint", func(t *testing.T) { testModWA(t, modArgs{"value": uint(math.MaxUint / 123456)}) })
	t.Run("hexFloat", func(t *testing.T) { testModWA(t, modArgs{"value": math.Pi}) })
}

func BenchmarkModEncode(b *testing.B) {
	b.Run("hexStr", func(b *testing.B) { benchModWA(b, modArgs{"value": "lorem ipsum dolor sit..."}) })
	b.Run("hexInt", func(b *testing.B) { benchModWA(b, modArgs{"value": -int(math.MaxInt - 5e12)}) })
	b.Run("hexUint", func(b *testing.B) { benchModWA(b, modArgs{"value": uint(math.MaxUint / 123456)}) })
	b.Run("hexFloat", func(b *testing.B) { benchModWA(b, modArgs{"value": math.Pi}) })
}

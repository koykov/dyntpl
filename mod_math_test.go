package dyntpl

import "testing"

func TestModMath(t *testing.T) {
	t.Run("abs", func(t *testing.T) { testModWA(t, modArgs{"num": -123456}) })
	t.Run("add", func(t *testing.T) { testModWA(t, modArgs{"num": 555}) })
	t.Run("sub", func(t *testing.T) { testModWA(t, modArgs{"num": 555}) })
	t.Run("inc", func(t *testing.T) { testModWA(t, modArgs{"num": 10}) })
	t.Run("dec", func(t *testing.T) { testModWA(t, modArgs{"num": 10}) })
	t.Run("mul", func(t *testing.T) { testModWA(t, modArgs{"num": 10}) })
	t.Run("div", func(t *testing.T) { testModWA(t, modArgs{"num": 10}) })
	t.Run("mod", func(t *testing.T) { testModWA(t, modArgs{"num": 10}) })
}

func BenchmarkModMath(b *testing.B) {
	b.Run("abs", func(b *testing.B) { benchModWA(b, modArgs{"num": -123456}) })
	b.Run("add", func(b *testing.B) { benchModWA(b, modArgs{"num": 555}) })
	b.Run("sub", func(b *testing.B) { benchModWA(b, modArgs{"num": 555}) })
	b.Run("inc", func(b *testing.B) { benchModWA(b, modArgs{"num": 10}) })
	b.Run("dec", func(b *testing.B) { benchModWA(b, modArgs{"num": 10}) })
	b.Run("mul", func(b *testing.B) { benchModWA(b, modArgs{"num": 10}) })
	b.Run("div", func(b *testing.B) { benchModWA(b, modArgs{"num": 10}) })
	b.Run("mod", func(b *testing.B) { benchModWA(b, modArgs{"num": 10}) })
}

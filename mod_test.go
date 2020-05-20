package dyntpl

import "testing"

var (
	tplModDef       = []byte(`Cost is: {%= user.Cost|default(999.99) %} USD`)
	tplModDefStatic = []byte(`{% ctx defaultCost = 999.99 %}Cost is: {%= user.Cost|default(defaultCost) %} USD`)
	expectModDef    = []byte(`Cost is: 999.99 USD`)
)

func TestTplModDef(t *testing.T) {
	testBase(t, "tplModDef", expectModDef, "mod def tpl mismatch")
}

func TestTplModDefStatic(t *testing.T) {
	testBase(t, "tplModDefStatic", expectModDef, "mod def static tpl mismatch")
}

func BenchmarkTplModDef(b *testing.B) {
	benchBase(b, "tplModDef", expectModDef, "mod def tpl mismatch")
}

func BenchmarkTplModDefStatic(b *testing.B) {
	benchBase(b, "tplModDefStatic", expectModDef, "mod def static tpl mismatch")
}

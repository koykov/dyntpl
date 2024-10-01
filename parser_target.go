package dyntpl

// Target is a storage of depths needed to provide proper out from conditions, loops and switches control structures.
type target struct {
	cc, cl, cs int
}

// Check if parser reached the target.
func (t *target) reached(p *parser) bool {
	return t.cc == p.cc &&
		t.cl == p.cl &&
		t.cs == p.cs
}

// Check if target is a root.
func (t *target) eqZero() bool {
	return t.cc == 0 &&
		t.cl == 0 &&
		t.cs == 0
}

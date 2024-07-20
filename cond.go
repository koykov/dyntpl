package dyntpl

// CondFn describes helper func signature.
type CondFn func(ctx *Ctx, args []any) bool

type CondFnTuple struct {
	docgen
	fn CondFn
}

var (
	// Registry of condition helpers.
	condRegistry = map[string]int{}
	condBuf      []CondFnTuple
)

// RegisterCondFn registers new condition helper.
func RegisterCondFn(name string, cond CondFn) *CondFnTuple {
	if idx, ok := condRegistry[name]; ok && idx >= 0 && idx < len(condBuf) {
		return &condBuf[idx]
	}
	condBuf = append(condBuf, CondFnTuple{
		docgen: docgen{name: name},
		fn:     cond,
	})
	idx := len(condBuf) - 1
	condRegistry[name] = idx
	return &condBuf[idx]
}

// RegisterCondFnNS registers new condition helper in given namespace.
func RegisterCondFnNS(namespace, name string, cond CondFn) *CondFnTuple {
	if len(namespace) > 0 {
		name = namespace + "::" + name
	}
	return RegisterCondFn(name, cond)
}

// GetCondFn returns condition helper from the registry.
func GetCondFn(name string) CondFn {
	if idx, ok := condRegistry[name]; ok && idx >= 0 && idx < len(condBuf) {
		return condBuf[idx].fn
	}
	return nil
}

var _ = RegisterCondFnNS

package dyntpl

// CondFn describes helper func signature.
type CondFn func(ctx *Ctx, args []any) bool

var (
	// Registry of condition helpers.
	condRegistry = map[string]CondFn{}
)

// RegisterCondFn registers new condition helper in registry.
func RegisterCondFn(name string, cond CondFn) {
	condRegistry[name] = cond
}

// GetCondFn returns condition helper from the registry.
func GetCondFn(name string) *CondFn {
	if fn, ok := condRegistry[name]; ok {
		return &fn
	}
	return nil
}

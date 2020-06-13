package dyntpl

// Condition helper func signature.
type CondFn func(ctx *Ctx, args []interface{}) bool

var (
	// Registry of condition helpers.
	condRegistry = map[string]CondFn{}
)

// Register new condition helper.
func RegisterCondFn(name string, cond CondFn) {
	condRegistry[name] = cond
}

// Get condition helper from the registry.
func GetCondFn(name string) *CondFn {
	if fn, ok := condRegistry[name]; ok {
		return &fn
	}
	return nil
}

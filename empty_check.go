package dyntpl

// Empty check helper func signature.
type EmptyCheckFn func(ctx *Ctx, val interface{}) bool

var (
	// Registry of empty check helpers.
	emptyCheckRegistry = map[string]EmptyCheckFn{}
	// Suppress go vet warning.
	_ = GetEmptyCheckFn
)

// Register new empty check helper.
func RegisterEmptyCheckFn(name string, cond EmptyCheckFn) {
	emptyCheckRegistry[name] = cond
}

// Get empty check helper from the registry.
func GetEmptyCheckFn(name string) *EmptyCheckFn {
	if fn, ok := emptyCheckRegistry[name]; ok {
		return &fn
	}
	return nil
}

// General empty check func.
func EmptyCheck(ctx *Ctx, val interface{}) bool {
	if val == nil {
		return true
	}
	for _, fn := range emptyCheckRegistry {
		if fn(ctx, val) {
			return true
		}
	}
	return false
}

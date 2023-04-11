package dyntpl

// EmptyCheckFn describes empty check helper func signature.
type EmptyCheckFn func(ctx *Ctx, val any) bool

var (
	// Registry of empty check helpers.
	emptyCheckRegistry = map[string]EmptyCheckFn{}
	// Suppress go vet warning.
	_ = GetEmptyCheckFn
)

// RegisterEmptyCheckFn registers new empty check helper.
func RegisterEmptyCheckFn(name string, cond EmptyCheckFn) {
	emptyCheckRegistry[name] = cond
}

// GetEmptyCheckFn gets empty check helper from the registry.
func GetEmptyCheckFn(name string) *EmptyCheckFn {
	if fn, ok := emptyCheckRegistry[name]; ok {
		return &fn
	}
	return nil
}

// EmptyCheck tries to apply all known helpers over the val.
//
// First acceptable helper will break next attempts.
func EmptyCheck(ctx *Ctx, val any) bool {
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

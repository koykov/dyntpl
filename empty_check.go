package dyntpl

// EmptyCheckFn describes empty check helper func signature.
type EmptyCheckFn func(ctx *Ctx, val any) bool

var (
	// Registry of empty check helpers.
	emptyCheckRegistry = map[string]EmptyCheckFn{}
)

// RegisterEmptyCheckFn registers new empty check helper.
func RegisterEmptyCheckFn(name string, cond EmptyCheckFn) {
	emptyCheckRegistry[name] = cond
}

// RegisterEmptyCheckFnNS registers new empty check helper.
func RegisterEmptyCheckFnNS(namespace, name string, cond EmptyCheckFn) {
	if len(namespace) > 0 {
		name = namespace + "::" + name
	}
	RegisterEmptyCheckFn(name, cond)
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

var _, _ = GetEmptyCheckFn, RegisterEmptyCheckFnNS

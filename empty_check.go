package dyntpl

// EmptyCheckFn describes empty check helper func signature.
type EmptyCheckFn func(ctx *Ctx, val any) bool

type EmptyCheckTuple struct {
	docgen
	fn EmptyCheckFn
}

var (
	// Registry of empty check helpers.
	emptyCheckRegistry = map[string]int{}
	emptyCheckBuf      []EmptyCheckTuple
)

// RegisterEmptyCheckFn registers new empty check helper.
func RegisterEmptyCheckFn(name string, cond EmptyCheckFn) *EmptyCheckTuple {
	if idx, ok := emptyCheckRegistry[name]; ok && idx >= 0 && idx < len(emptyCheckBuf) {
		return &emptyCheckBuf[idx]
	}
	emptyCheckBuf = append(emptyCheckBuf, EmptyCheckTuple{
		docgen: docgen{name: name},
		fn:     cond,
	})
	idx := len(emptyCheckBuf) - 1
	emptyCheckRegistry[name] = idx
	return &emptyCheckBuf[idx]
}

// RegisterEmptyCheckFnNS registers new empty check helper.
func RegisterEmptyCheckFnNS(namespace, name string, cond EmptyCheckFn) {
	if len(namespace) > 0 {
		name = namespace + "::" + name
	}
	RegisterEmptyCheckFn(name, cond)
}

// GetEmptyCheckFn gets empty check helper from the registry.
func GetEmptyCheckFn(name string) EmptyCheckFn {
	if idx, ok := emptyCheckRegistry[name]; ok && idx >= 0 && idx < len(emptyCheckBuf) {
		return emptyCheckBuf[idx].fn
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
	for i := 0; i < len(emptyCheckBuf); i++ {
		if emptyCheckBuf[i].fn(ctx, val) {
			return true
		}
	}
	return false
}

var _, _ = GetEmptyCheckFn, RegisterEmptyCheckFnNS

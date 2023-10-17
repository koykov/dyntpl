package dyntpl

import (
	"strconv"

	"github.com/koykov/fastconv"
)

// ModFn describes signature of the modifier functions.
//
// Arguments description:
// * ctx provides access to additional variables and various buffers to reduce allocations.
// * buf is a storage for final result after finishing modifier work.
// * val is a left side variable that preceded to call of modifier func, example: {%= val|mod(...) %}
// * args is a list of all arguments listed on modifier call.
type ModFn func(ctx *Ctx, buf *any, val any, args []any) error

// Internal modifier representation.
type mod struct {
	id  []byte
	fn  *ModFn
	arg []*arg
}

var (
	// Registry of modifiers.
	modRegistry = map[string]ModFn{}
)

// RegisterModFn registers new modifier function.
func RegisterModFn(name, alias string, mod ModFn) {
	modRegistry[name] = mod
	if len(alias) > 0 {
		modRegistry[alias] = mod
	}
}

// RegisterModFnNS registers new mod function in given namespace.
func RegisterModFnNS(namespace, name, alias string, mod ModFn) {
	if len(namespace) == 0 {
		RegisterModFn(name, alias, mod)
	}
	name = namespace + "::" + name
	if len(alias) > 0 {
		alias = namespace + "::" + alias
	}
	RegisterModFn(name, alias, mod)
}

// GetModFn gets modifier from the registry.
func GetModFn(name string) *ModFn {
	if fn, ok := modRegistry[name]; ok {
		return &fn
	}
	return nil
}

// Get count of print iterations.
func printIterations(args []any) int {
	itr := 1
	if len(args) > 0 {
		if itrRaw, ok := args[0].(*[]byte); ok {
			if itr64, err := strconv.ParseInt(fastconv.B2S(*itrRaw), 10, 64); err == nil {
				itr = int(itr64)
			}
		}
	}
	return itr
}

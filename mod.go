package dyntpl

import (
	"strconv"

	"github.com/koykov/byteconv"
)

// ModFn describes signature of the modifier functions.
//
// Arguments description:
// * ctx provides access to additional variables and various buffers to reduce allocations.
// * buf is a storage for final result after finishing modifier work.
// * val is a left side variable that preceded to call of modifier func, example: {%= val|mod(...) %}
// * args is a list of all arguments listed on modifier call.
type ModFn func(ctx *Ctx, buf *any, val any, args []any) error

type modFnParam struct {
	param string
	desc  string
}

type modFnTuple struct {
	id      string
	alias   string
	desc    string
	params  []modFnParam
	example string
	fn      ModFn
}

// Internal modifier representation.
type mod struct {
	id  []byte
	fn  ModFn
	arg []*arg
}

var (
	// Registry of modifiers.
	modRegistry = map[string]int{}
	modBuf      []modFnTuple
)

func (t *modFnTuple) WithDescription(desc string) *modFnTuple {
	t.desc = desc
	return t
}

func (t *modFnTuple) WithParam(param, desc string) *modFnTuple {
	t.params = append(t.params, modFnParam{
		param: param,
		desc:  desc,
	})
	return t
}

func (t *modFnTuple) WithExample(example string) *modFnTuple {
	t.example = example
	return t
}

// RegisterModFn registers new modifier function.
func RegisterModFn(name, alias string, mod ModFn) *modFnTuple {
	if idx, ok := modRegistry[name]; ok && idx >= 0 && idx < len(modBuf) {
		return &modBuf[idx]
	}
	modBuf = append(modBuf, modFnTuple{
		id:    name,
		alias: alias,
		fn:    mod,
	})
	idx := len(modBuf) - 1
	modRegistry[name] = idx
	if len(alias) > 0 {
		modRegistry[alias] = idx
	}
	return &modBuf[idx]
}

// RegisterModFnNS registers new mod function in given namespace.
func RegisterModFnNS(namespace, name, alias string, mod ModFn) {
	if len(namespace) == 0 {
		RegisterModFn(name, alias, mod)
		return
	}
	name = namespace + "::" + name
	if len(alias) > 0 {
		alias = namespace + "::" + alias
	}
	RegisterModFn(name, alias, mod)
}

// GetModFn gets modifier from the registry.
func GetModFn(name string) ModFn {
	if idx, ok := modRegistry[name]; ok && idx >= 0 && idx < len(modBuf) {
		return modBuf[idx].fn
	}
	return nil
}

// Get count of print iterations.
func printIterations(args []any) int {
	itr := 1
	if len(args) > 0 {
		if itrRaw, ok := args[0].(*[]byte); ok {
			if itr64, err := strconv.ParseInt(byteconv.B2S(*itrRaw), 10, 64); err == nil {
				itr = int(itr64)
			}
		}
	}
	return itr
}

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

type ModFnTuple struct {
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

func (t *ModFnTuple) WithDescription(desc string) *ModFnTuple {
	t.desc = desc
	return t
}

func (t *ModFnTuple) WithParam(param, desc string) *ModFnTuple {
	t.params = append(t.params, modFnParam{
		param: param,
		desc:  desc,
	})
	return t
}

func (t *ModFnTuple) WithExample(example string) *ModFnTuple {
	t.example = example
	return t
}

var (
	// Registry of modifiers.
	modRegistry = map[string]int{}
	modBuf      []ModFnTuple
)

// RegisterModFn registers new modifier function.
func RegisterModFn(name, alias string, mod ModFn) *ModFnTuple {
	if idx, ok := modRegistry[name]; ok && idx >= 0 && idx < len(modBuf) {
		return &modBuf[idx]
	}
	modBuf = append(modBuf, ModFnTuple{
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
func RegisterModFnNS(namespace, name, alias string, mod ModFn) *ModFnTuple {
	if len(namespace) == 0 {
		return RegisterModFn(name, alias, mod)
	}
	name = namespace + "::" + name
	if len(alias) > 0 {
		alias = namespace + "::" + alias
	}
	return RegisterModFn(name, alias, mod)
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

package dyntpl

import "github.com/koykov/inspector/testobj"

// CondOKFn describes helper func signature.
type CondOKFn func(ctx *Ctx, v *any, ok *bool, args []any)

var (
	// Registry of condition-OK helpers.
	condOKRegistry = map[string]CondOKFn{}
)

// RegisterCondOKFn registers new condition-OK helper.
func RegisterCondOKFn(name string, cond CondOKFn) {
	condOKRegistry[name] = cond
}

// RegisterCondOKFnNS registers new condition-OK helper in given namespace.
func RegisterCondOKFnNS(namespace, name string, cond CondOKFn) {
	if len(namespace) > 0 {
		name = namespace + "::" + name
	}
	RegisterCondOKFn(name, cond)
}

// GetCondOKFn returns condition-OK helper from the registry.
func GetCondOKFn(name string) *CondOKFn {
	if fn, ok := condOKRegistry[name]; ok {
		return &fn
	}
	return nil
}

// Simple example of condition-OK helper func.
func testCondOK(ctx *Ctx, v *any, ok *bool, args []any) {
	if len(args) == 0 {
		*ok = false
		return
	}
	if fin, ok1 := args[0].(**testobj.TestFinance); ok1 {
		c := ctx.GetCounter("__testUserNextHistory999counter")
		if c >= len((*fin).History) {
			*ok = false
			return
		}
		*v = &(*fin).History[c]
		*ok = true
		c++
		ctx.SetCounter("__testUserNextHistory999counter", c)
	}
}

var _ = RegisterCondOKFnNS

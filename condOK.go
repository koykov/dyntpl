package dyntpl

import "github.com/koykov/inspector/testobj"

// CondOKFn describes helper func signature.
type CondOKFn func(ctx *Ctx, v *interface{}, ok *bool, args []interface{})

var (
	// Registry of condition-OK helpers.
	condOKRegistry = map[string]CondOKFn{}
)

// RegisterCondOKFn registers new condition-OK helper in registry.
func RegisterCondOKFn(name string, cond CondOKFn) {
	condOKRegistry[name] = cond
}

// GetCondOKFn returns condition-OK helper from the registry.
func GetCondOKFn(name string) *CondOKFn {
	if fn, ok := condOKRegistry[name]; ok {
		return &fn
	}
	return nil
}

// Simple example of condition-OK helper func.
func testCondOK(ctx *Ctx, v *interface{}, ok *bool, args []interface{}) {
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

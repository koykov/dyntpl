package dyntpl

import "github.com/koykov/inspector/testobj"

// Condition helper func signature.
type CondOKFn func(ctx *Ctx, v *interface{}, ok *bool, args []interface{})

var (
	// Registry of condition-ok helpers.
	condOKRegistry = map[string]CondOKFn{}
)

// Register new condition-ok helper.
func RegisterCondOKFn(name string, cond CondOKFn) {
	condOKRegistry[name] = cond
}

// Get condition-ok helper from the registry.
func GetCondOKFn(name string) *CondOKFn {
	if fn, ok := condOKRegistry[name]; ok {
		return &fn
	}
	return nil
}

// Simple example of condition-ok helper func.
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

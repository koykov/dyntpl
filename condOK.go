package dyntpl

import "github.com/koykov/inspector/testobj"

// CondOKFn describes helper func signature.
type CondOKFn func(ctx *Ctx, v *any, ok *bool, args []any)

type CondOKTuple struct {
	docgen
	fn CondOKFn
}

var (
	// Registry of condition-OK helpers.
	condOKRegistry = map[string]int{}
	condOkBuf      []CondOKTuple
)

// RegisterCondOKFn registers new condition-OK helper.
func RegisterCondOKFn(name string, cond CondOKFn) *CondOKTuple {
	if idx, ok := condOKRegistry[name]; ok && idx >= 0 && idx < len(condOkBuf) {
		return &condOkBuf[idx]
	}
	condOkBuf = append(condOkBuf, CondOKTuple{
		docgen: docgen{name: name},
		fn:     cond,
	})
	idx := len(condOkBuf) - 1
	condOKRegistry[name] = idx
	return &condOkBuf[idx]
}

// RegisterCondOKFnNS registers new condition-OK helper in given namespace.
func RegisterCondOKFnNS(namespace, name string, cond CondOKFn) *CondOKTuple {
	if len(namespace) > 0 {
		name = namespace + "::" + name
	}
	return RegisterCondOKFn(name, cond)
}

// GetCondOKFn returns condition-OK helper from the registry.
func GetCondOKFn(name string) CondOKFn {
	if idx, ok := condOKRegistry[name]; ok && idx >= 0 && idx < len(condOkBuf) {
		return condOkBuf[idx].fn
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

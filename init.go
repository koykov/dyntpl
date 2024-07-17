package dyntpl

import (
	"github.com/koykov/clock"
	"github.com/koykov/inspector/testobj_ins"
)

func init() {
	// Register simple builtin modifiers.
	RegisterModFn("default", "def", modDefault).
		WithDescription("The `default` modifier returns the passed default value if the value is undefined or empty, otherwise the value of the variable.").
		WithParam("args ...any").
		WithExample(`{%= notExistingOrEmptyVar|default("N/D") %}"`)
	RegisterModFn("ifThen", "if", modIfThen).
		WithDescription("Modifier `ifThen` passes `arg` only if preceding condition is true.").
		WithParam("arg any").
		WithExample(`{%= obj.Active|ifThen("<button>Buy</button>") %}`)
	RegisterModFn("ifThenElse", "ifel", modIfThenElse).
		WithDescription("Modifier 'ifTheElse` passes `arg0` if preceding condition is true or `arg1` otherwise.").
		WithParam("arg0 any").
		WithParam("arg1 any").
		WithExample(`{%= user.AllowSell|ifThenElse("<button>Sell</button>", "not available!") %}`)

	// Register builtin escape/quote modifiers.
	RegisterModFn("jsonEscape", "je", modJSONEscape)
	RegisterModFn("jsonQuote", "jq", modJSONQuote)
	RegisterModFn("htmlEscape", "he", modHTMLEscape)
	RegisterModFn("linkEscape", "le", modLinkEscape)
	RegisterModFn("urlEncode", "ue", modURLEncode)
	RegisterModFn("attrEscape", "ae", modAttrEscape)
	RegisterModFn("cssEscape", "ce", modCSSEscape)
	RegisterModFn("jsEscape", "jse", modJSEscape)
	RegisterModFn("raw", "noesc", func(_ *Ctx, _ *any, _ any, _ []any) error { return nil })

	// Register builtin round modifiers.
	RegisterModFn("round", "round", modRound)
	RegisterModFn("roundPrec", "roundp", modRoundPrec)
	RegisterModFn("ceil", "ceil", modCeil)
	RegisterModFn("ceilPrec", "ceilp", modCeilPrec)
	RegisterModFn("floor", "floor", modFloor)
	RegisterModFn("floorPrec", "floorp", modFloorPrec)

	// Register time modifiers.
	RegisterModFnNS("time", "now", "", modNow)
	RegisterModFnNS("time", "format", "date", modDate)
	RegisterModFnNS("time", "add", "date_modify", modDateAdd)

	// Register math modifiers.
	RegisterModFnNS("math", "abs", "", modAbs)
	RegisterModFnNS("math", "inc", "", modInc)
	RegisterModFnNS("math", "dec", "", modDec)
	RegisterModFnNS("math", "add", "", modMathAdd)
	RegisterModFnNS("math", "sub", "", modMathSub)
	RegisterModFnNS("math", "mul", "", modMathMul)
	RegisterModFnNS("math", "div", "", modMathDiv)
	RegisterModFnNS("math", "mod", "", modMathMod)
	RegisterModFnNS("math", "sqrt", "", modMathSqrt)
	RegisterModFnNS("math", "cbrt", "", modMathCbrt)
	RegisterModFnNS("math", "radical", "rad", modMathRadical)
	RegisterModFnNS("math", "exp", "", modMathExp)
	RegisterModFnNS("math", "log", "", modMathLog)
	RegisterModFnNS("math", "factorial", "fact", modMathFact)
	RegisterModFnNS("math", "max", "", modMathMax)
	RegisterModFnNS("math", "min", "", modMathMin)
	RegisterModFnNS("math", "pow", "", modMathPow)

	// Register builtin condition helpers.
	RegisterCondFn("lenEq0", condLenEq0)
	RegisterCondFn("lenGt0", condLenGt0)
	RegisterCondFn("lenGtq0", condLenGtq0)

	// Register builtin empty check helpers.
	RegisterEmptyCheckFn("int", EmptyCheckInt)
	RegisterEmptyCheckFn("uint", EmptyCheckUint)
	RegisterEmptyCheckFn("float", EmptyCheckFloat)
	RegisterEmptyCheckFn("bytes", EmptyCheckBytes)
	RegisterEmptyCheckFn("bytes_slice", EmptyCheckBytesSlice)
	RegisterEmptyCheckFn("str", EmptyCheckStr)
	RegisterEmptyCheckFn("str_slice", EmptyCheckStrSlice)
	RegisterEmptyCheckFn("bool", EmptyCheckBool)

	// Register datetime layouts.
	RegisterGlobalNS("time", "Layout", "", clock.Layout)
	RegisterGlobalNS("time", "ANSIC", "", clock.ANSIC)
	RegisterGlobalNS("time", "UnixDate", "", clock.UnixDate)
	RegisterGlobalNS("time", "RubyDate", "", clock.RubyDate)
	RegisterGlobalNS("time", "RFC822", "", clock.RFC822)
	RegisterGlobalNS("time", "RFC822Z", "", clock.RFC822Z)
	RegisterGlobalNS("time", "RFC850", "", clock.RFC850)
	RegisterGlobalNS("time", "RFC1123", "", clock.RFC1123)
	RegisterGlobalNS("time", "RFC1123Z", "", clock.RFC1123Z)
	RegisterGlobalNS("time", "RFC3339", "", clock.RFC3339)
	RegisterGlobalNS("time", "RFC3339Nano", "", clock.RFC3339Nano)
	RegisterGlobalNS("time", "Kitchen", "", clock.Kitchen)
	RegisterGlobalNS("time", "Stamp", "", clock.Stamp)
	RegisterGlobalNS("time", "StampMilli", "", clock.StampMilli)
	RegisterGlobalNS("time", "StampMicro", "", clock.StampMicro)
	RegisterGlobalNS("time", "StampNano", "", clock.StampNano)

	// Register test modifiers.
	RegisterModFn("testNameOf", "", modTestNameOf)
	RegisterModFnNS("testns", "pack", "", func(_ *Ctx, _ *any, _ any, _ []any) error { return nil })
	RegisterModFnNS("testns", "extract", "", func(_ *Ctx, _ *any, _ any, _ []any) error { return nil })
	RegisterModFnNS("testns", "marshal", "", func(_ *Ctx, _ *any, _ any, _ []any) error { return nil })
	RegisterModFnNS("testns", "modCB", "", func(ctx *Ctx, _ *any, _ any, args []any) error {
		ctx.SetStatic("testVar", args[0])
		return nil
	})

	// Register test condition-ok helpers.
	RegisterCondOKFn("__testUserNextHistory999", testCondOK)

	// Register test variable-inspector pairs.
	RegisterVarInsPair("__testFin999", &testobj_ins.TestFinanceInspector{})
}

package dyntpl

import "github.com/koykov/inspector/testobj_ins"

func init() {
	// Register simple builtin modifiers.
	RegisterModFn("default", "def", modDefault)
	RegisterModFn("ifThen", "if", modIfThen)
	RegisterModFn("ifThenElse", "ifel", modIfThenElse)

	// Register i18n modifiers.
	RegisterModFn("translate", "t", modTranslate)
	RegisterModFn("translatePlural", "tp", modTranslatePlural)

	// Register builtin escape/quote modifiers.
	RegisterModFn("jsonEscape", "je", modJSONEscape)
	RegisterModFn("jsonQuote", "jq", modJSONQuote)
	RegisterModFn("htmlEscape", "he", modHTMLEscape)
	RegisterModFn("linkEscape", "le", modLinkEscape)
	RegisterModFn("urlEncode", "ue", modURLEncode)
	RegisterModFn("attrEscape", "ae", modAttrEscape)
	RegisterModFn("cssEscape", "ce", modCSSEscape)
	RegisterModFn("jsEscape", "jse", modJSEscape)

	// Register builtin round modifiers.
	RegisterModFn("round", "round", modRound)
	RegisterModFn("roundPrec", "roundp", modRoundPrec)
	RegisterModFn("ceil", "ceil", modCeil)
	RegisterModFn("ceilPrec", "ceilp", modCeilPrec)
	RegisterModFn("floor", "floor", modFloor)
	RegisterModFn("floorPrec", "floorp", modFloorPrec)

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

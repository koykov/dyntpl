package dyntpl

import "github.com/koykov/any2bytes"

func init() {
	RegisterModFn("default", "def", modDefault)
	RegisterModFn("ifThen", "if", modIfThen)
	RegisterModFn("ifThenElse", "ifel", modIfThenElse)

	RegisterModFn("jsonEscape", "je", modJsonEscape)
	RegisterModFn("jsonQuote", "jq", modJsonQuote)
	RegisterModFn("htmlEscape", "he", modHtmlEscape)

	RegisterModFn("round", "round", modRound)
	RegisterModFn("roundPrec", "roundp", modRoundPrec)
	RegisterModFn("ceil", "ceil", modCeil)
	RegisterModFn("ceilPrec", "ceilp", modCeilPrec)
	RegisterModFn("floor", "floor", modFloor)
	RegisterModFn("floorPrec", "floorp", modFloorPrec)

	RegisterCondFn("lenEq0", condLenEq0)
	RegisterCondFn("lenGt0", condLenGt0)
	RegisterCondFn("lenGtq0", condLenGtq0)

	any2bytes.RegisterAnyToBytesFn(ByteBufToBytes)
}

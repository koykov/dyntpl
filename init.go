package dyntpl

func init() {
	RegisterModFn("default", "def", modDefault)
	RegisterModFn("ifThen", "if", modIfThen)
	RegisterModFn("ifThenElse", "ifel", modIfThenElse)
	RegisterModFn("jsonQuote", "jq", modJsonQuote)
	RegisterModFn("htmlEscape", "he", modHtmlEscape)
}

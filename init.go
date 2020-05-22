package dyntpl

func init() {
	RegisterModFn("default", "def", modDefault)
	RegisterModFn("ifThen", "if", modIfThen)
	RegisterModFn("ifThenElse", "ifel", modIfThenElse)

	RegisterModFn("jsonEscape", "je", modJsonEscape)
	RegisterModFn("jsonQuote", "jq", modJsonQuote)
	RegisterModFn("htmlEscape", "he", modHtmlEscape)
}

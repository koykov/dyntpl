package dyntpl

func init() {
	RegisterModFn("default", "def", modDefault)
	RegisterModFn("jsonQuote", "jq", modJsonQuote)
	RegisterModFn("htmlEscape", "he", modHtmlEscape)
}

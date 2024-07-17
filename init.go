package dyntpl

import (
	"github.com/koykov/clock"
	"github.com/koykov/inspector/testobj_ins"
)

func init() {
	// Register simple builtin modifiers.
	RegisterModFn("default", "def", modDefault).
		WithDescription("Modifier `default` returns the passed `arg` if the preceding value is undefined or empty, otherwise the value of the variable.").
		WithParam("arg any", "").
		WithExample(`{%= notExistingOrEmptyVar|default("N/D") %}"`)
	RegisterModFn("ifThen", "if", modIfThen).
		WithDescription("Modifier `ifThen` passes `arg` only if preceding condition is true.").
		WithParam("arg any", "").
		WithExample(`{%= obj.Active|ifThen("<button>Buy</button>") %}`)
	RegisterModFn("ifThenElse", "ifel", modIfThenElse).
		WithDescription("Modifier `ifTheElse` passes `arg0` if preceding condition is true or `arg1` otherwise.").
		WithParam("arg0 any", "").
		WithParam("arg1 any", "").
		WithExample(`{%= user.AllowSell|ifThenElse("<button>Sell</button>", "not available!") %}`)

	// Register builtin escape/quote modifiers.
	RegisterModFn("jsonEscape", "je", modJSONEscape).
		WithDescription("Applies JSON escaping to printing value. May work together with `{% jsonquote %}...{% endjsonquote %}`.").
		WithExample(`{"name":"{%= user.name|jsonEscape %}"} // {"name":"Foo\"bar"}`)
	RegisterModFn("jsonQuote", "jq", modJSONQuote).
		WithDescription("Applies JSON quotation to printing value.").
		WithExample(`{"name":{%= user.name|jsonQuote %}} // {"name":"Foo\"bar"}`)
	RegisterModFn("htmlEscape", "he", modHTMLEscape).
		WithDescription("Applies HTML escaping to printing value.").
		WithExample(`<span data-title="{%= title|htmlEscape %}">{%= text|he %}</span> // <span data-title="&lt;h1&gt;Go is an open source programming language that makes it easy to build &lt;strong&gt;simple&lt;strong&gt;, &lt;strong&gt;reliable&lt;/strong&gt;, and &lt;strong&gt;efficient&lt;/strong&gt; software.&lt;/h1&gt;">Show more &gt;</span>`)
	RegisterModFn("linkEscape", "le", modLinkEscape).
		WithDescription("Applies Link escaping to printing value.").
		WithExample(`<a href="{%l= link %}"> // <a href="http://x.com/link-with-\"-and+space-symbol">`)
	RegisterModFn("urlEncode", "ue", modURLEncode).
		WithDescription("Applies URL encoding to printing value.").
		WithExample(`<a href="https://redir.com/{%u= url %}">go to >>></a> // <a href="https://redir.com/https%3A%2F%2Fgolang.org%2Fsrc%2Fnet%2Furl%2Furl.go%23L100">go to >>></a>`)
	RegisterModFn("attrEscape", "ae", modAttrEscape).
		WithDescription("Applies Attribute escaping to printing value.").
		WithExample(`<span font='{%a= var1 %}'> // <span font='foo&amp;&lt;&gt;&quot;&#x27;&#x60;&#x21;&#x40;&#x24;&#x25;&#x28;&#x29;&#x3d;&#x2b;&#x7b;&#x7d;&#x5b;&#x5d;&#x23;&#x3b;bar'>`)
	RegisterModFn("cssEscape", "ce", modCSSEscape).
		WithDescription("Applies CSS escaping to printing value.").
		WithExample(`background-image:url({%c= var1|escape('css') %}); // background-image:url(\3c \3e \27 \22 \26 \100 \2c \2e \5f aAzZ09\20 \21 \1f600 );`)
	RegisterModFn("jsEscape", "jse", modJSEscape).
		WithDescription("Applies Javascript escaping to printing value.").
		WithExample(`<script>{%J= var1 %}</script> // <script>\u003c\u003e\u0027\u0022\u0026\/,._aAzZ09\u0020\u0100\ud83d\ude00</script>`)
	RegisterModFn("raw", "noesc", func(_ *Ctx, _ *any, _ any, _ []any) error { return nil }).
		WithDescription("Disable value escaping/quoting inside bound tags (`{% jsonescape %}...{% endjsonescape %}`, ...)")

	// Register builtin round modifiers.
	RegisterModFn("round", "", modRound).
		WithDescription("Modifier `round` returns the nearest integer, rounding half away from zero.")
	RegisterModFn("roundPrec", "roundp", modRoundPrec).
		WithDescription("Modifier `roundPrec` rounds value to given precision.").
		WithExample(`// f = 3.1415
{%= f|roundPrec(3) %} // 3.141`)
	RegisterModFn("ceil", "", modCeil).
		WithDescription("Modifier `ceil` returns the least integer value greater than or equal to x.")
	RegisterModFn("ceilPrec", "ceilp", modCeilPrec).
		WithDescription("Modifier `ceilPrec` rounds value to ceil value with given precision").
		WithExample(`// f = 56.68734
{% = f|ceilPrec(3) %} // 56.688`)
	RegisterModFn("floor", "", modFloor).
		WithDescription("Modifier `floor` returns the greatest integer value less than or equal to x.")
	RegisterModFn("floorPrec", "floorp", modFloorPrec).
		WithDescription("Modifier `floorPrec` rounds value to floor value with given precision").
		WithExample(`// f = 20.214999
{% = f|floorPrec(3) %} // 20.214`)

	// Register time modifiers.
	RegisterModFnNS("time", "now", "", modNow).
		WithDescription("Returns the current local time.")
	RegisterModFnNS("time", "format", "date", modDate).
		WithParam("layout string", "See https://github.com/koykov/clock#format for possible patterns").
		WithDescription("Modifier `time::format` returns a textual representation of the time value formatted according given layout.").
		WithExample(`{%= date|time::date("%d %b %y %H:%M %z") %} // 05 Feb 09 07:00 +0200
{%= date|time::date("%b %e %H:%M:%S.%N") %} // Feb  5 07:00:57.012345600`)
	RegisterModFnNS("time", "add", "date_modify", modDateAdd).
		WithParam("duration string", "Textual representation of duration you want to add to the datetime. Possible units:\n"+
			"  * `nsec`, `ns`\n"+
			"  * `usec`, `us`, `µs`\n"+
			"  * `msec`, `ms`\n"+
			"  * `seconds`, `second`, `sec`, `s`\n"+
			"  * `minutes`, `minute`, `min`, `m`\n"+
			"  * `hours`, `hour`, `hr`, `h`\n"+
			"  * `days`, `day`, `d`\n"+
			"  * `weeks`, `week`, `w`\n"+
			"  * `months`, `month`, `M`\n"+
			"  * `years`, `year`, `y`\n"+
			"  * `century`, `cen`, `c`\n"+
			"  * `millennium`, `mil`\n").
		WithDescription("Modifier `time::add` returns time+duration.").
		WithExample(`{%= date|time::add("+1 m")|time::date(time::StampNano) %}{% endl %}      // Jan 21 20:05:26.000000555
{%= date|time::add("+1 min")|time::date(time::StampNano) %}{% endl %}	 // Jan 21 20:05:26.000000555
{%= date|time::add("+1 minute")|time::date(time::StampNano) %}{% endl %} // Jan 21 20:05:26.000000555
{%= date|time::add("+1 minutes")|time::date(time::StampNano) %}			 // Jan 21 20:05:26.000000555
`)

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

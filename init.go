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
		WithExample(`{%= notExistingOrEmptyVar|default("N/D") %}`)
	RegisterModFn("ifThen", "if", modIfThen).
		WithDescription("Modifier `ifThen` passes `arg` only if preceding condition is true.").
		WithParam("arg any", "").
		WithExample(`{%= obj.Active|ifThen("<button>Buy</button>") %}`)
	RegisterModFn("ifThenElse", "ifel", modIfThenElse).
		WithDescription("Modifier `ifTheElse` passes `arg0` if preceding condition is true or `arg1` otherwise.").
		WithParam("arg0 any", "").
		WithParam("arg1 any", "").
		WithExample(`{%= user.AllowSell|ifThenElse("<button>Sell</button>", "not available!") %}`)

	// Register fmt modifiers.
	RegisterModFnNS("fmt", "format", "f", modFmtFormat).
		WithDescription("Modifier `fmt::format` formats according to a format specifier and returns the resulting string.").
		WithParam("format string", "").
		WithParam("args ...any", "").
		WithExample("{%= fmt::format(\"Welcome %s\", user.Name) %}")

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
		WithExample(`<a href="{%l= link %}"> // <a href="https://x.com/link-with-\"-and+space-symbol">`)
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
{%= date|time::add("+1 minutes")|time::date(time::StampNano) %}		 // Jan 21 20:05:26.000000555
`)

	// Register math modifiers.
	RegisterModFnNS("math", "abs", "", modAbs).
		WithParam("val float", "Not mandatory").
		WithDescription("Returns absolute value.").
		WithExample(`num = -123456
{%= num|math::abs() %} // 123456
{%= math::abs(num) %}  // 123456`)
	RegisterModFnNS("math", "inc", "", modInc).
		WithParam("val float", "Not mandatory").
		WithDescription("Increments the value.").
		WithExample(`num = 10
{%= num|math::inc() %} // 11
{%= math::inc(num) %}  // 11`)
	RegisterModFnNS("math", "dec", "", modDec).
		WithParam("val float", "Not mandatory").
		WithDescription("Decrements the value.").
		WithExample(`num = 10
{%= num|math::dec() %} // 9
{%= math::dec(num) %}  // 9`)
	RegisterModFnNS("math", "add", "", modMathAdd).
		WithParam("val float", "Not mandatory").
		WithParam("arg float", "Value to add").
		WithDescription("Adds `arg` to value.").
		WithExample(`num = 555
{%= num|math::add(5.345) %}  // 560.345
{%= math::add(num, 5.345) %} // 560.345`)
	RegisterModFnNS("math", "sub", "", modMathSub).
		WithParam("val float", "Not mandatory").
		WithParam("arg float", "Value to subtract").
		WithDescription("Subtracts `arg` from value.").
		WithExample(`num = 555
{%= num|math::sub(5.345) %}  // 549.655
{%= math::sub(num, 5.345) %} // 549.655`)
	RegisterModFnNS("math", "mul", "", modMathMul).
		WithParam("val float", "Not mandatory").
		WithParam("arg float", "Value to multiply").
		WithDescription("Multiplies value with `arg`.").
		WithExample(`num = 10
{%= num|math::mul(5.345) %}  // 100
{%= math::mul(num, 5.345) %} // 100`)
	RegisterModFnNS("math", "div", "", modMathDiv).
		WithParam("val float", "Not mandatory").
		WithParam("arg float", "Value to divide").
		WithDescription("Divides value to `arg`.").
		WithExample(`num = 10
{%= num|math::div(5) %}  // 2
{%= math::div(num, 5) %} // 2`)
	RegisterModFnNS("math", "mod", "", modMathMod).
		WithParam("val float", "Not mandatory").
		WithParam("arg float", "Value to divide").
		WithDescription("Modifier 'mod' implement modulus operator. Returns the remainder or signed remainder of a division, after value is divided by `arg`.").
		WithExample(`num = 10
{%= num|math::mod(3) %}  // 1
{%= math::mod(num, 3) %} // 1`)
	RegisterModFnNS("math", "sqrt", "", modMathSqrt).
		WithParam("val float", "Not mandatory").
		WithDescription("Modifier 'sqrt' returns square root of value (`√val`).").
		WithExample(`num = 100
{%= num|math::sqrt() %} // 10
{%= math::sqrt(num) %}  // 10`)
	RegisterModFnNS("math", "cbrt", "", modMathCbrt).
		WithParam("val float", "Not mandatory").
		WithDescription("Modifier 'sqrt' computes the cube root of value (∛val).").
		WithExample(`num = 1000
{%= num|math::cqrt() %} // 10
{%= math::cqrt(num) %}  // 10`)
	RegisterModFnNS("math", "radical", "rad", modMathRadical).
		WithParam("val float", "Not mandatory").
		WithParam("root float", "Root of value required").
		WithDescription("Modifier 'radical' computes the radical (`root`-order root) of value.").
		WithExample(`num = 16
{%f.6= num|math::rad(4) %}  // 2
{%f.6= math::rad(num, 4) %} // 2`)
	RegisterModFnNS("math", "exp", "", modMathExp).
		WithParam("val float", "Not mandatory").
		WithDescription("Modifier 'exp' computes e**value, the base-e exponential of value.").
		WithExample(`num = 5
{%f.6= num|math::exp() %} // 148.4131591025766
{%f.6= math::exp(num) %}  // 148.4131591025766`)
	RegisterModFnNS("math", "log", "", modMathLog).
		WithParam("val float", "Not mandatory").
		WithDescription("Modifier 'log' computes the natural logarithm of value.").
		WithExample(`num = 6
{%f.6= num|math::log() %} // 1.791759469228055
{%f.6= math::log(num) %}  // 1.791759469228055`)
	RegisterModFnNS("math", "factorial", "fact", modMathFact).
		WithParam("val float", "Not mandatory").
		WithParam("root float", "Root of value required").
		WithDescription("Modifier 'radical' computes the radical (`root`-order root) of value.").
		WithExample(`num = 16
{%f.6= num|math::rad(4) %}  // 2
{%f.6= math::rad(num, 4) %} // 2`)
	RegisterModFnNS("math", "max", "", modMathMax).
		WithParam("arg0 float", "").
		WithParam("arg1 float", "").
		WithDescription("Modifier 'max' returns maximum value of `arg0` or `arg1`.")
	RegisterModFnNS("math", "min", "", modMathMin).
		WithParam("arg0 float", "").
		WithParam("arg1 float", "").
		WithDescription("Modifier 'max' returns minimum value of `arg0` or `arg1`.")
	RegisterModFnNS("math", "pow", "", modMathPow).
		WithParam("val float", "Not mandatory").
		WithParam("exp float", "Exponent").
		WithDescription("Modifier 'radical' computes `val**exp`, the base-`val` exponential of `exp`.").
		WithExample(`num = 2
{%= num|math::pow(4) %}  // 16
{%= math::pow(num, 4) %} // 16`)

	// Register builtin condition helpers.
	RegisterCondFn("lenEq0", condLenEq0).
		WithParam("arg bytes", "Possible types: `string`, `[]byte`, `[]string`, `[][]byte`.").
		WithDescription("Checks if length of `arg` equal to zero.").
		WithNote("DEPRECATED! Use native expression like `{% if len(arg) == 0 %}...{% endif %}`.")
	RegisterCondFn("lenGt0", condLenGt0).
		WithParam("arg bytes", "Possible types: `string`, `[]byte`, `[]string`, `[][]byte`.").
		WithDescription("Checks if length of `arg` greater or equal to zero.").
		WithNote("DEPRECATED! Use native expression like `{% if len(arg) >= 0 %}...{% endif %}`.")
	RegisterCondFn("lenGtq0", condLenGtq0).
		WithParam("arg bytes", "Possible types: `string`, `[]byte`, `[]string`, `[][]byte`.").
		WithDescription("Checks if length of `arg` less or equal to zero.").
		WithNote("DEPRECATED! Use native expression like `{% if len(arg) <= 0 %}...{% endif %}`.")

	// Register builtin empty check helpers.
	RegisterEmptyCheckFn("int", EmptyCheckInt).
		WithDescription("Checks if value is integer and contains zero value.")
	RegisterEmptyCheckFn("uint", EmptyCheckUint).
		WithDescription("Checks if value is unsigned integer and contains zero value.")
	RegisterEmptyCheckFn("float", EmptyCheckFloat).
		WithDescription("Checks if value is float and contains zero value.")
	RegisterEmptyCheckFn("bytes", EmptyCheckBytes).
		WithDescription("Checks if value is byte slice (`[]byte`) and its length equals to zero.")
	RegisterEmptyCheckFn("bytes_slice", EmptyCheckBytesSlice).
		WithDescription("Checks if value is bytes matrix (`[][]byte`) and contains no rows.")
	RegisterEmptyCheckFn("str", EmptyCheckStr).
		WithDescription("Checks if value is string and its length equals to zero.")
	RegisterEmptyCheckFn("str_slice", EmptyCheckStrSlice).
		WithDescription("Checks if value is slice of strings and its contains n string.")
	RegisterEmptyCheckFn("bool", EmptyCheckBool).
		WithDescription("Checks if value is boolean and equals to false.")

	// Register datetime layouts.
	RegisterGlobalNS("time", "Layout", "", clock.Layout).
		WithType("string").
		WithDescription("time.Layout presentation in strtotime format.")
	RegisterGlobalNS("time", "ANSIC", "", clock.ANSIC).
		WithType("string").
		WithDescription("time.ANSIC presentation in strtotime format.")
	RegisterGlobalNS("time", "UnixDate", "", clock.UnixDate).
		WithType("string").
		WithDescription("time.UnixDate presentation in strtotime format.")
	RegisterGlobalNS("time", "RubyDate", "", clock.RubyDate).
		WithType("string").
		WithDescription("time.RubyDate presentation in strtotime format.")
	RegisterGlobalNS("time", "RFC822", "", clock.RFC822).
		WithType("string").
		WithDescription("time.RFC822 presentation in strtotime format.")
	RegisterGlobalNS("time", "RFC822Z", "", clock.RFC822Z).
		WithType("string").
		WithDescription("time.RFC822Z presentation in strtotime format.")
	RegisterGlobalNS("time", "RFC850", "", clock.RFC850).
		WithType("string").
		WithDescription("time.RFC850 presentation in strtotime format.")
	RegisterGlobalNS("time", "RFC1123", "", clock.RFC1123).
		WithType("string").
		WithDescription("time.RFC1123 presentation in strtotime format.")
	RegisterGlobalNS("time", "RFC1123Z", "", clock.RFC1123Z).
		WithType("string").
		WithDescription("time.RFC1123Z presentation in strtotime format.")
	RegisterGlobalNS("time", "RFC3339", "", clock.RFC3339).
		WithType("string").
		WithDescription("time.RFC3339 presentation in strtotime format.")
	RegisterGlobalNS("time", "RFC3339Nano", "", clock.RFC3339Nano).
		WithType("string").
		WithDescription("time.RFC3339Nano presentation in strtotime format.")
	RegisterGlobalNS("time", "Kitchen", "", clock.Kitchen).
		WithType("string").
		WithDescription("time.Kitchen presentation in strtotime format.")
	RegisterGlobalNS("time", "Stamp", "", clock.Stamp).
		WithType("string").
		WithDescription("time.Stamp presentation in strtotime format.")
	RegisterGlobalNS("time", "StampMilli", "", clock.StampMilli).
		WithType("string").
		WithDescription("time.StampMilli presentation in strtotime format.")
	RegisterGlobalNS("time", "StampMicro", "", clock.StampMicro).
		WithType("string").
		WithDescription("time.StampMicro presentation in strtotime format.")
	RegisterGlobalNS("time", "StampNano", "", clock.StampNano).
		WithType("string").
		WithDescription("time.StampNano presentation in strtotime format.")

	// Register test modifiers.
	RegisterModFn("testNameOf", "", modTestNameOf).
		WithDescription("Testing stuff: don't use in production.")
	RegisterModFnNS("testns", "pack", "", func(_ *Ctx, _ *any, _ any, _ []any) error { return nil }).
		WithDescription("Testing namespace stuff: don't use in production.")
	RegisterModFnNS("testns", "extract", "", func(_ *Ctx, _ *any, _ any, _ []any) error { return nil }).
		WithDescription("Testing namespace stuff: don't use in production.")
	RegisterModFnNS("testns", "marshal", "", func(_ *Ctx, _ *any, _ any, _ []any) error { return nil }).
		WithDescription("Testing namespace stuff: don't use in production.")
	RegisterModFnNS("testns", "modCB", "", func(ctx *Ctx, _ *any, _ any, args []any) error {
		ctx.SetStatic("testVar", args[0])
		return nil
	}).
		WithDescription("Testing namespace stuff: don't use in production.")

	// Register test condition-ok helpers.
	RegisterCondOKFn("__testUserNextHistory999", testCondOK).
		WithDescription("Testing namespace stuff: don't use in production.")

	// Register test variable-inspector pairs.
	RegisterVarInsPair("__testFin999", &testobj_ins.TestFinanceInspector{}).
		WithDescription("Testing stuff: don't use in production.")
}

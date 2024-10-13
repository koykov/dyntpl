# Dynamic templates

Dynamic replacement for [quicktemplate](https://github.com/valyala/quicktemplate) template engine.

## Retrospective

We're used for a long time quicktemplate for building JSON to interchange data between microservices on high-load project,
and we were happy. Ерут we need to be able to change existing templates or add new templates on the fly. Unfortunately
quicktemplates doesn't support this and this package was developed as replacement.

It reproduces many of qtpl features and syntax.

## How it works

Working in templates divided to two phases - parsing and templating. The parsing phase builds from template a tree
(like AST) and registers it in templates registry by unique name afterward. This phase doesn't intended to use in highload
conditions due to high pressure to cpu/mem. The second phase - templating, against intended to use in highload.

Templating phase required a preparation to pass data to te template. There is special object [Ctx](ctx.go), that collects
variables to use in template. Each variable must have three params:
* unique name
* data - anything you need to use in template
* inspector type

Inspector type must be explained.

The biggest problem during development was how to get data from arbitrary structure without using reflection,
since `reflect` package produces a lot of allocations by design and is extremely slowly in general.

To solve that problem was developed code-generation framework [inspector](https://github.com/koykov/inspector) framework.
It provides primitive methods to read data from struct's fields, iterating, etc. without using reflection and makes it
fast. Framework generates for each required struct a special type with such methods. It isn't a pure dynamic like 
`reflect` provided, but works so [fast](https://github.com/koykov/versus/tree/master/inspector2) and makes zero allocations. 

You may check example of inspectors in package [testobj_ins](./testobj_ins) that represents testing structures
in [testobj](./testobj).

## Usage

The typical usage of dyntpl looks like this:
```go
package main

import (
	"bytes"

	"github.com/koykov/dyntpl"
	"path/to/inspector_lib_ins"
	"path/to/test_struct"
)

var (
	// Fill up test struct with data.
	data = &test_struct.Data{
		// ...
	}

	// Template code.
	tplData = []byte(`{"id":"{%=data.Id%}","hist":[{%for _,v:=range data.History separator ,%}"{%=v.Datetime%}"{%endfor%}]}`)
)

func init() {
	// Parse the template and register it.
	tree, _ := dyntpl.Parse(tplData, false)
	dyntpl.RegisterTplKey("tplData", tree)
}

func main() {
	// Prepare output buffer
	buf := bytes.Buffer{}
	// Prepare dyntpl context.
	ctx := dyntpl.AcquireCtx()
	ctx.Set("data", data, inspector_lib_ins.DataInspector{})
	// Execute the template and write result to buf.
	_ = dyntpl.Write(&buf, "tplData", ctx)
	// Release context.
	dyntpl.ReleaseCtx(ctx)
	// buf.Bytes() or buf.String() contains the result.
}
```
Content of `init()` function should be executed once (or periodically on fly from some source, eg DB).

Content of `main()` function is how to use dyntpl in general way in highload. Of course, byte buffer should take from the pool.

## Benchmarks

See [bench.md](bench.md) for result of internal benchmarks.

Highly recommend to check `*_test.go` files in the project, since them contains a lot of typical language constructions
that supports this engine.

See [versus/dyntpl](https://github.com/koykov/versus/tree/master/dyntpl) for comparison benchmarks with
[quicktemplate](https://github.com/valyala/quicktemplate) and native marshaler/template.

As you can see, dyntpl in ~3-4 times slower than [quicktemplates](https://github.com/valyala/quicktemplate).
That is a cost for dynamics. There is no way to write template engine that will faster than native Go code.

## Syntax

### Print

The most general syntax construction is printing a variable or structure field:
```
This is a simple statis variable: {%= var0 %}
This is a field of struct: {%= obj.Parent.Name %}
```
Construction `{%= ... %}` prints data as is without any type check, escaping or modification.

There are special escaping modifiers. They should use before `=`:
* `h` - HTML escape.
* `a` - HTML attribute escape.
* `j` - JSON escape.
* `q` - JSON quote.
* `J` - JS escape.
* `u` - URL encode.
* `l` - Link escape.
* `c` - CSS escape.
* `f.<num>` - float with precision, example: `{%f.3= 3.1415 %}` will output `3.141`.
* `F.<num>` - ceil rounded float with precision, example: `{%F.3= 3.1415 %}` will output `3.142`.

Note, that none of these directives doesn't apply by default. It's your responsibility to controls what and where you print.

All directives (except of `f` and `F`) supports multipliers, like `{%jj= ... %}`, `{%uu= ... %}`, `{%uuu= ... %}`, ...

For example, the following instruction `{%uu= someUrl %}` will print double url-encoded value of `someUrl`. It may be
helpful to build chain of redirects:
```
https://domain.com?redirect0={%u= url0 %}{%uu= url1 %}{%uuu= url2 %}
```

Also, you may combine directives in any combinations (`{%Ja= var1 %}`, `{%jachh= var1 %}`, ...). Modifier will apply
consecutive and each modifier will take to input the result of previous modifier.

### Bound tags

To apply escaping to some big block containing both text and printing variables there are special bound tags:
* `{% jsonquote %}...{% endjsonquote %}` apply JSON escape to all contents.
* `{% htmlescape %}...{% endhtmlescape %}` apply HTML escape.
* `{% urlencode %}...{% endurlencode %}` URL encode all text data.

Example:
```
{"key": "{% jsonquote %}Lorem ipsum "dolor sit amet", {%= var0 %}.{%endjsonquote%}"}
```

#### Prefix/suffix

Print construction supports prefix and suffix attributes, it may be handy when you print HTML or XML:
```html
<ul>
{%= var prefix <li> suffix </li> %}
</ul>
```
Prefix and suffix will print only if `var` isn't empty. Prefix/suffix has shorthands `pfx` and `sfx`.

Also print supports data modifiers. They calls typically for any template languages:
```
Name: {%= obj.Name|default("anonymous") %}
```
and may contains variadic list of arguments or doesn't contain them at all. See the full list of built-in modifiers
in [init.go](init.go) (calls of `RegisterModFn()`).
You may register your own modifiers, see section [Modifier helpers](#modifier-helpers).

#### Conditions

Conditions in dyntpl is pretty simple and supports only two types of record:
* `{% if leftVar [=|!=|>|>=|<|<=] rightVar %}...{% endif %}`
* `{% if conditionHelper(var0, obj.Name, "foo") %}...{% endif %}`

First type is for the simplest case, like:
```html
{% if user.Id == 0 %}
You should <a href="#">log in</a>.
{% endif %}
```
Left side or right side or both may be a variable. But you can't specify a condition with static values on both sides,
since it's senseless.

Second type of condition is for more complex conditions when any side of condition should contain Go code, like:
```
Welcome, {% if len(user.Name) > 0 %}{%= user.Name %}{% else %}anonymous{%endif%}!
```
Dyntpl can't handle that kind of records, but it supports special functions that may make a decision is given args
suitable or not and return true/false.
See the full list of built-in condition helpers in [init.go](init.go) (calls of `RegisterCondFn`). Of course,
you can register your own handlers to implement your logic.

For multiple conditions you can use `switch` statement, example 1:
```xml
<item type="{% switch item.Type %}
{% case 0 %}
    deny
{% case 1 %}
    allow
{% case 2 %}
    allow-by-permission
{% default %}
    unknown
{% endswitch %}">foo</item>
```
, example 2:
```xml
<item type="{% switch %}
{% case item.Type == 0 %}
    deny
{% case item.Type == 1 %}
    allow
{% case item.Type == 2 %}
    allow-by-permission
{% default %}
    unknown
{% endswitch %}">foo</item>
```

Switch can handle only primitive cases, condition helpers doesn't support.

#### Loops

Dyntpl supports both types of loops:
* conditional loop from three components separated by semicolon, like `{% for i:=0; i<5; i++ %}...{% endfor %}`
* range-loop, like `{% for k, v := range obj.Items %}...{% endfor %}`

Edge cases like `for k < 2000 {...}` or `for ; i < 10 ; {...}` isn't supported.
Also, you can't make infinite loop by using `for {...}`.

There is a special attribute `separator` that made special to build JSON output. Example of use:
```
[
  {% for _, a := range user.History separator , %}
    {
      "id": {%q= a.Id %},
      "date": {%q= a.Date %},
      "comment": {%q= a.Note %}
    }
  {% endfor %}
]
```
The output that will produced:
```json
[
  {"id":1, "date": "2020-01-01", "comment": "success"},
  {"id":2, "date": "2020-01-01", "comment": "failed"},
  {"id":3, "date": "2020-01-01", "comment": "rejected"}
]
```
As you see, commas between 2nd and last elements was added by dyntpl without any additional handling like
`...{% if i>0 %},{% endif %}{% endfor %}`.
Separator has shorthand variant `sep`.

## Include sub-templates

Just call `{% include subTplID %}` (example `{% include sidebar/right %}`) to render and include output of that template
inside current template.
Sub-template will use parent template's context to access the data.

## Modifier helpers

Modifiers is a special functions that may perform modifications over the data during print. These function have signature:
```go
func(ctx *Ctx, buf *any, val any, args []any) error
```
and should be registered using function `dyntpl.RegisterModFn()`. See [init.go](init.go) for examples.
See [mod.go](mod.go) for explanation of arguments.

Modifiers calls using pipeline symbol after a variable, example: `{%= var0|default(0) %}`.

You may specify a sequence of modifiers: `{%= var0|roundPrec(4)|default(1) %}`.

## Condition helpers

If you want to make a condition more complex than simple condition, you may declare a special function with signature:
```go
func(ctx *Ctx, args []any) bool
```
and register it using function `dyntpl.RegisterCondFn()`. See [init.go](init.go) for examples.
See [cond.go](cond.go) for explanation of arguments.

After declaring and registering you can use the helper in conditions:
```
{% if <condFnName>(var0, var1, "static val", 0, 15.234) %}...{% endif %}
```
Function will make a decision according arguments you take and will return true or false.

## Bound tags

Dyntpl support special tags to escape/quote the output. Currently, allows three types:
* `{% jsonquote %}...{% endjsonquote %}` apply JSON escape for all text data.
* `{% htmlescape %}...{% endhtmlescape %}` apply HTML escape.
* `{% urlencode %}...{% endurlencode %}` URL encode all text data.

Note, these tags escapes only text data inside. All variables should be escaped using corresponding modifiers. Example:
```json
{"key": "{% jsonquote %}Lorem ipsum "dolor sit amet", {%j= var0 %}.{%endjsonquote%}"}
```
Here, `{% end/jsonquote %}` applies only for text data `Lorem ipsum "dolor sit amet",`, whereas `var0` prints using
JSON-escape printing prefix.

`{% end/htmlescape %}` and `{% end/urlencode %}` works the same.

## Extensions

Dyntpl's features may be extended by including modules to the project. Currently supported modules:
* [dyntpl_vector](https://github.com/koykov/dyntpl_vector) provide support of vector parsers inside the templates.
* [dyntpl_i18n](https://github.com/koykov/dyntpl_i18n) provide support of i18n features.

To enable necessary module just import it to the project, eg:
```go
import (
	_ "https://github.com/koykov/dyntpl_vector"
)
```
and vector's [features](https://github.com/koykov/dyntpl_vector) will be available inside templates. 

Feel free to develop your own extensions. Strongly recommend to register new modifiers using namespaces, like
[this](https://github.com/koykov/dyntpl_vector/blob/master/init.go#L12).

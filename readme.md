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

### Print modifiers

Alongside of short modifiers using before `=` engine supports user defined modifiers. You may use them after printing
variable using `|` and looks lice function call:
```
Name: {%= obj.Name|default("anonymous") %}Welcome, {%= testNameOf(user, {"foo": "bar", "id": user.Id}, "qwe") %}
                  ^ simple example                     ^ call modifier without variable like simple function call
Chain of modifiers: {%= dateVariable|default("2022-10-04")|formatDate("%y-%m-%d") %}
                                    ^ first modifier      ^ second modifier
```
Modifiers may collect in chain with variadic length. In that case each modifier will take to input the result of
previous modifier. Each modifier may take arbitrary count of arguments.

In general modifier is a Go function with special signature:
```go
type ModFn func(ctx *Ctx, buf *any, val any, args []any) error
```
, where:
* ctx - context of the template
* buf - pointer to return value
* val - value to pass to input (eg `varName|modifier()` value of `varName`)
* args - list of all aguments

After writing your function you need to register it using one of functions:
* `RegisterModFn(name, alias string, mod ModFn)`
* `RegisterModFnNS(namespace, name, alias string, mod ModFn)`

They are the same, but NS version allows to specify the namespace of the function. In that case you should specify namespace
on modifiers call:
```
Print using ns: {%= varName|namespaceName::modifier() %}
```

#### Conditions

dyntpl supports classic syntax of conditions:
```
{% if leftVar [==|!=|>|>=|<|<=] rightVar %}
    true branch
{% else %}
    false branch
{% endif %}
```

Examples: [1](testdata/parser/condition.tpl), [2](testdata/parser/conditionNested.tpl), [3](testdata/parser/conditionStr.tpl).

dyntpl can't handle complicated conditions containing more than one comparison, like:
```
{% if user.Id == 0 || user.Finance.Balance == 0 %}You're not able to buy!{% endif %}
```
In the grim darkness of the far future this problem will be solved, but now you can make nested conditions or use
conditions helpers - functions with signature:
```go
type CondFn func(ctx *Ctx, args []any) bool
```
, where you may pass arbitrary amount of arguments and these functions will return bool to choose right execution branch.
These function is user-defined like modifiers and you may write your own and then register it using one of functions:
```go
func RegisterCondFn(name string, cond CondFn)
func RegisterCondFnNS(namespace, name string, cond CondFn) // namespace version
```

Then condition helper will accessible inside templates and you may use it using name:
```
{% if helperName(user.Id, user.Finance.Balance) %}You're not able to buy!{% endif %}
```

As exception there are two functions `len()` and `cap()` that works the same as builtin native Go functions. The result
of their execution may be compared
```
{% if len(user.Name) > 0 %}...{% endif %}
```
, whereas user-defined helpers doesn't allow comparisons.

For multiple conditions you can use `switch` statement, examples:
* [classic switch](testdata/parser/switch.tpl)
* [no-condition switch](testdata/parser/switchNoCondition.tpl)
* [no-condition switch with helpers](testdata/parser/switchNoConditionWithHelper.tpl)

### Loops

Dyntpl supports both types of loops:
* counter loops, like `{% for i:=0; i<5; i++ %}...{% endfor %}`
* range-loop, like `{% for k, v := range obj.Items %}...{% endfor %}`

Edge cases like `for k < 2000 {...}` or `for ; i < 10 ; {...}` isn't supported.
Also, you can't make infinite loop by using `for {...}`.

#### Separators

When separator between iterations required, there is a special attribute `separator` that made special to build JSON
output. Example of use:
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
The output that will be produced:
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

#### loop-else

Separator isn't the last exclusive feature of loops. For loops allows `else` branch like for conditions:
```
<select name="type">
  {% for k, v := range user.historyTags %}
    <option value="{%= k %}">{%= v %}</option>
  {% else %}
    <option>N/D</option>
  {% endfor %}
</select>
```
If loop's source is empty and there aren't data to iterate, then else branch will execute without manual handling. In the
example above, if `user.historyTags` is empty, the empty `<option>` will display.

#### Loop breaking

dyntpl supports default instructions `break` and `continue` to break loop/iteration, example:
```
{% for _, v := list %}
  {% if v.ID == 0 %}
    {% continue %}
  {% endif %}
  {% if v.Status == -1 %}
    {% break %}
  {% endif %}
{% endfor %}
```

These instructions works as intended, but they required condition wrapper and that's bulky. Therefore, dyntpl provide
combined `break if` and `continue if` that works the same:
```
{% for _, v := list %}
  {% continue if v.ID == 0 %}
  {% break if v.Status == -1 %}
{% endfor %}
```

The both examples is equal, but the second is more compact.

#### Lazy breaks

Imagine the case - you've decided in te middle of iteration that loop require break, but the iteration must finish its
work the end. Eg template printing some XML element and break inside it will produce unclosed tag:
```
<?xml version="1.0" encoding="UTF-8"?>
<users>
  {% for _, u := range users %}
    <user>
        <name>{%= u.Name %}</name>
        {% if u.Blocked == 1 %}
          {% break %} {# <-- unclosed tag reson #}
        {% endif %}
        <balance>{%= u.Balance }</balance>
    </user>
  {% endfor %}
</users>
```

Obviously, invalid XML document will build. For that case dyntpl supports special instruction `lazybreak`. It breaks the
loop but allows current iteration works till the end.

#### Nested loops break

In native Go to break nested loops you must use one of instructions:
* `goto <label>`
* `break <label>`
* `continue <label>`

Well, supporting of these things is too strange for template engine, isn't it? There is more handy [break, provided by php](https://www.php.net/manual/en/control-structures.break.php).
It allows to specify after keyword how many loops must be touched by instruction. dyntpl implements this case and provides
instructions:
* `break N`
* `lazybreak N`

```
{% for i:=0; i<10; i++ %}
  bar
  {% for j:=0; i<10; i++ %}
    foo
    {% if j == 8 %}
      {% break 2 %}
    {% endif %}
    {% if j == 7 %}
      {% lazybreak 2 %}
    {% endif %}
    {%= j %}
  {% endfor %}
  {%= i %}
{% endfor %}
```

`break/lazybreak N` instructions supports conditional versions:
* `break N if`
* `lazybreak N if`

The example above may be changed to:
```
    ...
    {% break 2 if j == 8 %}
    {% lazybreak 2 if j == 7 %}
    ...
```
and you will give the same output.

## Include sub-templates

To reuse templates exists instruction `include` that may be included directly from the template.
Just call `{% include subTplID %}` (example `{% include sidebar/right %}`) to render and include output of that template
inside current template.
Sub-template will use parent template's context to access the data.

Also, you may include sub-templates in bash-style using `.`.

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

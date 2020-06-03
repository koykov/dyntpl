# Dynamic templates

Dynamic replacement for [quicktemplate](https://github.com/valyala/quicktemplate) template engine.

## Retrospective

We're used for a long time quicktemplate for building JSON to interchange data between microservices on high-load project and we were happy.
But now we need to able change existing templates or add new templates on the fly. Unfortunately quicktemplates doesn't support this and this package was developed as replacement.

It reproduces many of qtpl features and syntax.

## How it works

The biggest problem during development was how to get data from arbitrary structure without using reflection, since `reflect` package produces a lot of allocations by design and is extremely slowly in general.

To solve that problem was developed [inspector](https://github.com/koykov/inspector) framework. It takes as argument path to the package with structures signatures and build an exact primitive methods to get data of any fields in them, loop over fields that support loops, etc...

You may check example of inspectors in subdirectory [testobj_ins](./testobj_ins) that represents testing structures in [testobj](./testobj).

## Usage

The typical usage of dyntpl looks like this:
```go
package main

import (
    "bytes"

    "github.com/koykov/dyntpl"
    "path/to/inspector_lib_ins"
)

var (
    // Test data.
    data = &Data{
        // ...
    }

    // Template code.
    tplData = []byte(`...`)
)

func init() {
    // Parse the template and register it.
    tree, _ := dyntpl.Parse(tplData, false)
    dyntpl.RegisterTpl("tplData", tree)
}

func main() {
    // Prepare output buffer
    buf := &bytes.Buffer{}
    // Prepare dyntpl context.
    ctx := dyntpl.AcquireCtx()
    ctx.Set("data", data, &inspector_lib_ins.DataInspector{})
    // Execute the template and write result to buf.
    _ = dyntpl.RenderTo(buf, "tplData", ctx)
    // Use result as buf.Bytes() or buf.String() ...
    // Release context.
    dyntpl.ReleaseCtx(ctx)
}
```
Content of `init()` function may be moved to scheduler and periodically take fresh template code from source data, e.g. DB table and update it on the fly.

Content of `main()` function is how to use dyntpl in general way. Of course, byte buffer should take from the pool.

## Benchmarks

Here is a result of internal benchmarks:
```
BenchmarkCtxGet-8                 20000000       110 ns/op       0 B/op       0 allocs/op
BenchmarkCtxPoolGet-8             10000000       152 ns/op       0 B/op       0 allocs/op
BenchmarkTplSimple-8               1000000      1102 ns/op       0 B/op       0 allocs/op
BenchmarkTplCond-8                 2000000       805 ns/op       0 B/op       0 allocs/op
BenchmarkTplCondNoStatic-8         2000000       927 ns/op       0 B/op       0 allocs/op
BenchmarkTplCondHlp-8              5000000       284 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitch-8               2000000       788 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitchNoCond-8         2000000       630 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopRange-8             500000      3777 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountStatic-8       300000      4731 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountBreak-8        300000      5159 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountContinue-8     200000      6426 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCount-8             300000      5784 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountCtx-8          200000      5450 ns/op       0 B/op       0 allocs/op
BenchmarkTplExit-8                 5000000       288 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDef-8               5000000       363 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDefStatic-8         3000000       486 ns/op       0 B/op       0 allocs/op
BenchmarkTplModJsonQuote-8         3000000       370 ns/op       0 B/op       0 allocs/op
BenchmarkTplModHtmlEscape-8        2000000       847 ns/op       0 B/op       0 allocs/op
BenchmarkTplModIfThen-8           10000000       209 ns/op       0 B/op       0 allocs/op
BenchmarkTplModIfThenElse-8        3000000       372 ns/op       0 B/op       0 allocs/op
```
Highly recommend to check `*_test.go` files in the project, since them contains a lot of typical language constructions that supports this engine.

`cmp_test` dir contains comparison tests with corresponding [quicktemplate's test](https://github.com/valyala/quicktemplate/tree/master/tests):
* for [templates timings](https://github.com/valyala/quicktemplate/blob/master/tests/templates_timing_test.go)
```
BenchmarkDyntpl1-8                 	 5000000	       328 ns/op	       0 B/op	       0 allocs/op
BenchmarkDyntpl10-8                	 1000000	      1427 ns/op	       0 B/op	       0 allocs/op
BenchmarkDyntpl100-8               	  100000	     14454 ns/op	       0 B/op	       0 allocs/op
``` 
* for [marshal timings](https://github.com/valyala/quicktemplate/blob/master/tests/marshal_timing_test.go)
```
BenchmarkMarshalJSONDyntpl1-8      	 3000000	       433 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalJSONDyntpl10-8     	 1000000	      2226 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalJSONDyntpl100-8    	  100000	     20331 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalJSONDyntpl1000-8   	   10000	    203189 ns/op	       8 B/op	       0 allocs/op
BenchmarkMarshalXMLDyntpl1-8       	 3000000	       444 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalXMLDyntpl10-8      	 1000000	      2182 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalXMLDyntpl100-8     	  100000	     19500 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalXMLDyntpl1000-8    	   10000	    196965 ns/op	       9 B/op	       0 allocs/op
```
As you can see, dyntpl in ~3-4 times slowest than quicktemplates. That is a cost for dynamics. There is no way to write template engine that will fastests than native Go code.

## Syntax

...

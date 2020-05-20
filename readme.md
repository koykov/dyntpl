# Dynamic templates

## Benchmarks
```
BenchmarkCtx_Get-8                20000000       119 ns/op       0 B/op       0 allocs/op
BenchmarkCtxPool_Get-8            10000000       142 ns/op       0 B/op       0 allocs/op
BenchmarkTplSimple-8               1000000      1241 ns/op       0 B/op       0 allocs/op
BenchmarkTplCond-8                 2000000       743 ns/op       0 B/op       0 allocs/op
BenchmarkTplCondNoStatic-8         1000000      1037 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitch-8               2000000       707 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitchNoCond-8         3000000       609 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopRange-8             300000      4606 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountStatic-8       300000      5113 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountBreak-8        200000      5617 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountContinue-8     200000      7824 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCount-8             300000      5455 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountCtx-8          200000      6074 ns/op       0 B/op       0 allocs/op
BenchmarkTplExit-8                 5000000       252 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDef-8               5000000       331 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDefStatic-8         3000000       462 ns/op       0 B/op       0 allocs/op
BenchmarkTplModJsonQuote-8         5000000       347 ns/op       0 B/op       0 allocs/op
BenchmarkTplModHtmlEscape-8        1000000      1003 ns/op       0 B/op       0 allocs/op

```

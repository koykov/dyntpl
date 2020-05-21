# Dynamic templates

## Benchmarks
```
BenchmarkCtxGet-8                 10000000       121 ns/op       0 B/op       0 allocs/op
BenchmarkCtxPoolGet-8             10000000       154 ns/op       0 B/op       0 allocs/op
BenchmarkTplSimple-8               1000000      1275 ns/op       0 B/op       0 allocs/op
BenchmarkTplCond-8                 2000000       888 ns/op       0 B/op       0 allocs/op
BenchmarkTplCondNoStatic-8         1000000      1023 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitch-8               2000000       805 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitchNoCond-8         2000000       659 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopRange-8             300000      4831 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountStatic-8       200000      5717 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountBreak-8        200000      6462 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountContinue-8     200000      7901 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCount-8             300000      5893 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountCtx-8          200000      6536 ns/op       0 B/op       0 allocs/op
BenchmarkTplExit-8                 5000000       294 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDef-8               5000000       352 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDefStatic-8         3000000       473 ns/op       0 B/op       0 allocs/op
BenchmarkTplModJsonQuote-8         5000000       381 ns/op       0 B/op       0 allocs/op
BenchmarkTplModHtmlEscape-8        1000000      1037 ns/op       0 B/op       0 allocs/op
BenchmarkTplModIfThen-8           10000000       220 ns/op       0 B/op       0 allocs/op
BenchmarkTplModIfThenElse-8        3000000       409 ns/op       0 B/op       0 allocs/op
```

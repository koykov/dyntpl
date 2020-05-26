# Dynamic templates

## Benchmarks
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

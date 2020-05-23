# Dynamic templates

## Benchmarks
```
BenchmarkCtxGet-8                 10000000       121 ns/op       0 B/op       0 allocs/op
BenchmarkCtxPoolGet-8             10000000       164 ns/op       0 B/op       0 allocs/op
BenchmarkTplSimple-8               1000000      1274 ns/op       0 B/op       0 allocs/op
BenchmarkTplCond-8                 2000000       906 ns/op       0 B/op       0 allocs/op
BenchmarkTplCondNoStatic-8         1000000      1046 ns/op       0 B/op       0 allocs/op
BenchmarkTplCondHlp-8              5000000       330 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitch-8               2000000       809 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitchNoCond-8         2000000       671 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopRange-8             300000      4754 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountStatic-8       200000      5882 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountBreak-8        200000      6262 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountContinue-8     200000      7970 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCount-8             200000      5954 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountCtx-8          200000      6258 ns/op       0 B/op       0 allocs/op
BenchmarkTplExit-8                 5000000       297 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDef-8               5000000       370 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDefStatic-8         3000000       483 ns/op       0 B/op       0 allocs/op
BenchmarkTplModJsonQuote-8         3000000       402 ns/op       0 B/op       0 allocs/op
BenchmarkTplModHtmlEscape-8        1000000      1047 ns/op       0 B/op       0 allocs/op
BenchmarkTplModIfThen-8           10000000       231 ns/op       0 B/op       0 allocs/op
BenchmarkTplModIfThenElse-8        3000000       423 ns/op       0 B/op       0 allocs/op
```

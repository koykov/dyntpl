# Dynamic templates

## Benchmarks
```
BenchmarkCtxGet-8                 20000000       109 ns/op       0 B/op       0 allocs/op
BenchmarkCtxPoolGet-8             10000000       148 ns/op       0 B/op       0 allocs/op
BenchmarkTplSimple-8               1000000      1080 ns/op       0 B/op       0 allocs/op
BenchmarkTplCond-8                 2000000       882 ns/op       0 B/op       0 allocs/op
BenchmarkTplCondNoStatic-8         2000000      1023 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitch-8               2000000       693 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitchNoCond-8         3000000       541 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopRange-8             300000      4529 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountStatic-8       300000      4781 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountBreak-8        300000      5173 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountContinue-8     200000      7374 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCount-8             200000      5858 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountCtx-8          200000      5913 ns/op       0 B/op       0 allocs/op
BenchmarkTplExit-8                 5000000       286 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDef-8               5000000       302 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDefStatic-8         3000000       460 ns/op       0 B/op       0 allocs/op
BenchmarkTplModJsonQuote-8         5000000       312 ns/op       0 B/op       0 allocs/op
BenchmarkTplModHtmlEscape-8        1000000      1033 ns/op       0 B/op       0 allocs/op
```

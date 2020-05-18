# Dynamic templates

## Benchmarks
```
BenchmarkTplSimple-8               1000000      1049 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDef-8               5000000       269 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDefStatic-8         3000000       460 ns/op       0 B/op       0 allocs/op
BenchmarkTplCond-8                 2000000       813 ns/op       0 B/op       0 allocs/op
BenchmarkTplCondNoStatic-8         1000000      1041 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitch-8               2000000       695 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitchNoCond-8         2000000       639 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopRange-8             300000      4043 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountStatic-8       300000      4597 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountBreak-8        300000      6215 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountContinue-8     200000      7105 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCount-8             300000      5849 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountCtx-8          200000      5222 ns/op       0 B/op       0 allocs/op
BenchmarkTplExit-8                 5000000       301 ns/op       0 B/op       0 allocs/op
BenchmarkCtx_Get-8                10000000       101 ns/op       0 B/op       0 allocs/op
BenchmarkCtxPool_Get-8            10000000       149 ns/op       0 B/op       0 allocs/op
```

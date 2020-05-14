# CbyteTpl

## Benchmarks
```
BenchmarkTplSimple-8             1000000      1252 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDef-8             5000000       338 ns/op       0 B/op       0 allocs/op
BenchmarkTplModDefStatic-8       3000000       461 ns/op       0 B/op       0 allocs/op
BenchmarkTplCond-8               2000000       893 ns/op       0 B/op       0 allocs/op
BenchmarkTplCondNoStatic-8       1000000      1021 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitch-8             2000000       762 ns/op       0 B/op       0 allocs/op
BenchmarkTplSwitchNoCond-8       2000000       650 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopRange-8           300000      4506 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountStatic-8     300000      5707 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCount-8           300000      5862 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountCtx-8        200000      6080 ns/op       0 B/op       0 allocs/op
BenchmarkCtx_Get-8              10000000       120 ns/op       0 B/op       0 allocs/op
BenchmarkCtxPool_Get-8          10000000       159 ns/op       0 B/op       0 allocs/op
```

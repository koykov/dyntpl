# CbyteTpl

## Benchmarks
```
BenchmarkTplSimple-8             2000000       945 ns/op       0 B/op       0 allocs/op
BenchmarkTplCond-8               2000000       669 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopRange-8           500000      3383 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCountStatic-8     300000      4523 ns/op       0 B/op       0 allocs/op
BenchmarkTplLoopCount-8           300000      4968 ns/op       0 B/op       0 allocs/op
BenchmarkCtx_Get-8              20000000       122 ns/op       0 B/op       0 allocs/op
BenchmarkCtxPool_Get-8          10000000       152 ns/op       0 B/op       0 allocs/op
```

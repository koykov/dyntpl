# CbyteTpl

## Benchmarks
```
BenchmarkTplSimple-8                 1000000          1045 ns/op           0 B/op          0 allocs/op
BenchmarkTplCond-8                   2000000           851 ns/op           0 B/op          0 allocs/op
BenchmarkTplCondNoStatic-8           2000000           874 ns/op           0 B/op          0 allocs/op
BenchmarkTplSwitch-8                 2000000           668 ns/op           0 B/op          0 allocs/op
BenchmarkTplSwitchNoCond-8           3000000           562 ns/op           0 B/op          0 allocs/op
BenchmarkTplLoopRange-8               300000          3593 ns/op           0 B/op          0 allocs/op
BenchmarkTplLoopCountStatic-8         300000          4843 ns/op           0 B/op          0 allocs/op
BenchmarkTplLoopCount-8               300000          4896 ns/op           0 B/op          0 allocs/op
BenchmarkTplLoopCountCtx-8            200000          5163 ns/op           1 B/op          0 allocs/op
BenchmarkCtx_Get-8                  10000000           125 ns/op           0 B/op          0 allocs/op
BenchmarkCtxPool_Get-8              10000000           140 ns/op           0 B/op          0 allocs/op
```

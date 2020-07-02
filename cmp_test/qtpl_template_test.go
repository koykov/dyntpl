package cmp

import (
	"fmt"
	"testing"

	"github.com/koykov/cbytebuf"

	"github.com/koykov/dyntpl"
	"github.com/koykov/dyntpl/testobj"
	"github.com/koykov/dyntpl/testobj_ins"
)

var (
	tplTemplate = []byte(`<html>
	<head><title>test</title></head>
	<body>
		<ul>
		{% for _, row := range bench.Rows %}
			{% if row.Print == true %}
				<li>ID={%= row.ID %}, Message={%h= row.Message %}</li>
			{% endif %}
		{% endfor %}
		</ul>
	</body>
</html>`)
)

func BenchmarkDyntpl1(b *testing.B) {
	benchmarkDyntpl(b, 1)
}

func BenchmarkDyntpl10(b *testing.B) {
	benchmarkDyntpl(b, 10)
}

func BenchmarkDyntpl100(b *testing.B) {
	benchmarkDyntpl(b, 100)
}

func benchmarkDyntpl(b *testing.B, rowsCount int) {
	tree, _ := dyntpl.Parse(tplTemplate, false)
	dyntpl.RegisterTpl("tplTemplate", tree)

	b.ResetTimer()
	b.ReportAllocs()

	bench := getBenchRows(rowsCount)
	b.RunParallel(func(pb *testing.PB) {
		buf := cbytebuf.LBAcquire()
		ctx := dyntpl.AcquireCtx()
		ctx.Set("bench", bench, &testobj_ins.BenchRowsInspector{})
		for pb.Next() {
			_ = dyntpl.RenderTo(buf, "tplTemplate", ctx)
		}
		dyntpl.ReleaseCtx(ctx)
		cbytebuf.LBRelease(buf)
	})
}

func getBenchRows(n int) *testobj.BenchRows {
	bench := &testobj.BenchRows{
		Rows: make([]testobj.BenchRow, n),
	}
	for i := 0; i < n; i++ {
		bench.Rows[i] = testobj.BenchRow{
			ID:      i,
			Message: fmt.Sprintf("message %d", i),
			Print:   (i & 1) == 0,
		}
	}
	return bench
}

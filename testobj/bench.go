package testobj

type BenchRow struct {
	ID      int
	Message string
	Print   bool
}

type BenchRows struct {
	Rows []BenchRow
}

package cmpobj

type MarshalRow struct {
	Msg string
	N   int
}

type MarshalData struct {
	Foo  int
	Bar  string
	Rows []MarshalRow
}

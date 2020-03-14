package cbytetpl

import (
	"bytes"
	"testing"
)

type testPermission map[int32]bool
type testFlag map[string]int32

type testHistory struct {
	DateUnix int64
	Cost     float64
	Comment  []byte
}

type testFinance struct {
	MoneyIn  float64
	MoneyOut float64
	Balance  float64
	History  []testHistory
}

type testObject struct {
	Id         string
	Name       []byte
	Cost       float64
	Permission *testPermission
	Flags      testFlag
	Finance    *testFinance
}

var (
	testO = testObject{
		Id:         "foo",
		Name:       []byte("bar"),
		Cost:       12.34,
		Permission: &testPermission{15: true, 23: false},
		Flags: testFlag{
			"export": 17,
			"ro":     4,
			"rw":     7,
			"Valid":  1,
		},
		Finance: &testFinance{
			MoneyIn:  3200,
			MoneyOut: 1500.637657,
			Balance:  9000,
			History: []testHistory{
				{
					152354345634,
					14.345241,
					[]byte("pay for domain"),
				},
				{
					153465345246,
					-3.0000342543,
					[]byte("got refund"),
				},
				{
					156436535640,
					2325242534.35324523,
					[]byte("maintenance"),
				},
			},
		},
	}

	expectFoo     = []byte("bar")
	expectComment = []byte("pay for domain")
)

func testReflect(t testing.TB) {
	i := ReflectInspector{}

	id := i.Get(testO, "Id")
	if id.(string) != "foo" {
		t.Error("object.Id: mismatch result and expectation")
	}

	name := i.Get(testO, "Name")
	if !bytes.Equal(name.([]byte), expectFoo) {
		t.Error("object.Name: mismatch result and expectation")
	}

	perm := i.Get(testO, "Permission", "23")
	if perm.(bool) != false {
		t.Error("object.Permission.23: mismatch result and expectation")
	}

	flag := i.Get(testO, "Flags", "export")
	if flag.(int32) != 17 {
		t.Error("object.Flags.export: mismatch result and expectation")
	}

	bal := i.Get(testO, "Finance", "Balance")
	if bal.(float64) != 9000 {
		t.Error("object.Finance.Balance: mismatch result and expectation")
	}

	date := i.Get(testO, "Finance", "History", 1, "DateUnix")
	if date.(int64) != 153465345246 {
		t.Error("object.Finance.History.1.DateUnix: mismatch result and expectation")
	}

	comment := i.Get(testO, "Finance", "History", 0, "Comment")
	if !bytes.Equal(comment.([]byte), expectComment) {
		t.Error("object.Finance.History.1.DateUnix: mismatch result and expectation")
	}
}

func TestReflectInspector_Get(t *testing.T) {
	testReflect(t)
}

func BenchmarkReflectInspector_Get(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testReflect(b)
	}
}

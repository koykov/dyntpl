package dyntpl

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/koykov/bytealg"
	"github.com/koykov/inspector/testobj"
	"github.com/koykov/inspector/testobj_ins"
)

type stage struct {
	key string
	origin,
	expect []byte
}

var (
	stages []stage

	user = &testobj.TestObject{
		Id:     "115",
		Name:   []byte("John"),
		Status: 78,
		Flags: testobj.TestFlag{
			"export": 17,
			"ro":     4,
			"rw":     7,
			"Valid":  1,
		},
		Finance: &testobj.TestFinance{
			Balance:  9000.015,
			AllowBuy: false,
			History: []testobj.TestHistory{
				{
					DateUnix: 152354345634,
					Cost:     14.345241,
					Comment:  []byte("pay for domain"),
				},
				{
					DateUnix: 153465345246,
					Cost:     -3.0000342543,
					Comment:  []byte("got refund"),
				},
				{
					DateUnix: 156436535640,
					Cost:     2325242534.35324523,
					Comment:  []byte("maintenance"),
				},
			},
		},
	}
	ins testobj_ins.TestObjectInspector

	buf bytes.Buffer
)

func loadStages() {
	if len(stages) > 0 {
		return
	}
	dirs := []string{"tpl", "mod", "i18n"}
	for _, dir := range dirs {
		_ = filepath.Walk("testdata/"+dir, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".tpl" {
				st := stage{}
				st.key = strings.Replace(filepath.Base(path), ".tpl", "", 1)

				st.origin, _ = ioutil.ReadFile(path)
				tree, _ := Parse(st.origin, false)

				st.expect, _ = ioutil.ReadFile(strings.Replace(path, ".tpl", ".txt", 1))
				st.expect = bytealg.Trim(st.expect, []byte("\n"))
				stages = append(stages, st)

				RegisterTplKey(st.key, tree)
			}
			return nil
		})
	}
}

func getStage(key string) (st *stage) {
	for i := 0; i < len(stages); i++ {
		st1 := &stages[i]
		if st1.key == key {
			st = st1
		}
	}
	return st
}

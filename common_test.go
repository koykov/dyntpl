package dyntpl

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koykov/bytealg"
	"github.com/koykov/byteconv"
	"github.com/koykov/inspector/testobj"
	"github.com/koykov/inspector/testobj_ins"
)

type stage struct {
	key, err            string
	origin, expect, raw []byte
}

var (
	stages    []stage
	stagesReg = make(map[string]int)

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

func init() {
	registerTestStages("tpl")
	registerTestStages("mod")
	registerTestStages("datetime")
	registerTestStages("math")
	registerTestStages("fmt")

	_ = filepath.Walk("testdata/parser", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".tpl" {
			st := stage{}
			st.key = strings.Replace(filepath.Base(path), ".tpl", "", 1)
			st.key = "parser/" + st.key

			st.origin, _ = os.ReadFile(path)
			tree, _ := Parse(st.origin, false)

			if raw, err := os.ReadFile(strings.Replace(path, ".tpl", ".xml", 1)); err == nil {
				st.expect = raw
			} else if raw, err := os.ReadFile(strings.Replace(path, ".tpl", ".raw", 1)); err == nil {
				st.raw = bytealg.Trim(raw, []byte("\n"))
			} else if raw, err := os.ReadFile(strings.Replace(path, ".tpl", ".err", 1)); err == nil {
				st.err = bytealg.Trim(byteconv.B2S(raw), "\n")
			}
			stages = append(stages, st)
			stagesReg[st.key] = len(stages) - 1

			RegisterTplKey(st.key, tree)
		}
		return nil
	})
}

func registerTestStages(dir string) {
	_ = filepath.Walk("testdata/"+dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".tpl" {
			st := stage{}
			st.key = strings.Replace(filepath.Base(path), ".tpl", "", 1)

			st.origin, _ = os.ReadFile(path)
			tree, _ := Parse(st.origin, false)

			st.expect, _ = os.ReadFile(strings.Replace(path, ".tpl", ".txt", 1))
			st.expect = bytealg.Trim(st.expect, []byte("\n"))
			stages = append(stages, st)
			stagesReg[st.key] = len(stages) - 1

			RegisterTplKey(st.key, tree)
		}
		return nil
	})
}

func getStage(key string) (st *stage) {
	if i, ok := stagesReg[key]; ok {
		st = &stages[i]
	}
	return st
}

func getTBName(tb testing.TB) string {
	key := tb.Name()
	return key[strings.Index(key, "/")+1:]
}

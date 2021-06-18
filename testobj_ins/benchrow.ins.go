// Code generated by inspc. DO NOT EDIT.
// source: github.com/koykov/dyntpl/testobj

package testobj_ins

import (
	"github.com/koykov/dyntpl/testobj"
	"github.com/koykov/inspector"
	"strconv"
)

type BenchRowInspector struct {
	inspector.BaseInspector
}

func (i0 *BenchRowInspector) Get(src interface{}, path ...string) (interface{}, error) {
	var buf interface{}
	err := i0.GetTo(src, &buf, path...)
	return buf, err
}

func (i0 *BenchRowInspector) GetTo(src interface{}, buf *interface{}, path ...string) (err error) {
	if src == nil {
		return
	}
	var x *testobj.BenchRow
	_ = x
	if p, ok := src.(**testobj.BenchRow); ok {
		x = *p
	} else if p, ok := src.(*testobj.BenchRow); ok {
		x = p
	} else if v, ok := src.(testobj.BenchRow); ok {
		x = &v
	} else {
		return
	}
	if len(path) == 0 {
		*buf = &(*x)
		return
	}

	if len(path) > 0 {
		if path[0] == "ID" {
			*buf = &x.ID
			return
		}
		if path[0] == "Message" {
			*buf = &x.Message
			return
		}
		if path[0] == "Print" {
			*buf = &x.Print
			return
		}
	}
	*buf = &(*x)
	return
}

func (i0 *BenchRowInspector) Cmp(src interface{}, cond inspector.Op, right string, result *bool, path ...string) (err error) {
	if len(path) == 0 {
		return
	}
	if src == nil {
		return
	}
	var x *testobj.BenchRow
	_ = x
	if p, ok := src.(**testobj.BenchRow); ok {
		x = *p
	} else if p, ok := src.(*testobj.BenchRow); ok {
		x = p
	} else if v, ok := src.(testobj.BenchRow); ok {
		x = &v
	} else {
		return
	}

	if len(path) > 0 {
		if path[0] == "ID" {
			var rightExact int
			t0, err0 := strconv.ParseInt(right, 0, 0)
			if err0 != nil {
				return err0
			}
			rightExact = int(t0)
			switch cond {
			case inspector.OpEq:
				*result = x.ID == rightExact
			case inspector.OpNq:
				*result = x.ID != rightExact
			case inspector.OpGt:
				*result = x.ID > rightExact
			case inspector.OpGtq:
				*result = x.ID >= rightExact
			case inspector.OpLt:
				*result = x.ID < rightExact
			case inspector.OpLtq:
				*result = x.ID <= rightExact
			}
			return
		}
		if path[0] == "Message" {
			var rightExact string
			rightExact = right

			switch cond {
			case inspector.OpEq:
				*result = x.Message == rightExact
			case inspector.OpNq:
				*result = x.Message != rightExact
			case inspector.OpGt:
				*result = x.Message > rightExact
			case inspector.OpGtq:
				*result = x.Message >= rightExact
			case inspector.OpLt:
				*result = x.Message < rightExact
			case inspector.OpLtq:
				*result = x.Message <= rightExact
			}
			return
		}
		if path[0] == "Print" {
			var rightExact bool
			t2, err2 := strconv.ParseBool(right)
			if err2 != nil {
				return err2
			}
			rightExact = bool(t2)
			if cond == inspector.OpEq {
				*result = x.Print == rightExact
			} else {
				*result = x.Print != rightExact
			}
			return
		}
	}
	return
}

func (i0 *BenchRowInspector) Loop(src interface{}, l inspector.Looper, buf *[]byte, path ...string) (err error) {
	if len(path) == 0 {
		return
	}
	if src == nil {
		return
	}
	var x *testobj.BenchRow
	_ = x
	if p, ok := src.(**testobj.BenchRow); ok {
		x = *p
	} else if p, ok := src.(*testobj.BenchRow); ok {
		x = p
	} else if v, ok := src.(testobj.BenchRow); ok {
		x = &v
	} else {
		return
	}

	if len(path) > 0 {
	}
	return
}

func (i0 *BenchRowInspector) SetWB(dst, value interface{}, buf inspector.AccumulativeBuffer, path ...string) error {
	if len(path) == 0 {
		return nil
	}
	if dst == nil {
		return nil
	}
	var x *testobj.BenchRow
	_ = x
	if p, ok := dst.(**testobj.BenchRow); ok {
		x = *p
	} else if p, ok := dst.(*testobj.BenchRow); ok {
		x = p
	} else if v, ok := dst.(testobj.BenchRow); ok {
		x = &v
	} else {
		return nil
	}

	if len(path) > 0 {
		if path[0] == "ID" {
			inspector.AssignBuf(&x.ID, value, buf)
			return nil
		}
		if path[0] == "Message" {
			inspector.AssignBuf(&x.Message, value, buf)
			return nil
		}
		if path[0] == "Print" {
			inspector.AssignBuf(&x.Print, value, buf)
			return nil
		}
	}
	return nil
}

func (i0 *BenchRowInspector) Set(dst, value interface{}, path ...string) error {
	return i0.SetWB(dst, value, nil, path...)
}

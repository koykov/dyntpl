// Code generated by inspc. DO NOT EDIT.
// source: github.com/koykov/dyntpl/testobj

package testobj_ins

import (
	"encoding/json"
	"github.com/koykov/dyntpl/testobj"
	"github.com/koykov/inspector"
	"strconv"
)

type BenchRowInspector struct {
	inspector.BaseInspector
}

func (i0 BenchRowInspector) TypeName() string {
	return "BenchRow"
}

func (i0 BenchRowInspector) Get(src any, path ...string) (any, error) {
	var buf any
	err := i0.GetTo(src, &buf, path...)
	return buf, err
}

func (i0 BenchRowInspector) GetTo(src any, buf *any, path ...string) (err error) {
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
	return
}

func (i0 BenchRowInspector) Compare(src any, cond inspector.Op, right string, result *bool, path ...string) (err error) {
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

func (i0 BenchRowInspector) Loop(src any, l inspector.Iterator, buf *[]byte, path ...string) (err error) {
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

func (i0 BenchRowInspector) SetWithBuffer(dst, value any, buf inspector.AccumulativeBuffer, path ...string) error {
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

func (i0 BenchRowInspector) Set(dst, value any, path ...string) error {
	return i0.SetWithBuffer(dst, value, nil, path...)
}

func (i0 BenchRowInspector) DeepEqual(l, r any) bool {
	return i0.DeepEqualWithOptions(l, r, nil)
}

func (i0 BenchRowInspector) DeepEqualWithOptions(l, r any, opts *inspector.DEQOptions) bool {
	var (
		lx, rx   *testobj.BenchRow
		leq, req bool
	)
	_, _, _, _ = lx, rx, leq, req
	if lp, ok := l.(**testobj.BenchRow); ok {
		lx, leq = *lp, true
	} else if lp, ok := l.(*testobj.BenchRow); ok {
		lx, leq = lp, true
	} else if lp, ok := l.(testobj.BenchRow); ok {
		lx, leq = &lp, true
	}
	if rp, ok := r.(**testobj.BenchRow); ok {
		rx, req = *rp, true
	} else if rp, ok := r.(*testobj.BenchRow); ok {
		rx, req = rp, true
	} else if rp, ok := r.(testobj.BenchRow); ok {
		rx, req = &rp, true
	}
	if !leq || !req {
		return false
	}
	if lx == nil && rx == nil {
		return true
	}
	if (lx == nil && rx != nil) || (lx != nil && rx == nil) {
		return false
	}

	if lx.ID != rx.ID && inspector.DEQMustCheck("ID", opts) {
		return false
	}
	if lx.Message != rx.Message && inspector.DEQMustCheck("Message", opts) {
		return false
	}
	if lx.Print != rx.Print && inspector.DEQMustCheck("Print", opts) {
		return false
	}
	return true
}

func (i0 BenchRowInspector) Unmarshal(p []byte, typ inspector.Encoding) (any, error) {
	var x testobj.BenchRow
	switch typ {
	case inspector.EncodingJSON:
		err := json.Unmarshal(p, &x)
		return &x, err
	default:
		return nil, inspector.ErrUnknownEncodingType
	}
}

func (i0 BenchRowInspector) Copy(x any) (any, error) {
	var r testobj.BenchRow
	switch x.(type) {
	case testobj.BenchRow:
		r = x.(testobj.BenchRow)
	case *testobj.BenchRow:
		r = *x.(*testobj.BenchRow)
	case **testobj.BenchRow:
		r = **x.(**testobj.BenchRow)
	default:
		return nil, inspector.ErrUnsupportedType
	}
	bc := i0.countBytes(&r)
	var l testobj.BenchRow
	err := i0.CopyTo(&r, &l, inspector.NewByteBuffer(bc))
	return &l, err
}

func (i0 BenchRowInspector) CopyTo(src, dst any, buf inspector.AccumulativeBuffer) error {
	var r testobj.BenchRow
	switch src.(type) {
	case testobj.BenchRow:
		r = src.(testobj.BenchRow)
	case *testobj.BenchRow:
		r = *src.(*testobj.BenchRow)
	case **testobj.BenchRow:
		r = **src.(**testobj.BenchRow)
	default:
		return inspector.ErrUnsupportedType
	}
	var l *testobj.BenchRow
	switch dst.(type) {
	case testobj.BenchRow:
		return inspector.ErrMustPointerType
	case *testobj.BenchRow:
		l = dst.(*testobj.BenchRow)
	case **testobj.BenchRow:
		l = *dst.(**testobj.BenchRow)
	default:
		return inspector.ErrUnsupportedType
	}
	bb := buf.AcquireBytes()
	var err error
	if bb, err = i0.cpy(bb, l, &r); err != nil {
		return err
	}
	buf.ReleaseBytes(bb)
	return nil
}

func (i0 BenchRowInspector) countBytes(x *testobj.BenchRow) (c int) {
	c += len(x.Message)
	return c
}

func (i0 BenchRowInspector) cpy(buf []byte, l, r *testobj.BenchRow) ([]byte, error) {
	l.ID = r.ID
	buf, l.Message = inspector.BufferizeString(buf, r.Message)
	l.Print = r.Print
	return buf, nil
}

func (i0 BenchRowInspector) Length(src any, result *int, path ...string) error {
	if src == nil {
		return nil
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
		return inspector.ErrUnsupportedType
	}

	*result = 0
	if len(path) == 0 {
		return nil
	}
	if path[0] == "Message" {
		*result = len(x.Message)
		return nil
	}
	return nil
}

func (i0 BenchRowInspector) Capacity(src any, result *int, path ...string) error {
	if src == nil {
		return nil
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
		return inspector.ErrUnsupportedType
	}

	*result = 0
	if len(path) == 0 {
		return nil
	}
	if path[0] == "Message" {
	}
	return nil
}

func (i0 BenchRowInspector) Reset(x any) error {
	var origin *testobj.BenchRow
	_ = origin
	switch x.(type) {
	case testobj.BenchRow:
		return inspector.ErrMustPointerType
	case *testobj.BenchRow:
		origin = x.(*testobj.BenchRow)
	case **testobj.BenchRow:
		origin = *x.(**testobj.BenchRow)
	default:
		return inspector.ErrUnsupportedType
	}
	origin.ID = 0
	origin.Message = ""
	origin.Print = false
	return nil
}

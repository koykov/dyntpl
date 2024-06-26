// Code generated by inspc. DO NOT EDIT.
// source: github.com/koykov/dyntpl/testobj

package testobj_ins

import (
	"encoding/json"
	"github.com/koykov/dyntpl/testobj"
	"github.com/koykov/inspector"
	"strconv"
)

func init() {
	inspector.RegisterInspector("MarshalRow", MarshalRowInspector{})
}

type MarshalRowInspector struct {
	inspector.BaseInspector
}

func (i3 MarshalRowInspector) TypeName() string {
	return "MarshalRow"
}

func (i3 MarshalRowInspector) Get(src any, path ...string) (any, error) {
	var buf any
	err := i3.GetTo(src, &buf, path...)
	return buf, err
}

func (i3 MarshalRowInspector) GetTo(src any, buf *any, path ...string) (err error) {
	if src == nil {
		return
	}
	var x *testobj.MarshalRow
	_ = x
	if p, ok := src.(**testobj.MarshalRow); ok {
		x = *p
	} else if p, ok := src.(*testobj.MarshalRow); ok {
		x = p
	} else if v, ok := src.(testobj.MarshalRow); ok {
		x = &v
	} else {
		return
	}
	if len(path) == 0 {
		*buf = &(*x)
		return
	}

	if len(path) > 0 {
		if path[0] == "Msg" {
			*buf = &x.Msg
			return
		}
		if path[0] == "N" {
			*buf = &x.N
			return
		}
	}
	return
}

func (i3 MarshalRowInspector) Compare(src any, cond inspector.Op, right string, result *bool, path ...string) (err error) {
	if len(path) == 0 {
		return
	}
	if src == nil {
		return
	}
	var x *testobj.MarshalRow
	_ = x
	if p, ok := src.(**testobj.MarshalRow); ok {
		x = *p
	} else if p, ok := src.(*testobj.MarshalRow); ok {
		x = p
	} else if v, ok := src.(testobj.MarshalRow); ok {
		x = &v
	} else {
		return
	}

	if len(path) > 0 {
		if path[0] == "Msg" {
			var rightExact string
			rightExact = right

			switch cond {
			case inspector.OpEq:
				*result = x.Msg == rightExact
			case inspector.OpNq:
				*result = x.Msg != rightExact
			case inspector.OpGt:
				*result = x.Msg > rightExact
			case inspector.OpGtq:
				*result = x.Msg >= rightExact
			case inspector.OpLt:
				*result = x.Msg < rightExact
			case inspector.OpLtq:
				*result = x.Msg <= rightExact
			}
			return
		}
		if path[0] == "N" {
			var rightExact int
			t21, err21 := strconv.ParseInt(right, 0, 0)
			if err21 != nil {
				return err21
			}
			rightExact = int(t21)
			switch cond {
			case inspector.OpEq:
				*result = x.N == rightExact
			case inspector.OpNq:
				*result = x.N != rightExact
			case inspector.OpGt:
				*result = x.N > rightExact
			case inspector.OpGtq:
				*result = x.N >= rightExact
			case inspector.OpLt:
				*result = x.N < rightExact
			case inspector.OpLtq:
				*result = x.N <= rightExact
			}
			return
		}
	}
	return
}

func (i3 MarshalRowInspector) Loop(src any, l inspector.Iterator, buf *[]byte, path ...string) (err error) {
	if len(path) == 0 {
		return
	}
	if src == nil {
		return
	}
	var x *testobj.MarshalRow
	_ = x
	if p, ok := src.(**testobj.MarshalRow); ok {
		x = *p
	} else if p, ok := src.(*testobj.MarshalRow); ok {
		x = p
	} else if v, ok := src.(testobj.MarshalRow); ok {
		x = &v
	} else {
		return
	}

	if len(path) > 0 {
	}
	return
}

func (i3 MarshalRowInspector) SetWithBuffer(dst, value any, buf inspector.AccumulativeBuffer, path ...string) error {
	if len(path) == 0 {
		return nil
	}
	if dst == nil {
		return nil
	}
	var x *testobj.MarshalRow
	_ = x
	if p, ok := dst.(**testobj.MarshalRow); ok {
		x = *p
	} else if p, ok := dst.(*testobj.MarshalRow); ok {
		x = p
	} else if v, ok := dst.(testobj.MarshalRow); ok {
		x = &v
	} else {
		return nil
	}

	if len(path) > 0 {
		if path[0] == "Msg" {
			inspector.AssignBuf(&x.Msg, value, buf)
			return nil
		}
		if path[0] == "N" {
			inspector.AssignBuf(&x.N, value, buf)
			return nil
		}
	}
	return nil
}

func (i3 MarshalRowInspector) Set(dst, value any, path ...string) error {
	return i3.SetWithBuffer(dst, value, nil, path...)
}

func (i3 MarshalRowInspector) DeepEqual(l, r any) bool {
	return i3.DeepEqualWithOptions(l, r, nil)
}

func (i3 MarshalRowInspector) DeepEqualWithOptions(l, r any, opts *inspector.DEQOptions) bool {
	var (
		lx, rx   *testobj.MarshalRow
		leq, req bool
	)
	_, _, _, _ = lx, rx, leq, req
	if lp, ok := l.(**testobj.MarshalRow); ok {
		lx, leq = *lp, true
	} else if lp, ok := l.(*testobj.MarshalRow); ok {
		lx, leq = lp, true
	} else if lp, ok := l.(testobj.MarshalRow); ok {
		lx, leq = &lp, true
	}
	if rp, ok := r.(**testobj.MarshalRow); ok {
		rx, req = *rp, true
	} else if rp, ok := r.(*testobj.MarshalRow); ok {
		rx, req = rp, true
	} else if rp, ok := r.(testobj.MarshalRow); ok {
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

	if lx.Msg != rx.Msg && inspector.DEQMustCheck("Msg", opts) {
		return false
	}
	if lx.N != rx.N && inspector.DEQMustCheck("N", opts) {
		return false
	}
	return true
}

func (i3 MarshalRowInspector) Unmarshal(p []byte, typ inspector.Encoding) (any, error) {
	var x testobj.MarshalRow
	switch typ {
	case inspector.EncodingJSON:
		err := json.Unmarshal(p, &x)
		return &x, err
	default:
		return nil, inspector.ErrUnknownEncodingType
	}
}

func (i3 MarshalRowInspector) Copy(x any) (any, error) {
	var r testobj.MarshalRow
	switch x.(type) {
	case testobj.MarshalRow:
		r = x.(testobj.MarshalRow)
	case *testobj.MarshalRow:
		r = *x.(*testobj.MarshalRow)
	case **testobj.MarshalRow:
		r = **x.(**testobj.MarshalRow)
	default:
		return nil, inspector.ErrUnsupportedType
	}
	bc := i3.countBytes(&r)
	var l testobj.MarshalRow
	err := i3.CopyTo(&r, &l, inspector.NewByteBuffer(bc))
	return &l, err
}

func (i3 MarshalRowInspector) CopyTo(src, dst any, buf inspector.AccumulativeBuffer) error {
	var r testobj.MarshalRow
	switch src.(type) {
	case testobj.MarshalRow:
		r = src.(testobj.MarshalRow)
	case *testobj.MarshalRow:
		r = *src.(*testobj.MarshalRow)
	case **testobj.MarshalRow:
		r = **src.(**testobj.MarshalRow)
	default:
		return inspector.ErrUnsupportedType
	}
	var l *testobj.MarshalRow
	switch dst.(type) {
	case testobj.MarshalRow:
		return inspector.ErrMustPointerType
	case *testobj.MarshalRow:
		l = dst.(*testobj.MarshalRow)
	case **testobj.MarshalRow:
		l = *dst.(**testobj.MarshalRow)
	default:
		return inspector.ErrUnsupportedType
	}
	bb := buf.AcquireBytes()
	var err error
	if bb, err = i3.cpy(bb, l, &r); err != nil {
		return err
	}
	buf.ReleaseBytes(bb)
	return nil
}

func (i3 MarshalRowInspector) countBytes(x *testobj.MarshalRow) (c int) {
	c += len(x.Msg)
	return c
}

func (i3 MarshalRowInspector) cpy(buf []byte, l, r *testobj.MarshalRow) ([]byte, error) {
	buf, l.Msg = inspector.BufferizeString(buf, r.Msg)
	l.N = r.N
	return buf, nil
}

func (i3 MarshalRowInspector) Length(src any, result *int, path ...string) error {
	if src == nil {
		return nil
	}
	var x *testobj.MarshalRow
	_ = x
	if p, ok := src.(**testobj.MarshalRow); ok {
		x = *p
	} else if p, ok := src.(*testobj.MarshalRow); ok {
		x = p
	} else if v, ok := src.(testobj.MarshalRow); ok {
		x = &v
	} else {
		return inspector.ErrUnsupportedType
	}

	*result = 0
	if len(path) == 0 {
		return nil
	}
	if path[0] == "Msg" {
		*result = len(x.Msg)
		return nil
	}
	return nil
}

func (i3 MarshalRowInspector) Capacity(src any, result *int, path ...string) error {
	if src == nil {
		return nil
	}
	var x *testobj.MarshalRow
	_ = x
	if p, ok := src.(**testobj.MarshalRow); ok {
		x = *p
	} else if p, ok := src.(*testobj.MarshalRow); ok {
		x = p
	} else if v, ok := src.(testobj.MarshalRow); ok {
		x = &v
	} else {
		return inspector.ErrUnsupportedType
	}

	*result = 0
	if len(path) == 0 {
		return nil
	}
	if path[0] == "Msg" {
	}
	return nil
}

func (i3 MarshalRowInspector) Reset(x any) error {
	var origin *testobj.MarshalRow
	_ = origin
	switch x.(type) {
	case testobj.MarshalRow:
		return inspector.ErrMustPointerType
	case *testobj.MarshalRow:
		origin = x.(*testobj.MarshalRow)
	case **testobj.MarshalRow:
		origin = *x.(**testobj.MarshalRow)
	default:
		return inspector.ErrUnsupportedType
	}
	origin.Msg = ""
	origin.N = 0
	return nil
}

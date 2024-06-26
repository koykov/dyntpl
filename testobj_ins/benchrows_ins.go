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
	inspector.RegisterInspector("BenchRows", BenchRowsInspector{})
}

type BenchRowsInspector struct {
	inspector.BaseInspector
}

func (i1 BenchRowsInspector) TypeName() string {
	return "BenchRows"
}

func (i1 BenchRowsInspector) Get(src any, path ...string) (any, error) {
	var buf any
	err := i1.GetTo(src, &buf, path...)
	return buf, err
}

func (i1 BenchRowsInspector) GetTo(src any, buf *any, path ...string) (err error) {
	if src == nil {
		return
	}
	var x *testobj.BenchRows
	_ = x
	if p, ok := src.(**testobj.BenchRows); ok {
		x = *p
	} else if p, ok := src.(*testobj.BenchRows); ok {
		x = p
	} else if v, ok := src.(testobj.BenchRows); ok {
		x = &v
	} else {
		return
	}
	if len(path) == 0 {
		*buf = &(*x)
		return
	}

	if len(path) > 0 {
		if path[0] == "Rows" {
			x0 := x.Rows
			_ = x0
			if len(path) > 1 {
				var i int
				t3, err3 := strconv.ParseInt(path[1], 0, 0)
				if err3 != nil {
					return err3
				}
				i = int(t3)
				if len(x0) > i {
					x1 := &(x0)[i]
					_ = x1
					if len(path) > 2 {
						if path[2] == "ID" {
							*buf = &x1.ID
							return
						}
						if path[2] == "Message" {
							*buf = &x1.Message
							return
						}
						if path[2] == "Print" {
							*buf = &x1.Print
							return
						}
					}
					*buf = x1
				}
			}
			*buf = &x.Rows
			return
		}
	}
	return
}

func (i1 BenchRowsInspector) Compare(src any, cond inspector.Op, right string, result *bool, path ...string) (err error) {
	if len(path) == 0 {
		return
	}
	if src == nil {
		return
	}
	var x *testobj.BenchRows
	_ = x
	if p, ok := src.(**testobj.BenchRows); ok {
		x = *p
	} else if p, ok := src.(*testobj.BenchRows); ok {
		x = p
	} else if v, ok := src.(testobj.BenchRows); ok {
		x = &v
	} else {
		return
	}

	if len(path) > 0 {
		if path[0] == "Rows" {
			x0 := x.Rows
			_ = x0
			if len(path) > 1 {
				var i int
				t4, err4 := strconv.ParseInt(path[1], 0, 0)
				if err4 != nil {
					return err4
				}
				i = int(t4)
				if len(x0) > i {
					x1 := &(x0)[i]
					_ = x1
					if len(path) > 2 {
						if path[2] == "ID" {
							var rightExact int
							t5, err5 := strconv.ParseInt(right, 0, 0)
							if err5 != nil {
								return err5
							}
							rightExact = int(t5)
							switch cond {
							case inspector.OpEq:
								*result = x1.ID == rightExact
							case inspector.OpNq:
								*result = x1.ID != rightExact
							case inspector.OpGt:
								*result = x1.ID > rightExact
							case inspector.OpGtq:
								*result = x1.ID >= rightExact
							case inspector.OpLt:
								*result = x1.ID < rightExact
							case inspector.OpLtq:
								*result = x1.ID <= rightExact
							}
							return
						}
						if path[2] == "Message" {
							var rightExact string
							rightExact = right

							switch cond {
							case inspector.OpEq:
								*result = x1.Message == rightExact
							case inspector.OpNq:
								*result = x1.Message != rightExact
							case inspector.OpGt:
								*result = x1.Message > rightExact
							case inspector.OpGtq:
								*result = x1.Message >= rightExact
							case inspector.OpLt:
								*result = x1.Message < rightExact
							case inspector.OpLtq:
								*result = x1.Message <= rightExact
							}
							return
						}
						if path[2] == "Print" {
							var rightExact bool
							t7, err7 := strconv.ParseBool(right)
							if err7 != nil {
								return err7
							}
							rightExact = bool(t7)
							if cond == inspector.OpEq {
								*result = x1.Print == rightExact
							} else {
								*result = x1.Print != rightExact
							}
							return
						}
					}
				}
			}
		}
	}
	return
}

func (i1 BenchRowsInspector) Loop(src any, l inspector.Iterator, buf *[]byte, path ...string) (err error) {
	if len(path) == 0 {
		return
	}
	if src == nil {
		return
	}
	var x *testobj.BenchRows
	_ = x
	if p, ok := src.(**testobj.BenchRows); ok {
		x = *p
	} else if p, ok := src.(*testobj.BenchRows); ok {
		x = p
	} else if v, ok := src.(testobj.BenchRows); ok {
		x = &v
	} else {
		return
	}

	if len(path) > 0 {
		if path[0] == "Rows" {
			x0 := x.Rows
			_ = x0
			for k := range x0 {
				if l.RequireKey() {
					*buf = strconv.AppendInt((*buf)[:0], int64(k), 10)
					l.SetKey(buf, &inspector.StaticInspector{})
				}
				l.SetVal(&(x0)[k], &BenchRowInspector{})
				ctl := l.Iterate()
				if ctl == inspector.LoopCtlBrk {
					break
				}
				if ctl == inspector.LoopCtlCnt {
					continue
				}
			}
			return
		}
	}
	return
}

func (i1 BenchRowsInspector) SetWithBuffer(dst, value any, buf inspector.AccumulativeBuffer, path ...string) error {
	if len(path) == 0 {
		return nil
	}
	if dst == nil {
		return nil
	}
	var x *testobj.BenchRows
	_ = x
	if p, ok := dst.(**testobj.BenchRows); ok {
		x = *p
	} else if p, ok := dst.(*testobj.BenchRows); ok {
		x = p
	} else if v, ok := dst.(testobj.BenchRows); ok {
		x = &v
	} else {
		return nil
	}

	if len(path) > 0 {
		if path[0] == "Rows" {
			x0 := x.Rows
			if uvalue, ok := value.(*[]testobj.BenchRow); ok {
				x0 = *uvalue
			}
			if x0 == nil {
				z := make([]testobj.BenchRow, 0)
				x0 = z
				x.Rows = x0
			}
			_ = x0
			if len(path) > 1 {
				var i int
				t8, err8 := strconv.ParseInt(path[1], 0, 0)
				if err8 != nil {
					return err8
				}
				i = int(t8)
				if len(x0) > i {
					x1 := &(x0)[i]
					_ = x1
					if len(path) > 2 {
						if path[2] == "ID" {
							inspector.AssignBuf(&x1.ID, value, buf)
							return nil
						}
						if path[2] == "Message" {
							inspector.AssignBuf(&x1.Message, value, buf)
							return nil
						}
						if path[2] == "Print" {
							inspector.AssignBuf(&x1.Print, value, buf)
							return nil
						}
					}
					(x0)[i] = *x1
					return nil
				}
			}
			x.Rows = x0
		}
	}
	return nil
}

func (i1 BenchRowsInspector) Set(dst, value any, path ...string) error {
	return i1.SetWithBuffer(dst, value, nil, path...)
}

func (i1 BenchRowsInspector) DeepEqual(l, r any) bool {
	return i1.DeepEqualWithOptions(l, r, nil)
}

func (i1 BenchRowsInspector) DeepEqualWithOptions(l, r any, opts *inspector.DEQOptions) bool {
	var (
		lx, rx   *testobj.BenchRows
		leq, req bool
	)
	_, _, _, _ = lx, rx, leq, req
	if lp, ok := l.(**testobj.BenchRows); ok {
		lx, leq = *lp, true
	} else if lp, ok := l.(*testobj.BenchRows); ok {
		lx, leq = lp, true
	} else if lp, ok := l.(testobj.BenchRows); ok {
		lx, leq = &lp, true
	}
	if rp, ok := r.(**testobj.BenchRows); ok {
		rx, req = *rp, true
	} else if rp, ok := r.(*testobj.BenchRows); ok {
		rx, req = rp, true
	} else if rp, ok := r.(testobj.BenchRows); ok {
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

	lx1 := lx.Rows
	rx1 := rx.Rows
	_, _ = lx1, rx1
	if inspector.DEQMustCheck("Rows", opts) {
		if len(lx1) != len(rx1) {
			return false
		}
		for i := 0; i < len(lx1); i++ {
			lx2 := (lx1)[i]
			rx2 := (rx1)[i]
			_, _ = lx2, rx2
			if lx2.ID != rx2.ID && inspector.DEQMustCheck("Rows.ID", opts) {
				return false
			}
			if lx2.Message != rx2.Message && inspector.DEQMustCheck("Rows.Message", opts) {
				return false
			}
			if lx2.Print != rx2.Print && inspector.DEQMustCheck("Rows.Print", opts) {
				return false
			}
		}
	}
	return true
}

func (i1 BenchRowsInspector) Unmarshal(p []byte, typ inspector.Encoding) (any, error) {
	var x testobj.BenchRows
	switch typ {
	case inspector.EncodingJSON:
		err := json.Unmarshal(p, &x)
		return &x, err
	default:
		return nil, inspector.ErrUnknownEncodingType
	}
}

func (i1 BenchRowsInspector) Copy(x any) (any, error) {
	var r testobj.BenchRows
	switch x.(type) {
	case testobj.BenchRows:
		r = x.(testobj.BenchRows)
	case *testobj.BenchRows:
		r = *x.(*testobj.BenchRows)
	case **testobj.BenchRows:
		r = **x.(**testobj.BenchRows)
	default:
		return nil, inspector.ErrUnsupportedType
	}
	bc := i1.countBytes(&r)
	var l testobj.BenchRows
	err := i1.CopyTo(&r, &l, inspector.NewByteBuffer(bc))
	return &l, err
}

func (i1 BenchRowsInspector) CopyTo(src, dst any, buf inspector.AccumulativeBuffer) error {
	var r testobj.BenchRows
	switch src.(type) {
	case testobj.BenchRows:
		r = src.(testobj.BenchRows)
	case *testobj.BenchRows:
		r = *src.(*testobj.BenchRows)
	case **testobj.BenchRows:
		r = **src.(**testobj.BenchRows)
	default:
		return inspector.ErrUnsupportedType
	}
	var l *testobj.BenchRows
	switch dst.(type) {
	case testobj.BenchRows:
		return inspector.ErrMustPointerType
	case *testobj.BenchRows:
		l = dst.(*testobj.BenchRows)
	case **testobj.BenchRows:
		l = *dst.(**testobj.BenchRows)
	default:
		return inspector.ErrUnsupportedType
	}
	bb := buf.AcquireBytes()
	var err error
	if bb, err = i1.cpy(bb, l, &r); err != nil {
		return err
	}
	buf.ReleaseBytes(bb)
	return nil
}

func (i1 BenchRowsInspector) countBytes(x *testobj.BenchRows) (c int) {
	for i1 := 0; i1 < len(x.Rows); i1++ {
		x1 := &(x.Rows)[i1]
		c += len(x1.Message)
	}
	return c
}

func (i1 BenchRowsInspector) cpy(buf []byte, l, r *testobj.BenchRows) ([]byte, error) {
	if len(r.Rows) > 0 {
		buf1 := (l.Rows)
		if buf1 == nil {
			buf1 = make([]testobj.BenchRow, 0, len(r.Rows))
		}
		for i1 := 0; i1 < len(r.Rows); i1++ {
			var b1 testobj.BenchRow
			x1 := &(r.Rows)[i1]
			b1.ID = x1.ID
			buf, b1.Message = inspector.BufferizeString(buf, x1.Message)
			b1.Print = x1.Print
			buf1 = append(buf1, b1)
		}
		l.Rows = buf1
	}
	return buf, nil
}

func (i1 BenchRowsInspector) Length(src any, result *int, path ...string) error {
	if src == nil {
		return nil
	}
	var x *testobj.BenchRows
	_ = x
	if p, ok := src.(**testobj.BenchRows); ok {
		x = *p
	} else if p, ok := src.(*testobj.BenchRows); ok {
		x = p
	} else if v, ok := src.(testobj.BenchRows); ok {
		x = &v
	} else {
		return inspector.ErrUnsupportedType
	}

	*result = 0
	if len(path) == 0 {
		return nil
	}
	if path[0] == "Rows" {
		if len(path) == 1 {
			*result = len(x.Rows)
			return nil
		}
		if len(path) < 2 {
			return nil
		}
		var i int
		t9, err9 := strconv.ParseInt(path[1], 0, 0)
		if err9 != nil {
			return err9
		}
		i = int(t9)
		if len(x.Rows) > i {
			x1 := &(x.Rows)[i]
			_ = x1
			if len(path) < 3 {
				return nil
			}
			if path[2] == "Message" {
				*result = len(x1.Message)
				return nil
			}
		}
	}
	return nil
}

func (i1 BenchRowsInspector) Capacity(src any, result *int, path ...string) error {
	if src == nil {
		return nil
	}
	var x *testobj.BenchRows
	_ = x
	if p, ok := src.(**testobj.BenchRows); ok {
		x = *p
	} else if p, ok := src.(*testobj.BenchRows); ok {
		x = p
	} else if v, ok := src.(testobj.BenchRows); ok {
		x = &v
	} else {
		return inspector.ErrUnsupportedType
	}

	*result = 0
	if len(path) == 0 {
		return nil
	}
	if path[0] == "Rows" {
		if len(path) == 1 {
			*result = cap(x.Rows)
			return nil
		}
		if len(path) < 2 {
			return nil
		}
		var i int
		t10, err10 := strconv.ParseInt(path[1], 0, 0)
		if err10 != nil {
			return err10
		}
		i = int(t10)
		if len(x.Rows) > i {
			x1 := &(x.Rows)[i]
			_ = x1
			if len(path) < 3 {
				return nil
			}
			if path[2] == "Message" {
			}
		}
	}
	return nil
}

func (i1 BenchRowsInspector) Reset(x any) error {
	var origin *testobj.BenchRows
	_ = origin
	switch x.(type) {
	case testobj.BenchRows:
		return inspector.ErrMustPointerType
	case *testobj.BenchRows:
		origin = x.(*testobj.BenchRows)
	case **testobj.BenchRows:
		origin = *x.(**testobj.BenchRows)
	default:
		return inspector.ErrUnsupportedType
	}
	if l := len((origin.Rows)); l > 0 {
		_ = (origin.Rows)[l-1]
		for i := 0; i < l; i++ {
			x1 := &(origin.Rows)[i]
			x1.ID = 0
			x1.Message = ""
			x1.Print = false
		}
		(origin.Rows) = (origin.Rows)[:0]
	}
	return nil
}

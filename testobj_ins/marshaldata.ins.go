// Code generated by inspc. DO NOT EDIT.
// source: github.com/koykov/dyntpl/testobj

package testobj_ins

import (
	"github.com/koykov/dyntpl/testobj"
	"github.com/koykov/inspector"
	"strconv"
)

type MarshalDataInspector struct {
	inspector.BaseInspector
}

func (i2 *MarshalDataInspector) Get(src interface{}, path ...string) (interface{}, error) {
	var buf interface{}
	err := i2.GetTo(src, &buf, path...)
	return buf, err
}

func (i2 *MarshalDataInspector) GetTo(src interface{}, buf *interface{}, path ...string) (err error) {
	if src == nil {
		return
	}
	var x *testobj.MarshalData
	_ = x
	if p, ok := src.(**testobj.MarshalData); ok {
		x = *p
	} else if p, ok := src.(*testobj.MarshalData); ok {
		x = p
	} else if v, ok := src.(testobj.MarshalData); ok {
		x = &v
	} else {
		return
	}
	if len(path) == 0 {
		*buf = &(*x)
		return
	}

	if len(path) > 0 {
		if path[0] == "Foo" {
			*buf = &x.Foo
			return
		}
		if path[0] == "Bar" {
			*buf = &x.Bar
			return
		}
		if path[0] == "Rows" {
			x0 := x.Rows
			_ = x0
			if len(path) > 1 {
				var i int
				t9, err9 := strconv.ParseInt(path[1], 0, 0)
				if err9 != nil {
					return err9
				}
				i = int(t9)
				if len(x0) > i {
					x1 := &(x0)[i]
					_ = x1
					if len(path) > 2 {
						if path[2] == "Msg" {
							*buf = &x1.Msg
							return
						}
						if path[2] == "N" {
							*buf = &x1.N
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

func (i2 *MarshalDataInspector) Cmp(src interface{}, cond inspector.Op, right string, result *bool, path ...string) (err error) {
	if len(path) == 0 {
		return
	}
	if src == nil {
		return
	}
	var x *testobj.MarshalData
	_ = x
	if p, ok := src.(**testobj.MarshalData); ok {
		x = *p
	} else if p, ok := src.(*testobj.MarshalData); ok {
		x = p
	} else if v, ok := src.(testobj.MarshalData); ok {
		x = &v
	} else {
		return
	}

	if len(path) > 0 {
		if path[0] == "Foo" {
			var rightExact int
			t10, err10 := strconv.ParseInt(right, 0, 0)
			if err10 != nil {
				return err10
			}
			rightExact = int(t10)
			switch cond {
			case inspector.OpEq:
				*result = x.Foo == rightExact
			case inspector.OpNq:
				*result = x.Foo != rightExact
			case inspector.OpGt:
				*result = x.Foo > rightExact
			case inspector.OpGtq:
				*result = x.Foo >= rightExact
			case inspector.OpLt:
				*result = x.Foo < rightExact
			case inspector.OpLtq:
				*result = x.Foo <= rightExact
			}
			return
		}
		if path[0] == "Bar" {
			var rightExact string
			rightExact = right

			switch cond {
			case inspector.OpEq:
				*result = x.Bar == rightExact
			case inspector.OpNq:
				*result = x.Bar != rightExact
			case inspector.OpGt:
				*result = x.Bar > rightExact
			case inspector.OpGtq:
				*result = x.Bar >= rightExact
			case inspector.OpLt:
				*result = x.Bar < rightExact
			case inspector.OpLtq:
				*result = x.Bar <= rightExact
			}
			return
		}
		if path[0] == "Rows" {
			x0 := x.Rows
			_ = x0
			if len(path) > 1 {
				var i int
				t12, err12 := strconv.ParseInt(path[1], 0, 0)
				if err12 != nil {
					return err12
				}
				i = int(t12)
				if len(x0) > i {
					x1 := &(x0)[i]
					_ = x1
					if len(path) > 2 {
						if path[2] == "Msg" {
							var rightExact string
							rightExact = right

							switch cond {
							case inspector.OpEq:
								*result = x1.Msg == rightExact
							case inspector.OpNq:
								*result = x1.Msg != rightExact
							case inspector.OpGt:
								*result = x1.Msg > rightExact
							case inspector.OpGtq:
								*result = x1.Msg >= rightExact
							case inspector.OpLt:
								*result = x1.Msg < rightExact
							case inspector.OpLtq:
								*result = x1.Msg <= rightExact
							}
							return
						}
						if path[2] == "N" {
							var rightExact int
							t14, err14 := strconv.ParseInt(right, 0, 0)
							if err14 != nil {
								return err14
							}
							rightExact = int(t14)
							switch cond {
							case inspector.OpEq:
								*result = x1.N == rightExact
							case inspector.OpNq:
								*result = x1.N != rightExact
							case inspector.OpGt:
								*result = x1.N > rightExact
							case inspector.OpGtq:
								*result = x1.N >= rightExact
							case inspector.OpLt:
								*result = x1.N < rightExact
							case inspector.OpLtq:
								*result = x1.N <= rightExact
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

func (i2 *MarshalDataInspector) Loop(src interface{}, l inspector.Looper, buf *[]byte, path ...string) (err error) {
	if len(path) == 0 {
		return
	}
	if src == nil {
		return
	}
	var x *testobj.MarshalData
	_ = x
	if p, ok := src.(**testobj.MarshalData); ok {
		x = *p
	} else if p, ok := src.(*testobj.MarshalData); ok {
		x = p
	} else if v, ok := src.(testobj.MarshalData); ok {
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
				l.SetVal(&(x0)[k], &MarshalRowInspector{})
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

func (i2 *MarshalDataInspector) SetWB(dst, value interface{}, buf inspector.AccumulativeBuffer, path ...string) error {
	if len(path) == 0 {
		return nil
	}
	if dst == nil {
		return nil
	}
	var x *testobj.MarshalData
	_ = x
	if p, ok := dst.(**testobj.MarshalData); ok {
		x = *p
	} else if p, ok := dst.(*testobj.MarshalData); ok {
		x = p
	} else if v, ok := dst.(testobj.MarshalData); ok {
		x = &v
	} else {
		return nil
	}

	if len(path) > 0 {
		if path[0] == "Foo" {
			inspector.AssignBuf(&x.Foo, value, buf)
			return nil
		}
		if path[0] == "Bar" {
			inspector.AssignBuf(&x.Bar, value, buf)
			return nil
		}
		if path[0] == "Rows" {
			x0 := x.Rows
			if uvalue, ok := value.(*[]testobj.MarshalRow); ok {
				x0 = *uvalue
			}
			if x0 == nil {
				z := make([]testobj.MarshalRow, 0)
				x0 = z
				x.Rows = x0
			}
			_ = x0
			if len(path) > 1 {
				var i int
				t15, err15 := strconv.ParseInt(path[1], 0, 0)
				if err15 != nil {
					return err15
				}
				i = int(t15)
				if len(x0) > i {
					x1 := &(x0)[i]
					_ = x1
					if len(path) > 2 {
						if path[2] == "Msg" {
							inspector.AssignBuf(&x1.Msg, value, buf)
							return nil
						}
						if path[2] == "N" {
							inspector.AssignBuf(&x1.N, value, buf)
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

func (i2 *MarshalDataInspector) Set(dst, value interface{}, path ...string) error {
	return i2.SetWB(dst, value, nil, path...)
}

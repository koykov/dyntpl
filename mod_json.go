package dyntpl

import "github.com/koykov/fastconv"

var (
	jqQd = byte('"')
	jqSl = byte('\\')
	jqNl = byte('\n')
	jqCr = byte('\r')
	jqT  = byte('\t')
	jqFf = byte('\f')
	jqBs = byte('\f')
	jqLt = byte('<')
	jqQs = byte('\'')
	jqZ  = byte(0)

	jqQdR = []byte(`\"`)
	jqSlR = []byte("\\")
	jqNlR = []byte("\n")
	jqCrR = []byte("\r")
	jqTR  = []byte("\t")
	jqFfR = []byte("\u000c")
	jqBsR = []byte("\u0008")
	jqLtR = []byte("\u003c")
	jqQsR = []byte("\u0027")
	jqZR  = []byte("\u0000")
)

func modJsonQuote(ctx *Ctx, val interface{}, _ []interface{}) (interface{}, error) {
	var (
		b    []byte
		dst  = ctx.GetBbuf()
		l, o int
	)
	switch val.(type) {
	case []byte:
		b = val.([]byte)
	case *[]byte:
		b = *val.(*[]byte)
	case string:
		b = fastconv.S2B(val.(string))
	case *string:
		b = fastconv.S2B(*val.(*string))
	default:
		return val, ErrModNoStr
	}
	l = len(b)
	if l == 0 {
		return val, nil
	}
	_ = b[l-1]
	for i := 0; i < l; i++ {
		switch b[i] {
		case jqQd:
			dst = append(dst, b[o:i]...)
			dst = append(dst, jqQdR...)
			o = i + 1
		case jqSl:
			dst = append(dst, b[o:i]...)
			dst = append(dst, jqSlR...)
			o = i + 1
		case jqNl:
			dst = append(dst, b[o:i]...)
			dst = append(dst, jqNlR...)
			o = i + 1
		case jqCr:
			dst = append(dst, b[o:i]...)
			dst = append(dst, jqCrR...)
			o = i + 1
		case jqT:
			dst = append(dst, b[o:i]...)
			dst = append(dst, jqTR...)
			o = i + 1
		case jqFf:
			dst = append(dst, b[o:i]...)
			dst = append(dst, jqFfR...)
			o = i + 1
		case jqBs:
			dst = append(dst, b[o:i]...)
			dst = append(dst, jqBsR...)
			o = i + 1
		case jqLt:
			dst = append(dst, b[o:i]...)
			dst = append(dst, jqLtR...)
			o = i + 1
		case jqQs:
			dst = append(dst, b[o:i]...)
			dst = append(dst, jqQsR...)
			o = i + 1
		case jqZ:
			dst = append(dst, b[o:i]...)
			dst = append(dst, jqZR...)
			o = i + 1
		}
	}
	dst = append(dst, b[o:]...)

	return dst, nil
}

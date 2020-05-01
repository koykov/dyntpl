package cbytetpl

type ByteConvFn func(buf []byte, val interface{}) ([]byte, error)

var (
	byteConvFnRegistry = make([]ByteConvFn, 0)
)

func RegisterByteConvFn(fn ByteConvFn) {
	for _, f := range byteConvFnRegistry {
		if &f == &fn {
			return
		}
	}
	byteConvFnRegistry = append(byteConvFnRegistry, fn)
}

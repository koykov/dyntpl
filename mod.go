package dyntpl

type ModFn func(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) error

type mod struct {
	id  []byte
	fn  *ModFn
	arg []*modArg
}

type modArg struct {
	val    []byte
	static bool
}

var (
	modRegistry = map[string]ModFn{}
)

func RegisterModFn(name, alias string, mod ModFn) {
	modRegistry[name] = mod
	if len(alias) > 0 {
		modRegistry[alias] = mod
	}
}

func GetModFn(name string) *ModFn {
	if fn, ok := modRegistry[name]; ok {
		return &fn
	}
	return nil
}

func ModInt(val interface{}) (i int64, ok bool) {
	ok = true
	switch val.(type) {
	case int:
		i = int64(val.(int))
	case *int:
		i = int64(*val.(*int))
	case int8:
		i = int64(val.(int8))
	case *int8:
		i = int64(*val.(*int8))
	case int16:
		i = int64(val.(int16))
	case *int16:
		i = int64(*val.(*int16))
	case int32:
		i = int64(val.(int32))
	case *int32:
		i = int64(*val.(*int32))
	case int64:
		i = val.(int64)
	case *int64:
		i = *val.(*int64)
	default:
		ok = false
	}
	return
}

func ModUint(val interface{}) (u uint64, ok bool) {
	ok = true
	switch val.(type) {
	case uint:
		u = uint64(val.(uint))
	case *uint:
		u = uint64(*val.(*uint))
	case uint8:
		u = uint64(val.(uint8))
	case *uint8:
		u = uint64(*val.(*uint8))
	case uint16:
		u = uint64(val.(uint16))
	case *uint16:
		u = uint64(*val.(*uint16))
	case uint32:
		u = uint64(val.(uint32))
	case *uint32:
		u = uint64(*val.(*uint32))
	case uint64:
		u = val.(uint64)
	case *uint64:
		u = *val.(*uint64)
	default:
		ok = false
	}
	return
}

func ModFloat(val interface{}) (f float64, ok bool) {
	ok = true
	switch val.(type) {
	case float32:
		f = float64(val.(float32))
	case *float32:
		f = float64(*val.(*float32))
	case float64:
		f = val.(float64)
	case *float64:
		f = *val.(*float64)
	default:
		ok = false
	}
	return
}

func ModBytes(val interface{}) (b []byte, ok bool) {
	ok = true
	switch val.(type) {
	case []byte:
		b = val.([]byte)
	case *[]byte:
		b = *val.(*[]byte)
	default:
		ok = false
	}
	return
}

func ModStr(val interface{}) (s string, ok bool) {
	ok = true
	switch val.(type) {
	case string:
		s = val.(string)
	case *string:
		s = *val.(*string)
	default:
		ok = false
	}
	return
}

func ModBool(val interface{}) (b bool, ok bool) {
	ok = true
	switch val.(type) {
	case bool:
		b = val.(bool)
	case *bool:
		b = *val.(*bool)
	default:
		ok = false
	}
	return
}

package dyntpl

func modDefault(_ *Ctx, buf *interface{}, val interface{}, args []interface{}) error {
	if len(args) == 0 {
		return ErrModNoArgs
	}
	b := false
	switch val.(type) {
	case int:
		b = val.(int) == 0
	case *int:
		b = *val.(*int) == 0
	case int8:
		b = val.(int8) == 0
	case *int8:
		b = *val.(*int8) == 0
	case int16:
		b = val.(int16) == 0
	case *int16:
		b = *val.(*int16) == 0
	case int32:
		b = val.(int32) == 0
	case *int32:
		b = *val.(*int32) == 0
	case int64:
		b = val.(int64) == 0
	case *int64:
		b = *val.(*int64) == 0
	case uint:
		b = val.(uint) == 0
	case *uint:
		b = *val.(*uint) == 0
	case uint8:
		b = val.(uint8) == 0
	case *uint8:
		b = *val.(*uint8) == 0
	case uint16:
		b = val.(uint16) == 0
	case *uint16:
		b = *val.(*uint16) == 0
	case uint32:
		b = val.(uint32) == 0
	case *uint32:
		b = *val.(*uint32) == 0
	case uint64:
		b = val.(uint64) == 0
	case *uint64:
		b = *val.(*uint64) == 0
	case float32:
		b = val.(float32) == 0
	case *float32:
		b = *val.(*float32) == 0
	case float64:
		b = val.(float64) == 0
	case *float64:
		b = *val.(*float64) == 0
	case []byte:
		b = len(val.([]byte)) == 0
	case *[]byte:
		b = len(*val.(*[]byte)) == 0
	case string:
		b = len(val.(string)) == 0
	case *string:
		b = len(*val.(*string)) == 0
	}
	if b {
		*buf = args[0]
	}
	return nil
}

package dyntpl

func modDefault(_ *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if len(args) == 0 {
		err = ErrModNoArgs
		return
	}
	if i, ok := ConvInt(val); ok {
		if i == 0 {
			*buf = args[0]
			return
		}
	}
	if u, ok := ConvUint(val); ok {
		if u == 0 {
			*buf = args[0]
			return
		}
	}
	if f, ok := ConvFloat(val); ok {
		if f == 0 {
			*buf = args[0]
			return
		}
	}
	if b, ok := ConvBytes(val); ok {
		if len(b) == 0 {
			*buf = args[0]
			return
		}
	}
	if s, ok := ConvStr(val); ok {
		if len(s) == 0 {
			*buf = args[0]
			return
		}
	}
	if b, ok := ConvBool(val); ok {
		if !b {
			*buf = args[0]
			return
		}
	}
	return nil
}

func modIfThen(_ *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if len(args) == 0 {
		err = ErrModNoArgs
		return
	}
	if b, ok := ConvBool(val); ok {
		if b {
			*buf = args[0]
		}
	}
	return
}

func modIfThenElse(_ *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if len(args) < 0 {
		err = ErrModPoorArgs
		return
	}
	if b, ok := ConvBool(val); ok {
		if b {
			*buf = args[0]
		} else {
			*buf = args[1]
		}
	}
	return
}

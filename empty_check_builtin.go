package dyntpl

func EmptyCheckInt(_ *Ctx, val interface{}) bool {
	if i, ok := ConvInt(val); ok && i == 0 {
		return true
	}
	return false
}

func EmptyCheckUint(_ *Ctx, val interface{}) bool {
	if u, ok := ConvUint(val); ok && u == 0 {
		return true
	}
	return false
}

func EmptyCheckFloat(_ *Ctx, val interface{}) bool {
	if f, ok := ConvFloat(val); ok && f == 0 {
		return true
	}
	return false
}

func EmptyCheckBytes(_ *Ctx, val interface{}) bool {
	if b, ok := ConvBytes(val); ok && len(b) == 0 {
		return true
	}
	return false
}

func EmptyCheckBytesSlice(_ *Ctx, val interface{}) bool {
	if b, ok := ConvBytesSlice(val); ok && len(b) == 0 {
		return true
	}
	return false
}

func EmptyCheckStr(_ *Ctx, val interface{}) bool {
	if s, ok := ConvStr(val); ok && len(s) == 0 {
		return true
	}
	return false
}

func EmptyCheckStrSlice(_ *Ctx, val interface{}) bool {
	if s, ok := ConvStrSlice(val); ok && len(s) == 0 {
		return true
	}
	return false
}

func EmptyCheckBool(_ *Ctx, val interface{}) bool {
	if b, ok := ConvBool(val); ok && !b {
		return true
	}
	return false
}

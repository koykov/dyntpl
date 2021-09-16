package dyntpl

// EmptyCheckInt checks is val is an empty integer.
func EmptyCheckInt(_ *Ctx, val interface{}) bool {
	if i, ok := ConvInt(val); ok && i == 0 {
		return true
	}
	return false
}

// EmptyCheckUint checks is val is an empty unsigned integer.
func EmptyCheckUint(_ *Ctx, val interface{}) bool {
	if u, ok := ConvUint(val); ok && u == 0 {
		return true
	}
	return false
}

// EmptyCheckFloat checks is val is an empty float number.
func EmptyCheckFloat(_ *Ctx, val interface{}) bool {
	if f, ok := ConvFloat(val); ok && f == 0 {
		return true
	}
	return false
}

// EmptyCheckBytes checks is val is an empty bytes array.
func EmptyCheckBytes(_ *Ctx, val interface{}) bool {
	if b, ok := ConvBytes(val); ok && len(b) == 0 {
		return true
	}
	return false
}

// EmptyCheckBytesSlice checks is val is an empty slice of bytes.
func EmptyCheckBytesSlice(_ *Ctx, val interface{}) bool {
	if b, ok := ConvBytesSlice(val); ok && len(b) == 0 {
		return true
	}
	return false
}

// EmptyCheckStr checks is val is an empty string.
func EmptyCheckStr(_ *Ctx, val interface{}) bool {
	if s, ok := ConvStr(val); ok && len(s) == 0 {
		return true
	}
	return false
}

// EmptyCheckStrSlice checks is val is an empty slice of strings.
func EmptyCheckStrSlice(_ *Ctx, val interface{}) bool {
	if s, ok := ConvStrSlice(val); ok && len(s) == 0 {
		return true
	}
	return false
}

// EmptyCheckBool checks is val is an empty bool.
func EmptyCheckBool(_ *Ctx, val interface{}) bool {
	if b, ok := ConvBool(val); ok && !b {
		return true
	}
	return false
}

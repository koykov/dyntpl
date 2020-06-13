package dyntpl

// Check if argument length equal zero.
func condLenEq0(_ *Ctx, args []interface{}) bool {
	if len(args) == 0 {
		return false
	}
	return getLen(args[0]) == 0
}

// Check if argument length is greater than zero.
func condLenGt0(_ *Ctx, args []interface{}) bool {
	if len(args) == 0 {
		return false
	}
	return getLen(args[0]) > 0
}

// Check if argument length is greater or equal than zero.
func condLenGtq0(_ *Ctx, args []interface{}) bool {
	if len(args) == 0 {
		return false
	}
	return getLen(args[0]) >= 0
}

// Get length of argument if it is a string or bytes.
func getLen(val interface{}) int {
	if b, ok := ConvBytes(val); ok {
		return len(b)
	}
	if s, ok := ConvStr(val); ok {
		return len(s)
	}
	return 0
}

package dyntpl

func condLenEq0(_ *Ctx, args []interface{}) bool {
	if len(args) == 0 {
		return false
	}
	return getLen(args[0]) == 0
}

func condLenGt0(_ *Ctx, args []interface{}) bool {
	if len(args) == 0 {
		return false
	}
	return getLen(args[0]) > 0
}

func condLenGtq0(_ *Ctx, args []interface{}) bool {
	if len(args) == 0 {
		return false
	}
	return getLen(args[0]) >= 0
}

func getLen(val interface{}) int {
	if b, ok := ConvBytes(val); ok {
		return len(b)
	}
	if s, ok := ConvStr(val); ok {
		return len(s)
	}
	return 0
}

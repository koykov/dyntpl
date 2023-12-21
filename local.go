package dyntpl

// Local describes value of local ctx variable.
type Local any

type localDB struct {
	idx map[string]int
	buf []Local
}

func (l *localDB) init() {
	if l.idx == nil {
		l.idx = make(map[string]int)
	}
}

// SetLocal registers new local variable.
func (l *localDB) SetLocal(name string, val Local) {
	l.init()
	l.buf = append(l.buf, val)
	l.idx[name] = len(l.buf) - 1
}

// GetLocal returns local variable from ctx.
func (l *localDB) GetLocal(name string) Local {
	l.init()
	if idx, ok := l.idx[name]; ok && idx >= 0 && idx < len(l.buf) {
		return l.buf[idx]
	}
	return nil
}

// WriteLocal writes value of local variable to dst.
func (l *localDB) WriteLocal(dst *Local, name string) {
	l.init()
	if idx, ok := l.idx[name]; ok && idx >= 0 && idx < len(l.buf) {
		*dst = l.buf[idx]
	}
}

func (l *localDB) reset() {
	l.buf = l.buf[:0]
	for k := range l.idx {
		delete(l.idx, k)
	}
}

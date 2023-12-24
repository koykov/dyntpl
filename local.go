package dyntpl

// Local describes value of local ctx variable.
type Local any

type localDB struct {
	idx map[string]int
	buf []Local
	ln  int
}

func (l *localDB) init() {
	if l.idx == nil {
		l.idx = make(map[string]int)
	}
}

// SetLocal registers new local variable.
func (l *localDB) SetLocal(name string, val Local) {
	l.init()
	if idx, ok := l.idx[name]; ok && idx >= 0 && idx < len(l.buf) {
		l.buf[idx] = val
		return
	}
	if l.ln < len(l.buf) {
		l.buf[l.ln] = val
		l.idx[name] = l.ln
	} else {
		l.buf = append(l.buf, val)
		l.idx[name] = len(l.buf) - 1
	}
	l.ln++
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
	l.ln = 0
	for k := range l.idx {
		delete(l.idx, k)
	}
}

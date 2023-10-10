package dyntpl

// Pool represents internal pool.
// In addition to native sync.Pool requires Reset() method.
type Pool interface {
	Get() any
	Put(any)
	// Reset cleanups data before putting to the pool.
	Reset(any)
}

type ipools struct {
	index map[string]int
	buf   []Pool
}

type ipoolVar struct {
	key string
	val any
}

var ipoolRegistry ipools

func (p *ipools) init() {
	if p.index == nil {
		p.index = make(map[string]int)
	}
}

func (p *ipools) acquire(key string) (any, error) {
	ipoolRegistry.init()
	i, ok := p.index[key]
	if !ok {
		return nil, ErrUnknownPool
	}
	return p.buf[i].Get(), nil
}

func (p *ipools) release(key string, x any) error {
	ipoolRegistry.init()
	i, ok := p.index[key]
	if !ok {
		return ErrUnknownPool
	}
	p.buf[i].Reset(x)
	p.buf[i].Put(x)
	return nil
}

// RegisterPool adds new internal pool to the registry by given key.
func RegisterPool(key string, pool Pool) error {
	ipoolRegistry.init()
	if _, ok := ipoolRegistry.index[key]; ok {
		return nil
	}
	ipoolRegistry.buf = append(ipoolRegistry.buf, pool)
	ipoolRegistry.index[key] = len(ipoolRegistry.buf) - 1
	return nil
}

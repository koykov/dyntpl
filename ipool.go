package dyntpl

type Pool interface {
	Get() any
	Reset(any)
	Put(any)
}

type ipools struct {
	index map[string]int
	buf   []Pool
}

type ipoolVar struct {
	key string
	val any
}

var ipools_ ipools

func (p *ipools) init() {
	if p.index == nil {
		p.index = make(map[string]int)
	}
}

func (p *ipools) acquire(key string) (any, error) {
	ipools_.init()
	i, ok := p.index[key]
	if !ok {
		return nil, ErrUnknownPool
	}
	return p.buf[i].Get(), nil
}

func (p *ipools) release(key string, x any) error {
	ipools_.init()
	i, ok := p.index[key]
	if !ok {
		return ErrUnknownPool
	}
	p.buf[i].Reset(x)
	p.buf[i].Put(x)
	return nil
}

func RegisterPool(key string, pool Pool) error {
	ipools_.init()
	if _, ok := ipools_.index[key]; ok {
		return nil
	}
	ipools_.buf = append(ipools_.buf, pool)
	ipools_.index[key] = len(ipools_.buf) - 1
	return nil
}

var _ = RegisterPool

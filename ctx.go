package cbytetpl

type Ctx struct {
	vars []ctxVar
}

type ctxVar struct {
	key string
	val interface{}
	ins Inspector
}

func NewCtx() *Ctx {
	ctx := Ctx{vars: make([]ctxVar, 0)}
	return &ctx
}

func (c *Ctx) Set(key string, val interface{}) {
	c.SetWithInspector(key, val, ReflectInspector{})
}

func (c *Ctx) SetWithInspector(key string, val interface{}, ins Inspector) {
	c.vars = append(c.vars, ctxVar{
		key: key,
		val: val,
		ins: ins,
	})
}

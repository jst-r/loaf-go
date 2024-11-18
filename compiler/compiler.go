package compiler

type Compliler struct {
	locals     []*Local
	localCount int
	scopeDepth int
}

type Local struct {
	name  *Token
	depth int
}

func (c *Compliler) beginScope() {
	c.scopeDepth++
}

func (c *Compliler) endScope() {
	c.scopeDepth--
}

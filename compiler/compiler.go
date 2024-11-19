package compiler

import "errors"

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

	i := len(c.locals) - 1
	for i >= 0 {
		if c.locals[i].depth == c.scopeDepth {
			break
		}
		i--
	}
	c.locals = c.locals[:i+1]
	c.localCount = len(c.locals)
}

func (c *Compliler) addLocal(name *Token) error {
	if c.localCount >= 255 {
		return errors.New("too many local variables")
	}

	c.locals = append(c.locals, &Local{name, c.scopeDepth})
	c.localCount++
	return nil
}

func (c *Compliler) resolveLocal(name *Token) int {
	for i := len(c.locals) - 1; i >= 0; i-- {
		if c.locals[i].name.Lexeme == name.Lexeme {
			return i
		}
	}
	return -1
}

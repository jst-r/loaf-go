package compiler

import "testing"

func TestLocals(t *testing.T) {
	c := Compliler{locals: make([]*Local, 0)}

	c.beginScope()
	c.addLocal(nil)
	c.addLocal(nil)
	c.beginScope()
	c.addLocal(nil)
	c.addLocal(nil)
	if len(c.locals) != 4 {
		t.Errorf("Expected 4 locals, got %d", len(c.locals))
	}
	c.endScope()
	if len(c.locals) != 2 {
		t.Errorf("Expected 2 locals, got %d", len(c.locals))
	}
	c.endScope()
	if len(c.locals) != 0 {
		t.Errorf("Expected 0 locals, got %d", len(c.locals))
	}
}

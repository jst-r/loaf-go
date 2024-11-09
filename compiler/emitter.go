package compiler

import (
	"github.com/jst-r/loaf-go/bytecode"
	"github.com/jst-r/loaf-go/value"
)

func (p *Parser) emitConstant(v value.Value) {
	p.emityBytes(bytecode.OpConstant, p.makeConstant(v))
}

func (p *Parser) makeConstant(v value.Value) uint8 {
	constant := p.currentChunk().AddConstant(v)
	if constant > 255 {
		p.error("Too many constants")
		return 0
	}
	return uint8(constant)
}

func (p *Parser) emitByte(b uint8) {
	p.currentChunk().Write(b, p.previous.Line)
}

func (p *Parser) emityBytes(bs ...uint8) {
	for _, b := range bs {
		p.emitByte(b)
	}
}

func (p *Parser) emitReturn() {
	p.emitByte(bytecode.OpReturn)
}

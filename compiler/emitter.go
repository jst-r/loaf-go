package compiler

import (
	"encoding/binary"
	"math"

	"github.com/jst-r/loaf-go/bytecode"
	"github.com/jst-r/loaf-go/value"
)

func (p *Parser) emitConstant(v value.Value) {
	p.emitBytes(bytecode.OpConstant, p.makeConstant(v))
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

func (p *Parser) emitBytes(bs ...uint8) {
	for _, b := range bs {
		p.emitByte(b)
	}
}

func (p *Parser) emitReturn() {
	p.emitByte(bytecode.OpReturn)
}

func (p *Parser) emitJump(op uint8) int {
	p.emitBytes(op, 0xff, 0xff)
	return len(p.currentChunk().Code) - 2
}

func (p *Parser) patchJump(offset int) {
	jump := len(p.currentChunk().Code) - offset - 2
	if jump > int(math.MaxUint16) {
		p.error("Too much code to jump over")
	}
	binary.LittleEndian.PutUint16(p.currentChunk().Code[offset:offset+2], uint16(jump))
}

func (p *Parser) emitLoop(start int) {
	p.emitByte(bytecode.OpLoop)
	offset := len(p.currentChunk().Code) - start + 2
	if offset > int(math.MaxUint16) {
		p.error("Loop body too large")
	}
	p.emitByte(byte(offset & 0xff))
	p.emitByte(byte((offset >> 8) & 0xff))
}

package compiler

import "github.com/jst-r/loaf-go/bytecode"

func (p *Parser) declaration() {
	if p.match(TokenVar) {
		p.varDeclaration()
	} else {
		p.statement()
	}

	if p.panicMode {
		p.syncronize()
	}
}

func (p *Parser) varDeclaration() {
	global := p.parseVariable("Expected variable name")

	if p.match(TokenEqual) {
		p.expression()
	} else {
		p.emitByte(bytecode.OpNil) // var x; => var x = nil
	}

	p.consume(TokenSemicolon, "Expected ; after variable declaration")
	p.defineVariable(global)
}

func (p *Parser) statement() {
	if p.match(TokenIf) {
		p.ifStatement()
	} else if p.match(TokenWhile) {
		p.whileStatement()
	} else if p.match(TokenPrint) {
		p.printStatement()
	} else if p.match(TokenLeftBrace) {
		p.compiler.beginScope()
		p.block()
		p.compiler.endScope()
	} else {
		p.expressionStatement()
	}
}

func (p *Parser) whileStatement() {
	loopStart := len(p.currentChunk().Code)
	p.expression() // condition
	p.consume(TokenLeftBrace, "Expected { after while condition")
	exitJump := p.emitJump(bytecode.OpJumpIfFalse)
	p.emitByte(bytecode.OpPop)

	p.compiler.beginScope()
	p.block()
	p.compiler.endScope()

	p.emitLoop(loopStart)

	p.patchJump(exitJump)
	p.emitByte(bytecode.OpPop)
}

func (p *Parser) ifStatement() {
	p.expression()
	p.consume(TokenLeftBrace, "Expected { after if condition")

	thenJump := p.emitJump(bytecode.OpJumpIfFalse)
	p.emitByte(bytecode.OpPop)
	p.compiler.beginScope()
	p.block()
	p.compiler.endScope()

	elseJump := p.emitJump(bytecode.OpJump)
	p.patchJump(thenJump)
	p.emitByte(bytecode.OpPop)

	if p.match(TokenElse) {
		p.consume(TokenLeftBrace, "Expected { after else")
		p.compiler.beginScope()
		p.block()
		p.compiler.endScope()

	}
	p.patchJump(elseJump)
}

func (p *Parser) block() {
	for !p.check(TokenRightBrace) && !p.check(TokenEof) {
		p.declaration()
	}
	p.consume(TokenRightBrace, "Expected } after block")
}

func (p *Parser) printStatement() {
	p.expression()
	p.consume(TokenSemicolon, "Expected ; after print statement")
	p.emitByte(bytecode.OpPrint)
}

func (p *Parser) expressionStatement() {
	p.expression()
	p.consume(TokenSemicolon, "Expected ; after expression")
	p.emitByte(bytecode.OpPop)
}

func (p *Parser) expression() {
	p.parsePrecedence(PrecedenceAssignment)
}

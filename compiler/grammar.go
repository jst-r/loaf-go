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
	if p.match(TokenPrint) {
		p.printStatement()
	} else if p.match(TokenLeftBrace) {
		p.compiler.beginScope()
		p.block()
		p.compiler.endScope()
	} else {
		p.expressionStatement()
	}
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

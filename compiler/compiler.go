package compiler

import (
	"fmt"
	"strings"

	"github.com/jst-r/loaf-go/bytecode"
)

func Compile(source string) (chunk *bytecode.Chunk, errors []string) {
	parser := NewParser(source)
	parser.advance()

	for !parser.match(TokenEof) {
		parser.declaration()
	}
	parser.endCompiler()
	return parser.compilingChunk, parser.errors
}

type Compliler struct {
	locals     []*Local
	localCount int
	scopeDepth int
}

type Local struct {
	name  *Token
	depth int
}

type Parser struct {
	compilingChunk *bytecode.Chunk
	scanner        *Scanner
	current        Token
	previous       Token
	errors         []string
	panicMode      bool
	rules          []ParseRule
	complier       Compliler
}

func NewParser(input string) *Parser {
	p := &Parser{
		compilingChunk: bytecode.NewChunk(),
		scanner:        NewScanner(input),
		current:        Token{}, previous: Token{},
		errors: nil, panicMode: false,
		complier: Compliler{
			locals: make([]*Local, 0),
		},
	}

	p.initRules()

	return p
}

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

func (p *Parser) defineVariable(global uint8) {
	p.emitBytes(bytecode.OpDefineGlobal, global)
}

func (p *Parser) statement() {
	if p.match(TokenPrint) {
		p.printStatement()
	} else if p.match(TokenRightBrace) {
		p.compiler.beginScope()
		p.block()
		p.compiler.endScope()
	} else {
		p.expressionStatement()
	}
}

func (p *Parser) block() {
	for !p.check(TokenRightBrace) && !p.check(TokenEof) {
		p.statement()
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

func (p *Parser) parseVariable(errorMessage string) uint8 {
	p.consume(TokenIdentifier, errorMessage)
	return p.identifierConstant(&p.previous)
}

func (p *Parser) identifierConstant(token *Token) uint8 {
	index := p.currentChunk().AddConstant(p.compilingChunk.Objects.NewString(strings.Clone(token.Lexeme)))
	return uint8(index)
}

func (p *Parser) endCompiler() {
	p.emitReturn()
}

func (p *Parser) currentChunk() *bytecode.Chunk {
	return p.compilingChunk
}

func (p *Parser) advance() {
	p.previous = p.current

	for {
		p.current = p.scanner.Scan()
		if p.current.Type != TokenError {
			break
		}
		p.errorAtCurrent(p.current.Lexeme)
	}
}

func (p *Parser) consume(tokenType TokenType, message string) {
	if p.current.Type == tokenType {
		p.advance()
	} else {
		p.errorAtCurrent(message)
	}
}

func (p *Parser) match(tokenType TokenType) bool {
	if !p.check(tokenType) {
		return false
	}
	p.advance()
	return true
}

func (p *Parser) check(tokenType TokenType) bool {
	return p.current.Type == tokenType
}

func (p *Parser) errorAtCurrent(msg string) {
	p.errorAt(p.current, msg)
}

func (p *Parser) error(msg string) {
	p.errorAt(p.previous, msg)
}

func (p *Parser) errorAt(token Token, msg string) {
	if p.panicMode {
		return
	}
	p.panicMode = true

	var loc string
	switch token.Type {
	case TokenEof:
		loc = " at end"
	case TokenError:
		loc = ""
	default:
		loc = fmt.Sprintf(" at %s", token.Lexeme)
	}
	err := fmt.Sprintf("[line %d]error%s: %s", token.Line, loc, msg)
	p.errors = append(p.errors, err)
}

func (p *Parser) hadError() bool {
	return len(p.errors) > 0
}

func (p *Parser) syncronize() {
	p.panicMode = false
	for p.current.Type != TokenEof {
		if p.previous.Type == TokenSemicolon {
			return
		}
		switch p.current.Type {
		case TokenClass, TokenFun, TokenVar, TokenFor, TokenIf, TokenWhile, TokenPrint, TokenReturn:
			return
		}
		p.advance()
	}
}

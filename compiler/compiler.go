package compiler

import (
	"fmt"
	"strconv"

	"github.com/jst-r/loaf-go/bytecode"
	"github.com/jst-r/loaf-go/value"
)

func Compile(source string) (chunk *bytecode.Chunk, errors []string) {
	parser := NewParser(source)
	parser.advance()
	parser.expression()
	parser.consume(TokenEof, "Expect end of input")
	parser.endCompiler()
	return parser.compilingChunk, parser.errors
}

type Parser struct {
	compilingChunk *bytecode.Chunk
	scanner        *Scanner
	current        Token
	previous       Token
	errors         []string
	panicMode      bool
}

func NewParser(input string) *Parser {
	return &Parser{
		compilingChunk: &bytecode.Chunk{},
		scanner:        NewScanner(input),
		current:        Token{}, previous: Token{},
		errors: nil, panicMode: false}
}

func (p *Parser) expression() {
	p.parsePrecedence(PrecedenceAssignment)
}

type Precedence int

const (
	PrecedenceNone Precedence = iota
	PrecedenceAssignment
	PrecedenceOr
	PrecedenceAnd
	PrecedenceComparison
	PrecedenceTerm
	PrecedenceFactor
	PrecedenceUnary
	PrecedenceCall
	PrecedencePrimary
)

type ParseRule struct {
	precedence Precedence
}

func (p *Parser) getRule(tokenType TokenType) ParseRule {
	return ParseRule{PrecedenceNone}
}

func (p *Parser) parsePrecedence(precedence Precedence) {}

func (p *Parser) number() {
	v, err := strconv.ParseFloat(p.previous.Lexeme, 64)
	if err != nil {
		p.error(err.Error())
		return
	}
	p.emitConstant(value.Float(v))
}

func (p *Parser) grouping() {
	p.expression()
	p.consume(TokenRightParen, "Expected ) after expression")
}

func (p *Parser) unary() {
	operatorType := p.previous.Type

	p.parsePrecedence(PrecedenceUnary) // compile operand first because of how the stack works

	switch operatorType {
	case TokenMinus:
		p.emitByte(bytecode.OpNegate)
	default:
		panic("unreachable case in unary")
	}
}

func (p *Parser) binary() {
	operatorType := p.previous.Type

	rule := p.getRule(operatorType)
	p.parsePrecedence(rule.precedence)

	switch operatorType {
	case TokenPlus:
		p.emitByte(bytecode.OpAdd)
	case TokenMinus:
		p.emitByte(bytecode.OpSubtract)
	case TokenStar:
		p.emitByte(bytecode.OpMultiply)
	case TokenSlash:
		p.emitByte(bytecode.OpDivide)
	default:
		panic("unreachable case in binary")
	}

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

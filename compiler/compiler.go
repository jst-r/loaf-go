package compiler

import (
	"fmt"

	"github.com/jst-r/loaf-go/bytecode"
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
	rules          []ParseRule
}

func NewParser(input string) *Parser {
	p := &Parser{
		compilingChunk: bytecode.NewChunk(),
		scanner:        NewScanner(input),
		current:        Token{}, previous: Token{},
		errors: nil, panicMode: false}

	p.initRules()

	return p
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

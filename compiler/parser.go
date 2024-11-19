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

type Parser struct {
	compilingChunk *bytecode.Chunk
	scanner        *Scanner
	current        Token
	previous       Token
	errors         []string
	panicMode      bool
	rules          []ParseRule
	compiler       Compliler
}

func NewParser(input string) *Parser {
	p := &Parser{
		compilingChunk: bytecode.NewChunk(),
		scanner:        NewScanner(input),
		current:        Token{}, previous: Token{},
		errors: nil, panicMode: false,
		compiler: Compliler{
			locals: make([]*Local, 0),
		},
	}

	p.initRules()

	return p
}

func (p *Parser) defineVariable(global uint8) {
	if p.compiler.scopeDepth > 0 {
		p.compiler.markInitialized()
		return
	}

	p.emitBytes(bytecode.OpDefineGlobal, global)
}

func (p *Parser) parseVariable(errorMessage string) uint8 {
	p.consume(TokenIdentifier, errorMessage)

	p.declareVariable()
	if p.compiler.scopeDepth > 0 {
		return 0
	}

	return p.identifierConstant(&p.previous)
}

func (p *Parser) identifierConstant(token *Token) uint8 {
	index := p.currentChunk().AddConstant(p.compilingChunk.Objects.NewString(strings.Clone(token.Lexeme)))
	return uint8(index)
}

func (p *Parser) declareVariable() {
	if p.compiler.scopeDepth == 0 {
		return
	}
	name := &p.previous

	for i := len(p.compiler.locals) - 1; i >= 0; i-- {
		local := p.compiler.locals[i]
		if local.depth != -1 && local.depth < p.compiler.scopeDepth {
			break
		}

		if local.name.Lexeme == name.Lexeme {
			p.errorAt(*name, "Variable with this name already declared in this scope")
			return
		}
	}

	err := p.compiler.addLocal(name)
	if err != nil {
		p.errorAt(*name, err.Error())
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

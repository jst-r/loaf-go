package compiler

import (
	"strconv"
	"strings"

	"github.com/jst-r/loaf-go/bytecode"
	"github.com/jst-r/loaf-go/value"
)

type Precedence int

const (
	PrecedenceNone Precedence = iota
	PrecedenceAssignment
	PrecedenceOr
	PrecedenceAnd
	PrecedenceEquality
	PrecedenceComparison
	PrecedenceTerm
	PrecedenceFactor
	PrecedenceUnary
	PrecedenceCall
	PrecedencePrimary
)

type ParseFun func()

type ParseRule struct {
	prefix     ParseFun
	infix      ParseFun
	precedence Precedence
}

func (p *Parser) initRules() {
	p.rules = make([]ParseRule, TokenEof+1)

	p.rules[TokenLeftParen] = ParseRule{prefix: p.grouping, precedence: PrecedenceNone}
	p.rules[TokenMinus] = ParseRule{prefix: p.unary, infix: p.binary, precedence: PrecedenceTerm}
	p.rules[TokenPlus] = ParseRule{infix: p.binary, precedence: PrecedenceTerm}
	p.rules[TokenStar] = ParseRule{infix: p.binary, precedence: PrecedenceFactor}
	p.rules[TokenSlash] = ParseRule{infix: p.binary, precedence: PrecedenceFactor}

	p.rules[TokenNil] = ParseRule{prefix: p.literal, precedence: PrecedenceNone}
	p.rules[TokenTrue] = ParseRule{prefix: p.literal, precedence: PrecedenceNone}
	p.rules[TokenFalse] = ParseRule{prefix: p.literal, precedence: PrecedenceNone}
	p.rules[TokenNumber] = ParseRule{prefix: p.number, precedence: PrecedenceNone}
	p.rules[TokenString] = ParseRule{prefix: p.string, precedence: PrecedenceNone}

	p.rules[TokenBang] = ParseRule{prefix: p.unary, precedence: PrecedenceNone}
	p.rules[TokenEqualEqual] = ParseRule{infix: p.binary, precedence: PrecedenceEquality}
	p.rules[TokenBangEqual] = ParseRule{infix: p.binary, precedence: PrecedenceEquality}
	p.rules[TokenLess] = ParseRule{infix: p.binary, precedence: PrecedenceComparison}
	p.rules[TokenLessEqual] = ParseRule{infix: p.binary, precedence: PrecedenceComparison}
	p.rules[TokenGreater] = ParseRule{infix: p.binary, precedence: PrecedenceComparison}
	p.rules[TokenGreaterEqual] = ParseRule{infix: p.binary, precedence: PrecedenceComparison}
}

func (p *Parser) declaration() {
	p.statement()

	if p.panicMode {
		p.syncronize()
	}
}

func (p *Parser) statement() {
	if p.match(TokenPrint) {
		p.printStatement()
	} else {
		p.expressionStatement()
	}
}

func (p *Parser) printStatement() {
	p.expression()
	p.consume(TokenSemicolon, "Expected ; after print statement")
	p.emitByte(bytecode.OpPrint)
	p.emitByte(bytecode.OpPop)
}

func (p *Parser) expressionStatement() {
	p.expression()
	p.consume(TokenSemicolon, "Expected ; after expression")
}

func (p *Parser) expression() {
	p.parsePrecedence(PrecedenceAssignment)
}

func (p *Parser) getRule(tokenType TokenType) *ParseRule {
	return &p.rules[tokenType]
}

func (p *Parser) parsePrecedence(precedence Precedence) {
	p.advance()
	rule := p.getRule(p.previous.Type)
	if rule.prefix == nil {
		p.error("Expected expression")
		return
	}

	rule.prefix()

	for precedence <= p.getRule(p.current.Type).precedence {
		p.advance()
		rule = p.getRule(p.previous.Type)
		if rule.infix == nil {
			p.error("Expected infix operation")
			return
		}
		rule.infix()
	}
}

func (p *Parser) number() {
	v, err := strconv.ParseFloat(p.previous.Lexeme, 64)
	if err != nil {
		p.error(err.Error())
		return
	}
	p.emitConstant(value.Float(v))
}

func (p *Parser) string() {
	// Without the clone this will point to the source file, which is a hassle to deal with
	v := strings.Clone(p.previous.Lexeme[1 : len(p.previous.Lexeme)-1])
	p.emitConstant(p.compilingChunk.Objects.NewString(v))
}

func (p *Parser) grouping() {
	p.expression()
	p.consume(TokenRightParen, "Expected ) after expression")
}

func (p *Parser) unary() {
	operatorType := p.previous.Type

	p.parsePrecedence(PrecedenceUnary) // compile operand first because of how the stack works

	switch operatorType {
	case TokenBang:
		p.emitByte(bytecode.OpNot)
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
	case TokenEqualEqual:
		p.emitByte(bytecode.OpEqual)
	case TokenBangEqual:
		p.emitBytes(bytecode.OpEqual, bytecode.OpNot)
	case TokenGreater:
		p.emitByte(bytecode.OpGreater)
	case TokenGreaterEqual:
		p.emitBytes(bytecode.OpGreater, bytecode.OpNot)
	case TokenLess:
		p.emitByte(bytecode.OpLess)
	case TokenLessEqual:
		p.emitBytes(bytecode.OpLess, bytecode.OpNot)
	default:
		panic("unreachable case in binary")
	}

}

func (p *Parser) literal() {
	switch p.previous.Type {
	case TokenNil:
		p.emitByte(bytecode.OpNil)
	case TokenTrue:
		p.emitByte(bytecode.OpTrue)
	case TokenFalse:
		p.emitByte(bytecode.OpFalse)
	default:
		panic("unreachable case in literal")
	}
}

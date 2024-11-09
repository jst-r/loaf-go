package compiler

import (
	"strconv"

	"github.com/jst-r/loaf-go/bytecode"
	"github.com/jst-r/loaf-go/value"
)

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
	p.rules[TokenNumber] = ParseRule{prefix: p.number, precedence: PrecedenceNone}
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

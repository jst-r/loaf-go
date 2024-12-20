package compiler

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
}

type TokenType int

const (
	// Single-character tokens
	TokenLeftParen TokenType = iota
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenComma
	TokenDot
	TokenMinus
	TokenPlus
	TokenSemicolon
	TokenSlash
	TokenStar
	// One or two character tokens
	TokenBang
	TokenBangEqual
	TokenEqual
	TokenEqualEqual
	TokenGreater
	TokenGreaterEqual
	TokenLess
	TokenLessEqual
	TokenPlusEqual
	TokenMinusEqual
	// Literals
	TokenIdentifier
	TokenString
	TokenNumber
	// Keywords
	TokenAnd
	TokenClass
	TokenElse
	TokenFalse
	TokenFor
	TokenFun
	TokenIf
	TokenNil
	TokenOr
	TokenPrint
	TokenReturn
	TokenSuper
	TokenThis
	TokenTrue
	TokenVar
	TokenWhile
	// Special
	TokenError
	TokenEof
)

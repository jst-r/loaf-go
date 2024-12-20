package compiler

type Scanner struct {
	input   string
	current int
	line    int
}

func NewScanner(input string) *Scanner {
	return &Scanner{input, 0, 1}
}

func (s *Scanner) Scan() (token Token) {
	s.skipWhitespace()
	s.input = s.input[s.current:]
	s.current = 0

	if len(s.input) == 0 {
		return Token{TokenEof, "", s.line}
	}

	c := s.advance()

	if isAlpha(c) {
		return s.identifier()
	}
	if isDigit(c) {
		return s.number()
	}

	switch c {
	case '(':
		return s.makeToken(TokenLeftParen)
	case ')':
		return s.makeToken(TokenRightParen)
	case '{':
		return s.makeToken(TokenLeftBrace)
	case '}':
		return s.makeToken(TokenRightBrace)
	case ',':
		return s.makeToken(TokenComma)
	case '.':
		return s.makeToken(TokenDot)
	case '-':
		if s.match('=') {
			return s.makeToken(TokenMinusEqual)
		}
		return s.makeToken(TokenMinus)
	case '+':
		if s.match('=') {
			return s.makeToken(TokenPlusEqual)
		}
		return s.makeToken(TokenPlus)
	case ';':
		return s.makeToken(TokenSemicolon)
	case '/':
		return s.makeToken(TokenSlash)
	case '*':
		return s.makeToken(TokenStar)
	case '!':
		if s.match('=') {
			return s.makeToken(TokenBangEqual)
		}
		return s.makeToken(TokenBang)
	case '=':
		if s.match('=') {
			return s.makeToken(TokenEqualEqual)
		}
		return s.makeToken(TokenEqual)
	case '<':
		if s.match('=') {
			return s.makeToken(TokenLessEqual)
		}
		return s.makeToken(TokenLess)
	case '>':
		if s.match('=') {
			return s.makeToken(TokenGreaterEqual)
		}
		return s.makeToken(TokenGreater)
	case '"':
		return s.string()
	case 0:
		return s.makeToken(TokenEof)
	}

	return s.errorToken("unexpected character")
}

func (s *Scanner) string() Token {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return s.errorToken("unterminated string")
	}
	s.advance()
	return s.makeToken(TokenString)
}

func (s *Scanner) number() Token {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	return s.makeToken(TokenNumber)
}

func (s *Scanner) identifier() Token {
	for isAlpha(s.peek()) || isDigit(s.peek()) {
		s.advance()
	}
	return s.makeToken(s.identifierType())
}

func (s *Scanner) identifierType() TokenType {
	switch s.input[:s.current] {
	case "and":
		return TokenAnd
	case "class":
		return TokenClass
	case "else":
		return TokenElse
	case "false":
		return TokenFalse
	case "for":
		return TokenFor
	case "fun":
		return TokenFun
	case "if":
		return TokenIf
	case "nil":
		return TokenNil
	case "or":
		return TokenOr
	case "print":
		return TokenPrint
	case "return":
		return TokenReturn
	case "super":
		return TokenSuper
	case "this":
		return TokenThis
	case "true":
		return TokenTrue
	case "var":
		return TokenVar
	case "while":
		return TokenWhile
	default:
		return TokenIdentifier
	}
}

func (s *Scanner) errorToken(msg string) Token {
	return Token{TokenError, msg, s.line}
}

func (s *Scanner) makeToken(tokenType TokenType) Token {
	return Token{tokenType, s.input[:s.current], s.line}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.input)
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.input[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.input[s.current+1]
}

func (s *Scanner) advance() byte {
	if s.isAtEnd() {
		return 0
	}
	s.current++
	return s.input[s.current-1]
}

func (s *Scanner) match(c byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.peek() == c {
		s.current++
		return true
	}
	return false
}

func (s *Scanner) skipWhitespace() {
	for {
		c := s.peek()
		switch c {
		case ' ', '\t', '\r':
			s.advance()
		case '\n':
			s.line++
			s.advance()
		case '/':
			if s.peekNext() != '/' {
				return
			}
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		default:
			return
		}
	}
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

package main

import (
	"strconv"
)

var keywords = map[string]TokenType{
	"and":   AND,
	"class": CLASS,
	"else":  ELSE,
	"false": FALSE,
	"for":   FOR,
	"fun":   FUN,
	"if":    IF,
	"nil":   NIL,
	"or":    OR,
	"print": PRINT,
	"return":RETURN,
	"super": SUPER,
	"this":  THIS,
	"true":  TRUE,
	"var":   VAR,
	"while": WHILE,
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func (s *Scanner) Scanner(source string) *Scanner {
	return &Scanner{
		source: source,
		line:   1,
	}
}

func (s *Scanner) scanTokens() []Token {
	// Does this really work?

	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(
		s.tokens,
		NewToken(EOF, "", nil, s.line),
	)
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func trn[T any](b bool, el1, el2 T) T {
	if b {
		return el1
	}
	return el2
}

func (s *Scanner) scanToken() {
	var c byte = s.advance()

	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		s.addToken(trn(s.match('='), BANG_EQUAL, BANG))
	case '=':
		s.addToken(trn(s.match('='), EQUAL_EQUAL, EQUAL))
	case '<':
		s.addToken(trn(s.match('='), LESS_EQUAL, LESS))
	case '>':
		s.addToken(trn(s.match('='), GREATER_EQUAL, GREATER))

	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}

	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		s.line++

	case '"':
		s.string()
		break

	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			loxerror(s.line, "Unexpected character.")
		}
		break

	}
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
    tokentype, ok := keywords[text];
    if !ok {
		tokentype = IDENTIFIER;
	}

    s.addToken(tokentype);

	s.addToken(IDENTIFIER)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		loxerror(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[(s.start + 1):(s.current - 1)]
	s.addToken2(STRING, value)
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current += 1
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func isAlpha(c byte) bool {
    return (c >= 'a' && c <= 'z') ||
           (c >= 'A' && c <= 'Z') ||
            c == '_';
  }

func isAlphaNumeric(c byte) bool {
    return isAlpha(c) || isDigit(c);
  }

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	fl, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		panic(err)
	}

	s.addToken2(NUMBER, fl)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) advance() byte {
	res := s.source[s.current]
	s.current += 1
	return res
}

func (s *Scanner) addToken(tokentype TokenType) {
	s.addToken2(tokentype, nil)
}

func (s *Scanner) addToken2(tokentype TokenType, literal any) {
	// TODO: Verify that this works in the same way as java substring()
	text := s.source[s.start:s.current]
	s.tokens = append(
		s.tokens,
		NewToken(tokentype, text, literal, s.line),
	)

}

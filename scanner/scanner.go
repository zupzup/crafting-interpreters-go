package scanner

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zupzup/crafting-interpreters-go/constants"
	"github.com/zupzup/crafting-interpreters-go/token"
	"strconv"
)

// Scanner scans tokens
type Scanner struct {
	Source  string
	Tokens  []token.Token
	start   int
	current int
	line    int
}

// ScanTokens returns a list of tokens or an error
func (s *Scanner) ScanTokens() ([]token.Token, error) {
	hadError := false
	s.Tokens = []token.Token{}
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			hadError = true
		}
	}
	s.Tokens = append(s.Tokens, token.Token{
		TokenType: constants.EOF,
		Lexeme:    "",
		Literal:   nil,
		Line:      s.line,
	})
	if hadError {
		return nil, errors.New("error")
	}
	return s.Tokens, nil
}

func (s *Scanner) scanToken() error {
	b := s.advance()
	switch b {
	case '(':
		s.addToken(constants.LeftParen, nil)
	case ')':
		s.addToken(constants.RightParen, nil)
	case '{':
		s.addToken(constants.LeftBrace, nil)
	case '}':
		s.addToken(constants.RightBrace, nil)
	case ',':
		s.addToken(constants.Comma, nil)
	case '.':
		s.addToken(constants.Dot, nil)
	case '-':
		s.addToken(constants.Minus, nil)
	case '+':
		s.addToken(constants.Plus, nil)
	case ';':
		s.addToken(constants.Semicolon, nil)
	case '*':
		s.addToken(constants.Star, nil)
	case '!':
		if s.match('=') {
			s.addToken(constants.BangEqual, nil)
		}
		s.addToken(constants.Bang, nil)
	case '=':
		if s.match('=') {
			s.addToken(constants.BangEqual, nil)
		}
		s.addToken(constants.Equal, nil)
	case '<':
		if s.match('=') {
			s.addToken(constants.LessEqual, nil)
		}
		s.addToken(constants.Less, nil)
	case '>':
		if s.match('=') {
			s.addToken(constants.GreateEqual, nil)
		}
		s.addToken(constants.Greater, nil)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		}
		s.addToken(constants.Slash, nil)
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line = s.line + 1
	case '"':
		err := s.string()
		if err != nil {
			logError(s.line, err.Error())
			return err
		}
	default:
		if s.isDigit(b) {
			s.number()
		} else if s.isAlpha(b) {
			s.identifier()
		} else {
			logError(s.line, "Unexpected character.")
			return errors.New("unexpected character")
		}
	}
	return nil
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line = s.line + 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		return errors.New("unterminated string")
	}

	s.advance()

	value := s.Source[s.start+1 : s.current-1]
	s.addToken(constants.Str, value)
	return nil
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
	}
	for s.isDigit(s.peek()) {
		s.advance()
	}
	v, _ := strconv.ParseFloat(s.Source[s.start:s.current], 64)
	s.addToken(constants.Number, v)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.Source[s.start:s.current]
	tokenType := constants.Keywords[text]
	fmt.Println(tokenType)
	if tokenType == 0 {
		tokenType = constants.Identifier
	}
	s.addToken(tokenType, nil)
}

func (s *Scanner) isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (s *Scanner) isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func (s *Scanner) isAlphaNumeric(b byte) bool {
	return s.isAlpha(b) || s.isDigit(b)
}

func (s *Scanner) advance() byte {
	s.current = s.current + 1
	return s.Source[s.current-1]
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.Source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.Source) {
		return '\000'
	}
	return s.Source[s.current+1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.Source[s.current] != expected {
		return false
	}
	s.current = s.current + 1
	return true
}

func (s *Scanner) addToken(tokenType int, literal interface{}) {
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, token.Token{
		TokenType: tokenType,
		Lexeme:    text,
		Literal:   literal,
		Line:      s.line,
	})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func logError(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
}

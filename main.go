package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zupzup/crafting-interpreters-go/constants"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// Scanner scans tokens
type Scanner struct {
	Source  string
	Tokens  []Token
	start   int
	current int
	line    int
}

func newScanner(source string) *Scanner {
	return &Scanner{
		Source: source,
	}
}

func (s *Scanner) scanTokens() ([]Token, error) {
	hadError := false
	s.Tokens = []Token{}
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			hadError = true
		}
	}
	s.Tokens = append(s.Tokens, Token{
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
	s.Tokens = append(s.Tokens, Token{
		TokenType: tokenType,
		Lexeme:    text,
		Literal:   literal,
		Line:      s.line,
	})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

// Token is a token
type Token struct {
	TokenType int
	Lexeme    string
	Literal   interface{}
	Line      int
}

func newToken(tokenType int, lexeme string, literal interface{}, line int) Token {
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%d %s %v", t.TokenType, t.Lexeme, t.Literal)
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [script]")
	} else if len(os.Args) > 1 {
		if err := runFile(os.Args[1]); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := runPrompt(); err != nil {
			log.Fatal(err)
		}
	}
}

func runPrompt() error {
	fmt.Println("prompt")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "could not read from stdin")
		}
		if err := run(text); err != nil {
			return err // TODO: don't kill session on one error, just log it
		}
	}
}

func runFile(path string) error {
	fmt.Println("file: " + path)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "could not read file at %s", path)
	}
	if err := run(string(bytes)); err != nil {
		return err
	}
	return nil
}

func run(code string) error {
	scanner := Scanner{Source: code}
	tokens, err := scanner.scanTokens()
	if err != nil {
		return err
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
	return nil
}

func logError(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
}

package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zupzup/crafting-interpreters-go/constants"
	"io/ioutil"
	"log"
	"os"
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

func (s *Scanner) scanTokens() []Token {
	s.Tokens = []Token{}
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.Tokens = append(s.Tokens, Token{
		TokenType: constants.EOF,
		Lexeme:    "",
		Literal:   nil,
		Line:      s.line,
	})
	return s.Tokens
}

func (s *Scanner) scanToken() {
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
	default:
		logError(s.line, "Unexpected character.")
	}
}

func (s *Scanner) advance() byte {
	s.current = s.current + 1
	return s.Source[s.current-1]
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
	scanner := Scanner{}
	tokens := scanner.scanTokens()

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

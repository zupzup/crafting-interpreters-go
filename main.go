package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zupzup/crafting-interpreters-go/ast"
	"github.com/zupzup/crafting-interpreters-go/scanner"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	blub := ast.Binary{}
	bla := ast.Unary{}
	p := ast.Printer{}
	fmt.Println(p.Print(&blub))
	fmt.Println(p.Print(&bla))
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
	scanner := scanner.Scanner{Source: code}
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return err
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
	return nil
}

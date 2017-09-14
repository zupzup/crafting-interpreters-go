package ast

import (
	"github.com/zupzup/crafting-interpreters-go/token"
)

// Printer prints an Ast
type Printer struct{}

// Print prints the ast
func (p *Printer) Print(expr Expr) string {
	v, _ := expr.accept(p).(string)
	return v
}

func (p *Printer) visitBinaryExpr(expr *Binary) interface{} {
	return "Binary"
}
func (p *Printer) visitGroupingExpr(expr *Grouping) interface{} {
	return "Grouping"
}
func (p *Printer) visitUnaryExpr(expr *Unary) interface{} {
	return "Unary"
}
func (p *Printer) visitLiteralExpr(expr *Literal) interface{} {
	return "Literal"
}

// Visitor is the visitor pattern interface
type Visitor interface {
	visitBinaryExpr(expr *Binary) interface{}
	visitGroupingExpr(expr *Grouping) interface{}
	visitUnaryExpr(expr *Unary) interface{}
	visitLiteralExpr(expr *Literal) interface{}
}

// Expr is the interface all expressions implement
type Expr interface {
	exprNode()
	accept(v Visitor) interface{}
}

// Binary is a binary expression
type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *Binary) exprNode() {}
func (b *Binary) accept(v Visitor) interface{} {
	return v.visitBinaryExpr(b)
}

// Grouping is a grouping expression
type Grouping struct {
	Expression Expr
}

func (g *Grouping) exprNode() {}
func (g *Grouping) accept(v Visitor) interface{} {
	return v.visitGroupingExpr(g)
}

// Unary is a unary expression
type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u *Unary) exprNode() {}
func (u *Unary) accept(v Visitor) interface{} {
	return v.visitUnaryExpr(u)
}

// Literal is a literal expression
type Literal struct {
	Value interface{}
}

func (l *Literal) exprNode() {}
func (l *Literal) accept(v Visitor) interface{} {
	return v.visitLiteralExpr(l)
}

// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

import (
	"fmt"
	"golem/ast"
	"golem/scanner"
	"runtime"
)

//--------------------------------------------------------------
// Parser

type Parser struct {
	scn  *scanner.Scanner
	cur  *ast.Token
	next *ast.Token
}

func NewParser(scn *scanner.Scanner) *Parser {
	return &Parser{scn, nil, nil}
}

func (p *Parser) ParseModule() (fn *ast.FnExpr, err error) {

	// In a recursive descent parser, errors can be generated deep
	// in the call stack.  We are going to use panic-recover to handle them.
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			fn = nil
			err = r.(error)
		}
	}()

	// read the first two tokens
	p.cur = p.advance()
	p.next = p.advance()

	// parse the module
	nodes := p.nodeSequence(ast.EOF)

	params := []*ast.IdentExpr{}
	return &ast.FnExpr{nil, params, &ast.Block{nodes}, 0, 0, nil}, err
}

func (p *Parser) parseExpression() (expr ast.Expr, err error) {

	// In a recursive descent parser, errors can be generated deep
	// in the call stack.  We are going to use panic-recover to handle them.
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			expr = nil
			err = r.(error)
		}
	}()

	// read the first two tokens
	p.cur = p.advance()
	p.next = p.advance()

	// parse the expression
	expr = p.expression()
	p.expect(ast.EOF)

	return expr, err
}

// parse a statement, or return nil if there is no statement
// waiting to be parsed
func (p *Parser) statement() ast.Stmt {

	switch {

	case p.accept(ast.CONST):
		return p.constStmt()

	case p.accept(ast.LET):
		return p.letStmt()

	case p.accept(ast.IF):
		return p.ifStmt()

	case p.accept(ast.WHILE):
		return p.whileStmt()

	case p.accept(ast.BREAK):
		return p.breakStmt()

	case p.accept(ast.CONTINUE):
		return p.continueStmt()

	case p.accept(ast.RETURN):
		return p.returnStmt()

	default:
		return nil
	}
}

func (p *Parser) constStmt() *ast.Const {

	sym := p.expect(ast.IDENT)
	p.expect(ast.EQ)
	ident := &ast.IdentExpr{sym, nil}

	expr := p.expression()
	p.expect(ast.SEMICOLON)

	return &ast.Const{ident, expr}
}

func (p *Parser) letStmt() *ast.Let {

	sym := p.expect(ast.IDENT)
	p.expect(ast.EQ)
	ident := &ast.IdentExpr{sym, nil}

	expr := p.expression()
	p.expect(ast.SEMICOLON)

	return &ast.Let{ident, expr}
}

func (p *Parser) ifStmt() *ast.If {

	cond := p.expression()
	then := p.block()

	if p.accept(ast.ELSE) {

		switch p.cur.Kind {

		case ast.LBRACE:
			return &ast.If{cond, then, p.block()}

		case ast.IF:
			p.expect(ast.IF)
			return &ast.If{cond, then, p.ifStmt()}

		default:
			panic(p.unexpected())
		}

	} else {
		return &ast.If{cond, then, nil}
	}
}

func (p *Parser) whileStmt() *ast.While {

	return &ast.While{p.expression(), p.block()}
}

func (p *Parser) breakStmt() *ast.Break {

	p.expect(ast.SEMICOLON)
	return &ast.Break{}
}

func (p *Parser) continueStmt() *ast.Continue {

	p.expect(ast.SEMICOLON)
	return &ast.Continue{}
}

func (p *Parser) returnStmt() *ast.Return {

	if p.accept(ast.SEMICOLON) {
		return &ast.Return{nil}
	} else {
		val := p.expression()
		p.expect(ast.SEMICOLON)
		return &ast.Return{val}
	}
}

// parse a sequence of nodes that are wrapped in curly braces
func (p *Parser) block() *ast.Block {

	p.expect(ast.LBRACE)
	return &ast.Block{p.nodeSequence(ast.RBRACE)}
}

// Parse a sequence of statements or expressions.
func (p *Parser) nodeSequence(endKind ast.TokenKind) []ast.Node {

	nodes := []ast.Node{}

	for {
		if p.accept(endKind) {
			break
		}

		// see if there is a statement on tap
		var node ast.Node = p.statement()

		// if there isn't, read an expression instead
		if node == nil {
			node = p.expression()
			p.expect(ast.SEMICOLON)
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func (p *Parser) expression() ast.Expr {

	if p.cur.Kind == ast.IDENT && p.next.Kind == ast.EQ {

		sym := p.expect(ast.IDENT)
		p.expect(ast.EQ)

		ident := &ast.IdentExpr{sym, nil}
		return &ast.Assignment{ident, p.expression()}

	} else {
		lhs := p.andExpr()
		for p.cur.Kind == ast.DBL_PIPE {
			tok := p.cur
			p.consume()
			lhs = &ast.BinaryExpr{lhs, tok, p.andExpr()}
		}
		return lhs

	}
}

func (p *Parser) andExpr() ast.Expr {

	lhs := p.comparativeExpr()
	for p.cur.Kind == ast.DBL_AMP {
		tok := p.cur
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.comparativeExpr()}
	}

	return lhs
}

func (p *Parser) comparativeExpr() ast.Expr {

	lhs := p.additiveExpr()
	for isComparative(p.cur) {
		tok := p.cur
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.additiveExpr()}
	}

	return lhs
}

func (p *Parser) additiveExpr() ast.Expr {

	lhs := p.multiplicativeExpr()
	for isAdditive(p.cur) {
		tok := p.cur
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.multiplicativeExpr()}
	}

	return lhs
}

func (p *Parser) multiplicativeExpr() ast.Expr {

	lhs := p.unaryExpr()
	for isMultiplicative(p.cur) {
		tok := p.cur
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.unaryExpr()}
	}

	return lhs
}

func (p *Parser) unaryExpr() ast.Expr {

	if isUnary(p.cur) {
		tok := p.cur
		p.consume()
		return &ast.UnaryExpr{tok, p.unaryExpr()}

	} else {
		return p.primaryExpr()
	}
}

func (p *Parser) primaryExpr() ast.Expr {
	prm := p.primary()

	for {
		// look for suffixes: Invoke, Select, Index
		switch p.cur.Kind {

		case ast.LPAREN:
			actual, last := p.actualParams()
			prm = &ast.InvokeExpr{last, prm, actual}

		case ast.DOT:
			p.expect(ast.DOT)
			key := p.expect(ast.IDENT)

			// TODO: is it correct to parse PutExpr here, rather than in p.expression()?
			// Something doesn't seem quite right.
			if p.accept(ast.EQ) {
				prm = &ast.PutExpr{prm, key, p.expression()}
			} else {
				prm = &ast.SelectExpr{prm, key}
			}

		default:
			return prm
		}
	}
}

func (p *Parser) primary() ast.Expr {

	switch p.cur.Kind {

	case ast.LPAREN:
		p.consume()
		expr := p.expression()
		p.expect(ast.RPAREN)
		return expr

	case ast.IDENT:
		return p.identExpr()

	case ast.THIS:
		return &ast.ThisExpr{p.consume(), nil}

	case ast.FN:
		return p.fnExpr(p.consume())

	case ast.OBJ:
		return p.objExpr(p.consume())

	default:
		return p.literalExpr()
	}
}

func (p *Parser) identExpr() *ast.IdentExpr {
	tok := p.cur
	p.expect(ast.IDENT)
	return &ast.IdentExpr{tok, nil}
}

func (p *Parser) fnExpr(first *ast.Token) ast.Expr {

	p.expect(ast.LPAREN)

	params := []*ast.IdentExpr{}

	switch p.cur.Kind {

	case ast.IDENT:
		params = append(params, p.identExpr())
	loop:
		for {
			switch p.cur.Kind {

			case ast.COMMA:
				p.consume()
				params = append(params, p.identExpr())

			case ast.RPAREN:
				p.consume()
				break loop

			default:
				panic(p.unexpected())
			}
		}

	case ast.RPAREN:
		p.consume()

	default:
		panic(p.unexpected())
	}

	blk := p.block()
	return &ast.FnExpr{first, params, blk, 0, 0, nil}
}

func (p *Parser) objExpr(first *ast.Token) ast.Expr {

	keys := []*ast.Token{}
	values := []ast.Expr{}
	var last *ast.Token

	p.expect(ast.LBRACE)

	switch p.cur.Kind {

	case ast.IDENT:
		keys = append(keys, p.cur)
		p.consume()
		p.expect(ast.COLON)
		values = append(values, p.expression())
	loop:
		for {
			switch p.cur.Kind {

			case ast.COMMA:
				p.consume()
				keys = append(keys, p.cur)
				p.consume()
				p.expect(ast.COLON)
				values = append(values, p.expression())

			case ast.RBRACE:
				last = p.consume()
				break loop

			default:
				panic(p.unexpected())
			}
		}

	case ast.RBRACE:
		last = p.consume()

	default:
		panic(p.unexpected())
	}

	return &ast.ObjExpr{first, last, keys, values, -1}
}

func (p *Parser) literalExpr() ast.Expr {

	tok := p.cur

	switch {

	case tok.IsBasic():
		p.consume()
		return &ast.BasicExpr{tok}

	default:
		panic(p.unexpected())
	}
}

func (p *Parser) actualParams() ([]ast.Expr, *ast.Token) {

	p.expect(ast.LPAREN)

	params := []ast.Expr{}
	switch p.cur.Kind {

	case ast.RPAREN:
		last := p.consume()
		return params, last

	default:
		params = append(params, p.expression())
		for {
			switch p.cur.Kind {

			case ast.COMMA:
				p.consume()
				params = append(params, p.expression())

			case ast.RPAREN:
				last := p.consume()
				return params, last

			default:
				panic(p.unexpected())
			}

		}
	}
}

// consume the current token if it has the given kind
func (p *Parser) accept(kind ast.TokenKind) bool {
	if p.cur.Kind == kind {
		p.consume()
		return true
	} else {
		return false
	}
}

// consume the current token if it has the given kind, else panic
func (p *Parser) expect(kind ast.TokenKind) *ast.Token {
	if p.cur.Kind == kind {
		result := p.cur
		p.consume()
		return result
	} else {
		panic(p.unexpected())
	}
}

// consume the current token
func (p *Parser) consume() *ast.Token {
	result := p.cur
	p.cur, p.next = p.next, p.advance()
	return result
}

func (p *Parser) advance() *ast.Token {

	tok := p.scn.Next()
	if tok.IsBad() {
		switch tok.Kind {

		case ast.UNEXPECTED_CHAR:
			panic(&parserError{UNEXPECTED_CHAR, tok})

		case ast.UNEXPECTED_EOF:
			panic(&parserError{UNEXPECTED_EOF, tok})

		default:
			panic("unreachable")
		}
	}
	return tok
}

// create a error that we will panic with
func (p *Parser) unexpected() error {
	switch p.cur.Kind {
	case ast.EOF:
		return &parserError{UNEXPECTED_EOF, p.cur}

	default:
		return &parserError{UNEXPECTED_TOKEN, p.cur}
	}
}

func isComparative(t *ast.Token) bool {
	return t.Kind == ast.DBL_EQ || t.Kind == ast.NOT_EQ ||
		t.Kind == ast.GT || t.Kind == ast.GT_EQ ||
		t.Kind == ast.LT || t.Kind == ast.LT_EQ ||
		t.Kind == ast.CMP
}

func isAdditive(t *ast.Token) bool {
	return t.Kind == ast.PLUS || t.Kind == ast.MINUS
}

func isMultiplicative(t *ast.Token) bool {
	return t.Kind == ast.MULT || t.Kind == ast.DIV
}

func isUnary(t *ast.Token) bool {
	return t.Kind == ast.MINUS || t.Kind == ast.NOT
}

//--------------------------------------------------------------
// parserError

type parserErrorKind int

const (
	UNEXPECTED_CHAR parserErrorKind = iota
	UNEXPECTED_TOKEN
	UNEXPECTED_EOF
)

type parserError struct {
	kind  parserErrorKind
	token *ast.Token
}

func (e *parserError) Error() string {

	switch e.kind {

	case UNEXPECTED_CHAR:
		return fmt.Sprintf("Unexpected Character '%v' at %v", e.token.Text, e.token.Position)

	case UNEXPECTED_TOKEN:
		return fmt.Sprintf("Unexpected Token '%v' at %v", e.token.Text, e.token.Position)

	case UNEXPECTED_EOF:
		return fmt.Sprintf("Unexpected EOF at %v", e.token.Position)

	default:
		panic("unreachable")
	}
}

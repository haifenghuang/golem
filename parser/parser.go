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
	p.expect(ast.EOF)

	params := []*ast.IdentExpr{}
	return &ast.FnExpr{nil, params, &ast.Block{nil, nodes, nil}, 0, 0, nil}, err
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

	switch p.cur.Kind {

	case ast.CONST:
		return p.constStmt()

	case ast.LET:
		return p.letStmt()

	case ast.IF:
		return p.ifStmt()

	case ast.WHILE:
		return p.whileStmt()

	case ast.BREAK:
		return p.breakStmt()

	case ast.CONTINUE:
		return p.continueStmt()

	case ast.RETURN:
		return p.returnStmt()

	default:
		return nil
	}
}

func (p *Parser) constStmt() *ast.Const {

	token := p.expect(ast.CONST)

	decls := []*ast.Decl{p.decl()}
	for {
		switch p.cur.Kind {
		case ast.COMMA:
			p.consume()
			decls = append(decls, p.decl())
		case ast.SEMICOLON:
			return &ast.Const{token, decls, p.consume()}
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) letStmt() *ast.Let {

	token := p.expect(ast.LET)

	decls := []*ast.Decl{p.decl()}
	for {
		switch p.cur.Kind {
		case ast.COMMA:
			p.consume()
			decls = append(decls, p.decl())
		case ast.SEMICOLON:
			return &ast.Let{token, decls, p.consume()}
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) decl() *ast.Decl {

	ident := &ast.IdentExpr{p.expect(ast.IDENT), nil}
	if p.accept(ast.EQ) {
		return &ast.Decl{ident, p.expression()}
	} else {
		return &ast.Decl{ident, nil}
	}
}

func (p *Parser) ifStmt() *ast.If {

	token := p.expect(ast.IF)
	cond := p.expression()
	then := p.block()

	if p.accept(ast.ELSE) {

		switch p.cur.Kind {

		case ast.LBRACE:
			return &ast.If{token, cond, then, p.block()}

		case ast.IF:
			return &ast.If{token, cond, then, p.ifStmt()}

		default:
			panic(p.unexpected())
		}

	} else {
		return &ast.If{token, cond, then, nil}
	}
}

func (p *Parser) whileStmt() *ast.While {

	return &ast.While{p.expect(ast.WHILE), p.expression(), p.block()}
}

func (p *Parser) breakStmt() *ast.Break {
	return &ast.Break{
		p.expect(ast.BREAK),
		p.expect(ast.SEMICOLON)}
}

func (p *Parser) continueStmt() *ast.Continue {
	return &ast.Continue{
		p.expect(ast.CONTINUE),
		p.expect(ast.SEMICOLON)}
}

func (p *Parser) returnStmt() *ast.Return {

	token := p.expect(ast.RETURN)

	if p.cur.Kind == ast.SEMICOLON {
		return &ast.Return{token, nil, p.expect(ast.SEMICOLON)}
	} else {
		val := p.expression()
		return &ast.Return{token, val, p.expect(ast.SEMICOLON)}
	}
}

// parse a sequence of nodes that are wrapped in curly braces
func (p *Parser) block() *ast.Block {

	lbrace := p.expect(ast.LBRACE)
	nodes := p.nodeSequence(ast.RBRACE)
	rbrace := p.expect(ast.RBRACE)
	return &ast.Block{lbrace, nodes, rbrace}
}

// Parse a sequence of statements or expressions.
func (p *Parser) nodeSequence(endKind ast.TokenKind) []ast.Node {

	nodes := []ast.Node{}

	for {
		if p.cur.Kind == endKind {
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

	exp := p.ternaryExpr()

	if asn, ok := exp.(ast.Assignable); ok {

		if p.cur.Kind == ast.EQ {

			// assignment
			eq := p.expect(ast.EQ)
			exp = &ast.Assignment{asn, eq, p.expression()}

		} else if isAssignOp(p.cur) {

			// assignment operation
			op := p.consume()
			exp = &ast.Assignment{
				asn,
				op,
				&ast.BinaryExpr{
					asn,
					fromAssignOp(op),
					p.expression()}}
		}
	}

	return exp
}

func (p *Parser) ternaryExpr() ast.Expr {

	lhs := p.orExpr()

	if p.cur.Kind == ast.HOOK {

		p.consume()
		then := p.expression()
		p.expect(ast.COLON)
		_else := p.ternaryExpr()
		return &ast.TernaryExpr{lhs, then, _else}

	} else {
		return lhs
	}
}

func (p *Parser) orExpr() ast.Expr {

	lhs := p.andExpr()
	for p.cur.Kind == ast.DBL_PIPE {
		tok := p.cur
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.andExpr()}
	}
	return lhs
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
		return p.postfixExpr()
	}
}

func (p *Parser) postfixExpr() ast.Expr {

	exp := p.primaryExpr()

	for isPostfix(p.cur) {

		if asn, ok := exp.(ast.Assignable); ok {
			tok := p.cur
			p.consume()
			exp = &ast.PostfixExpr{asn, tok}
		} else {
			panic(&parserError{INVALID_POSTFIX, p.cur})
		}
	}

	return exp
}

func (p *Parser) primaryExpr() ast.Expr {
	prm := p.primary()

	for {
		// look for suffixes: Invoke, Select, Index, Slice
		switch p.cur.Kind {

		case ast.LPAREN:
			lparen, actual, rparen := p.actualParams()
			prm = &ast.InvokeExpr{prm, lparen, actual, rparen}

		case ast.LBRACKET:
			lbracket := p.consume()

			switch p.cur.Kind {
			case ast.COLON:
				p.consume()
				prm = &ast.SliceToExpr{
					prm,
					lbracket,
					p.expression(),
					p.expect(ast.RBRACKET)}

			default:
				exp := p.expression()

				switch p.cur.Kind {
				case ast.RBRACKET:
					prm = &ast.IndexExpr{
						prm,
						lbracket,
						exp,
						p.expect(ast.RBRACKET)}

				case ast.COLON:
					p.consume()

					switch p.cur.Kind {
					case ast.RBRACKET:
						prm = &ast.SliceFromExpr{
							prm,
							lbracket,
							exp,
							p.expect(ast.RBRACKET)}
					default:
						prm = &ast.SliceExpr{
							prm,
							lbracket,
							exp,
							p.expression(),
							p.expect(ast.RBRACKET)}
					}

				default:
					panic(p.unexpected())
				}
			}

		case ast.DOT:
			p.expect(ast.DOT)
			prm = &ast.FieldExpr{prm, p.expect(ast.IDENT)}

		default:
			return prm
		}
	}
}

func (p *Parser) primary() ast.Expr {

	switch {

	case p.cur.Kind == ast.LPAREN:
		lparen := p.consume()
		expr := p.expression()

		switch p.cur.Kind {
		case ast.RPAREN:
			p.expect(ast.RPAREN)
			return expr

		case ast.COMMA:
			p.expect(ast.COMMA)
			return p.tupleExpr(lparen, expr)

		default:
			panic(p.unexpected())
		}

	case p.cur.Kind == ast.IDENT:
		return p.identExpr()

	case isBuiltIn(p.cur):
		return &ast.BuiltinExpr{p.consume()}

	case p.cur.Kind == ast.THIS:
		return &ast.ThisExpr{p.consume(), nil}

	case p.cur.Kind == ast.FN:
		return p.fnExpr(p.consume())

	case p.cur.Kind == ast.OBJ:
		return p.objExpr(p.consume())

	case p.cur.Kind == ast.DICT:
		return p.dictExpr(p.consume())

	case p.cur.Kind == ast.LBRACKET:
		return p.listExpr(p.consume())

	default:
		return p.basicExpr()
	}
}

func (p *Parser) identExpr() *ast.IdentExpr {
	tok := p.cur
	p.expect(ast.IDENT)
	return &ast.IdentExpr{tok, nil}
}

func (p *Parser) fnExpr(token *ast.Token) ast.Expr {

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
	return &ast.FnExpr{token, params, blk, 0, 0, nil}
}

func (p *Parser) objExpr(objToken *ast.Token) ast.Expr {

	keys := []*ast.Token{}
	values := []ast.Expr{}
	var rbrace *ast.Token

	lbrace := p.expect(ast.LBRACE)

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
				rbrace = p.consume()
				break loop

			default:
				panic(p.unexpected())
			}
		}

	case ast.RBRACE:
		rbrace = p.consume()

	default:
		panic(p.unexpected())
	}

	return &ast.ObjExpr{objToken, lbrace, keys, values, rbrace, -1}
}

func (p *Parser) dictExpr(dictToken *ast.Token) ast.Expr {

	entries := []*ast.DictEntryExpr{}
	var rbrace *ast.Token

	lbrace := p.expect(ast.LBRACE)

	switch p.cur.Kind {

	case ast.RBRACE:
		rbrace = p.consume()

	default:
		key := p.expression()
		p.expect(ast.COLON)
		value := p.expression()
		entries = append(entries, &ast.DictEntryExpr{key, value})

	loop:
		for {
			switch p.cur.Kind {

			case ast.COMMA:
				p.consume()

				key = p.expression()
				p.expect(ast.COLON)
				value = p.expression()
				entries = append(entries, &ast.DictEntryExpr{key, value})

			case ast.RBRACE:
				rbrace = p.consume()
				break loop

			default:
				panic(p.unexpected())
			}
		}
	}

	return &ast.DictExpr{dictToken, lbrace, entries, rbrace}
}

func (p *Parser) listExpr(lbracket *ast.Token) ast.Expr {

	if p.cur.Kind == ast.RBRACKET {
		return &ast.ListExpr{lbracket, []ast.Expr{}, p.consume()}
	} else {

		elems := []ast.Expr{p.expression()}
		for {
			switch p.cur.Kind {
			case ast.RBRACKET:
				return &ast.ListExpr{lbracket, elems, p.consume()}
			case ast.COMMA:
				p.consume()
				elems = append(elems, p.expression())
			default:
				panic(p.unexpected())
			}
		}
	}
}

func (p *Parser) tupleExpr(lparen *ast.Token, expr ast.Expr) ast.Expr {

	elems := []ast.Expr{expr, p.expression()}

	for {
		switch p.cur.Kind {
		case ast.RPAREN:
			return &ast.TupleExpr{lparen, elems, p.consume()}
		case ast.COMMA:
			p.consume()
			elems = append(elems, p.expression())
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) basicExpr() ast.Expr {

	tok := p.cur

	switch {

	case tok.IsBasic():
		p.consume()
		return &ast.BasicExpr{tok}

	default:
		panic(p.unexpected())
	}
}

func (p *Parser) actualParams() (*ast.Token, []ast.Expr, *ast.Token) {

	lparen := p.expect(ast.LPAREN)

	params := []ast.Expr{}
	switch p.cur.Kind {

	case ast.RPAREN:
		return lparen, params, p.consume()

	default:
		params = append(params, p.expression())
		for {
			switch p.cur.Kind {

			case ast.COMMA:
				p.consume()
				params = append(params, p.expression())

			case ast.RPAREN:
				return lparen, params, p.consume()

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
	switch t.Kind {
	case
		ast.DBL_EQ,
		ast.NOT_EQ,
		ast.GT,
		ast.GT_EQ,
		ast.LT,
		ast.LT_EQ,
		ast.CMP,
		ast.HAS:

		return true
	default:
		return false
	}
}

func isAdditive(t *ast.Token) bool {
	switch t.Kind {
	case
		ast.PLUS,
		ast.MINUS,
		ast.PIPE,
		ast.CARET:

		return true
	default:
		return false
	}
}

func isMultiplicative(t *ast.Token) bool {
	switch t.Kind {
	case
		ast.STAR,
		ast.SLASH,
		ast.PERCENT,
		ast.AMP,
		ast.DBL_LT,
		ast.DBL_GT:

		return true
	default:
		return false
	}
}

func isUnary(t *ast.Token) bool {

	switch t.Kind {
	case
		ast.MINUS,
		ast.NOT,
		ast.TILDE:

		return true
	default:
		return false
	}
}

func isPostfix(t *ast.Token) bool {

	switch t.Kind {
	case
		ast.DBL_PLUS,
		ast.DBL_MINUS:

		return true
	default:
		return false
	}
}

func isAssignOp(t *ast.Token) bool {
	switch t.Kind {
	case
		ast.PLUS_EQ,
		ast.MINUS_EQ,
		ast.STAR_EQ,
		ast.SLASH_EQ,
		ast.PERCENT_EQ,
		ast.CARET_EQ,
		ast.AMP_EQ,
		ast.PIPE_EQ,
		ast.DBL_LT_EQ,
		ast.DBL_GT_EQ:

		return true
	default:
		return false
	}
}

func isBuiltIn(t *ast.Token) bool {
	switch t.Kind {
	case
		ast.FN_PRINT,
		ast.FN_PRINTLN,
		ast.FN_STR,
		ast.FN_LEN,
		ast.FN_RANGE:
		return true
	default:
		return false
	}
}

func fromAssignOp(t *ast.Token) *ast.Token {

	switch t.Kind {
	case ast.PLUS_EQ:
		return &ast.Token{ast.PLUS, "+", t.Position}
	case ast.MINUS_EQ:
		return &ast.Token{ast.MINUS, "-", t.Position}
	case ast.STAR_EQ:
		return &ast.Token{ast.STAR, "*", t.Position}
	case ast.SLASH_EQ:
		return &ast.Token{ast.SLASH, "/", t.Position}
	case ast.PERCENT_EQ:
		return &ast.Token{ast.PERCENT, "%", t.Position}
	case ast.CARET_EQ:
		return &ast.Token{ast.CARET, "^", t.Position}
	case ast.AMP_EQ:
		return &ast.Token{ast.AMP, "&", t.Position}
	case ast.PIPE_EQ:
		return &ast.Token{ast.PIPE, "|", t.Position}
	case ast.DBL_LT_EQ:
		return &ast.Token{ast.DBL_LT, "<<", t.Position}
	case ast.DBL_GT_EQ:
		return &ast.Token{ast.DBL_GT, ">>", t.Position}

	default:
		panic("invalid op")
	}
}

//--------------------------------------------------------------
// parserError

type parserErrorKind int

const (
	UNEXPECTED_CHAR parserErrorKind = iota
	UNEXPECTED_TOKEN
	UNEXPECTED_EOF
	INVALID_POSTFIX
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

	case INVALID_POSTFIX:
		return fmt.Sprintf("Invalid Postfix Expression at %v", e.token.Position)

	default:
		panic("unreachable")
	}
}

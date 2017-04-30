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

package compiler

import (
	//"fmt"
	"golem/analyzer"
	"golem/ast"
	g "golem/core"
	"strconv"
)

type Compiler interface {
	ast.Visitor
	Compile() *g.Module
}

type compiler struct {
	anl  analyzer.Analyzer
	pool []g.Value
	opc  []byte
	opln []g.OpcLine

	funcs     []*ast.FnExpr
	templates []*g.Template
	defs      []*g.ObjDef
	idx       int
}

func NewCompiler(anl analyzer.Analyzer) Compiler {

	funcs := []*ast.FnExpr{anl.Module()}
	templates := []*g.Template{}
	defs := []*g.ObjDef{}
	return &compiler{anl, []g.Value{}, nil, nil, funcs, templates, defs, 0}
}

func (c *compiler) Compile() *g.Module {

	for c.idx < len(c.funcs) {
		c.templates = append(
			c.templates,
			c.compileFunc(c.funcs[c.idx]))
		c.idx += 1
	}

	return &g.Module{c.pool, nil, c.defs, c.templates}
}

func (c *compiler) compileFunc(fe *ast.FnExpr) *g.Template {

	arity := len(fe.FormalParams)
	tpl := &g.Template{arity, fe.NumCaptures, fe.NumLocals, nil, nil}

	c.opc = []byte{}
	c.opln = []g.OpcLine{}

	// TODO LOAD_NULL and RETURN are workarounds for the fact that
	// we have not yet written a Control Flow Graph
	c.push(ast.Pos{}, g.LOAD_NULL)
	c.Visit(fe.Body)
	c.push(ast.Pos{}, g.RETURN)

	tpl.OpCodes = c.opc
	tpl.OpcLines = c.opln

	return tpl
}

func (c *compiler) Visit(node ast.Node) {
	switch t := node.(type) {

	case *ast.Const:
		c.visitDecls(t.Decls)

	case *ast.Let:
		c.visitDecls(t.Decls)

	case *ast.Assignment:
		c.visitAssignment(t)

	case *ast.If:
		c.visitIf(t)

	case *ast.While:
		c.visitWhile(t)

	case *ast.For:
		c.visitFor(t)

	case *ast.Break:
		c.visitBreak(t)

	case *ast.Continue:
		c.visitContinue(t)

	case *ast.Return:
		c.visitReturn(t)

	case *ast.TernaryExpr:
		c.visitTernaryExpr(t)

	case *ast.BinaryExpr:
		c.visitBinaryExpr(t)

	case *ast.UnaryExpr:
		c.visitUnaryExpr(t)

	case *ast.PostfixExpr:
		c.visitPostfixExpr(t)

	case *ast.BasicExpr:
		c.visitBasicExpr(t)

	case *ast.IdentExpr:
		c.visitIdentExpr(t)

	case *ast.BuiltinExpr:
		c.visitBuiltinExpr(t)

	case *ast.FnExpr:
		c.visitFunc(t)

	case *ast.InvokeExpr:
		c.visitInvoke(t)

	case *ast.ObjExpr:
		c.visitObjExpr(t)

	case *ast.ThisExpr:
		c.visitThisExpr(t)

	case *ast.FieldExpr:
		c.visitFieldExpr(t)

	case *ast.IndexExpr:
		c.visitIndexExpr(t)

	case *ast.SliceExpr:
		c.visitSliceExpr(t)

	case *ast.SliceFromExpr:
		c.visitSliceFromExpr(t)

	case *ast.SliceToExpr:
		c.visitSliceToExpr(t)

	case *ast.ListExpr:
		c.visitListExpr(t)

	case *ast.TupleExpr:
		c.visitTupleExpr(t)

	case *ast.DictExpr:
		c.visitDictExpr(t)

	default:
		t.Traverse(c)
	}
}

func (c *compiler) visitDecls(decls []*ast.Decl) {

	for _, d := range decls {
		if d.Val == nil {
			c.push(d.Ident.Begin(), g.LOAD_NULL)
		} else {
			c.Visit(d.Val)
		}

		c.assignIdent(d.Ident)
	}
}

func (c *compiler) assignIdent(ident *ast.IdentExpr) {

	v := ident.Variable
	high, low := index(v.Index)
	if v.IsCapture {
		c.push(ident.Begin(), g.STORE_CAPTURE, high, low)
	} else {
		c.push(ident.Begin(), g.STORE_LOCAL, high, low)
	}
}

func (c *compiler) visitAssignment(asn *ast.Assignment) {

	switch t := asn.Assignee.(type) {

	case *ast.IdentExpr:

		c.Visit(asn.Val)

		// TODO doesn't DUP-ing have the potential to fill up the operand stack?
		c.push(asn.Eq.Position, g.DUP)
		c.assignIdent(t)

	case *ast.FieldExpr:

		c.Visit(t.Operand)
		c.Visit(asn.Val)

		high, low := index(len(c.pool))
		c.pool = append(c.pool, g.MakeStr(t.Key.Text))
		c.push(t.Key.Position, g.PUT_FIELD, high, low)

	case *ast.IndexExpr:

		c.Visit(t.Operand)
		c.Visit(t.Index)
		c.Visit(asn.Val)
		c.push(t.Index.Begin(), g.SET_INDEX)

	default:
		panic("invalid assignee type")
	}
}

func (c *compiler) visitPostfixExpr(pe *ast.PostfixExpr) {

	switch t := pe.Assignee.(type) {

	case *ast.IdentExpr:

		c.visitIdentExpr(t)
		c.push(t.Begin(), g.DUP)

		switch pe.Op.Text {
		case "++":
			c.push(pe.Op.Position, g.LOAD_ONE)
		case "--":
			c.push(pe.Op.Position, g.LOAD_NEG_ONE)
		default:
			panic("invalid postfix operator")
		}

		c.push(pe.Op.Position, g.ADD)
		c.assignIdent(t)

	case *ast.FieldExpr:

		c.Visit(t.Operand)

		switch pe.Op.Text {
		case "++":
			c.push(pe.Op.Position, g.LOAD_ONE)
		case "--":
			c.push(pe.Op.Position, g.LOAD_NEG_ONE)
		default:
			panic("invalid postfix operator")
		}

		high, low := index(len(c.pool))
		c.pool = append(c.pool, g.MakeStr(t.Key.Text))
		c.push(t.Key.Position, g.INC_FIELD, high, low)

	case *ast.IndexExpr:

		c.Visit(t.Operand)
		c.Visit(t.Index)

		switch pe.Op.Text {
		case "++":
			c.push(pe.Op.Position, g.LOAD_ONE)
		case "--":
			c.push(pe.Op.Position, g.LOAD_NEG_ONE)
		default:
			panic("invalid postfix operator")
		}

		c.push(t.Index.Begin(), g.INC_INDEX)

	default:
		panic("invalid assignee type")
	}
}

func (c *compiler) visitIf(f *ast.If) {

	c.Visit(f.Cond)

	j0 := c.push(f.Cond.End(), g.JUMP_FALSE, 0xFF, 0xFF)
	f.Then.Traverse(c)

	if f.Else == nil {

		c.setJump(j0, c.opcLen())

	} else {

		j1 := c.push(f.Else.Begin(), g.JUMP, 0xFF, 0xFF)
		c.setJump(j0, c.opcLen())

		f.Else.Traverse(c)
		c.setJump(j1, c.opcLen())
	}
}

func (c *compiler) visitTernaryExpr(f *ast.TernaryExpr) {

	c.Visit(f.Cond)
	j0 := c.push(f.Cond.End(), g.JUMP_FALSE, 0xFF, 0xFF)

	c.Visit(f.Then)
	j1 := c.push(f.Else.Begin(), g.JUMP, 0xFF, 0xFF)
	c.setJump(j0, c.opcLen())

	c.Visit(f.Else)
	c.setJump(j1, c.opcLen())
}

func (c *compiler) visitWhile(w *ast.While) {

	begin := c.opcLen()
	c.Visit(w.Cond)
	j0 := c.push(w.Cond.End(), g.JUMP_FALSE, 0xFF, 0xFF)

	body := c.opcLen()
	w.Body.Traverse(c)
	c.push(w.Body.End(), g.JUMP, begin.high, begin.low)

	end := c.opcLen()
	c.setJump(j0, end)

	c.fixBreakContinue(begin, body, end)
}

func (c *compiler) visitFor(f *ast.For) {

	tok := f.Iterable.Begin()

	c.Visit(f.Iterable)
	c.push(tok, g.ITER)

	begin := c.opcLen()
	c.push(tok, g.DUP)
	c.push(tok, g.ITER_NEXT)
	j0 := c.push(tok, g.JUMP_FALSE, 0xFF, 0xFF)

	c.push(tok, g.DUP)
	c.push(tok, g.ITER_GET)

	if len(f.Idents) == 1 {
		c.assignIdent(f.Idents[0])
	} else {
		panic("not yet implemented")
	}

	body := c.opcLen()
	f.Body.Traverse(c)
	c.push(f.Body.End(), g.JUMP, begin.high, begin.low)

	end := c.opcLen()
	c.setJump(j0, end)

	c.fixBreakContinue(begin, body, end)
}

func (c *compiler) fixBreakContinue(begin *instPtr, body *instPtr, end *instPtr) {

	// replace BREAK and CONTINUE with JUMP
	for i := body.ip; i < end.ip; {
		switch c.opc[i] {
		case g.BREAK:
			c.opc[i] = g.JUMP
			c.opc[i+1] = end.high
			c.opc[i+2] = end.low
		case g.CONTINUE:
			c.opc[i] = g.JUMP
			c.opc[i+1] = begin.high
			c.opc[i+2] = begin.low
		}
		i += g.OpCodeSize(c.opc[i])
	}
}

func (c *compiler) visitBreak(br *ast.Break) {
	c.push(br.Begin(), g.BREAK, 0xFF, 0xFF)
}

func (c *compiler) visitContinue(cn *ast.Continue) {
	c.push(cn.Begin(), g.CONTINUE, 0xFF, 0xFF)
}

func (c *compiler) visitReturn(rt *ast.Return) {
	if rt.Val != nil {
		c.Visit(rt.Val)
	}
	c.push(rt.Begin(), g.RETURN)
}

func (c *compiler) visitBinaryExpr(b *ast.BinaryExpr) {

	switch b.Op.Kind {

	case ast.DBL_PIPE:
		c.visitOr(b.Lhs, b.Rhs)
	case ast.DBL_AMP:
		c.visitAnd(b.Lhs, b.Rhs)

	case ast.DBL_EQ:
		b.Traverse(c)
		c.push(b.Op.Position, g.EQ)
	case ast.NOT_EQ:
		b.Traverse(c)
		c.push(b.Op.Position, g.NE)

	case ast.GT:
		b.Traverse(c)
		c.push(b.Op.Position, g.GT)
	case ast.GT_EQ:
		b.Traverse(c)
		c.push(b.Op.Position, g.GTE)
	case ast.LT:
		b.Traverse(c)
		c.push(b.Op.Position, g.LT)
	case ast.LT_EQ:
		b.Traverse(c)
		c.push(b.Op.Position, g.LTE)
	case ast.CMP:
		b.Traverse(c)
		c.push(b.Op.Position, g.CMP)
	case ast.HAS:
		b.Traverse(c)
		c.push(b.Op.Position, g.HAS)

	case ast.PLUS:
		b.Traverse(c)
		c.push(b.Op.Position, g.ADD)
	case ast.MINUS:
		b.Traverse(c)
		c.push(b.Op.Position, g.SUB)
	case ast.STAR:
		b.Traverse(c)
		c.push(b.Op.Position, g.MUL)
	case ast.SLASH:
		b.Traverse(c)
		c.push(b.Op.Position, g.DIV)

	case ast.PERCENT:
		b.Traverse(c)
		c.push(b.Op.Position, g.REM)
	case ast.AMP:
		b.Traverse(c)
		c.push(b.Op.Position, g.BIT_AND)
	case ast.PIPE:
		b.Traverse(c)
		c.push(b.Op.Position, g.BIT_OR)
	case ast.CARET:
		b.Traverse(c)
		c.push(b.Op.Position, g.BIT_XOR)
	case ast.DBL_LT:
		b.Traverse(c)
		c.push(b.Op.Position, g.LEFT_SHIFT)
	case ast.DBL_GT:
		b.Traverse(c)
		c.push(b.Op.Position, g.RIGHT_SHIFT)

	default:
		panic("unreachable")
	}
}

func (c *compiler) visitOr(lhs ast.Expr, rhs ast.Expr) {

	c.Visit(lhs)
	j0 := c.push(lhs.End(), g.JUMP_TRUE, 0xFF, 0xFF)

	c.Visit(rhs)
	j1 := c.push(rhs.End(), g.JUMP_FALSE, 0xFF, 0xFF)

	c.setJump(j0, c.opcLen())
	c.push(rhs.End(), g.LOAD_TRUE)
	j2 := c.push(rhs.End(), g.JUMP, 0xFF, 0xFF)

	c.setJump(j1, c.opcLen())
	c.push(rhs.End(), g.LOAD_FALSE)

	c.setJump(j2, c.opcLen())
}

func (c *compiler) visitAnd(lhs ast.Expr, rhs ast.Expr) {

	c.Visit(lhs)
	j0 := c.push(lhs.End(), g.JUMP_FALSE, 0xFF, 0xFF)

	c.Visit(rhs)
	j1 := c.push(rhs.End(), g.JUMP_FALSE, 0xFF, 0xFF)

	c.push(rhs.End(), g.LOAD_TRUE)
	j2 := c.push(rhs.End(), g.JUMP, 0xFF, 0xFF)

	c.setJump(j0, c.opcLen())
	c.setJump(j1, c.opcLen())
	c.push(rhs.End(), g.LOAD_FALSE)

	c.setJump(j2, c.opcLen())
}

func (c *compiler) visitUnaryExpr(u *ast.UnaryExpr) {

	switch u.Op.Kind {
	case ast.MINUS:
		opn := u.Operand

		switch t := opn.(type) {
		case *ast.BasicExpr:
			switch t.Token.Kind {

			case ast.INT:
				i := parseInt(t.Token.Text)
				switch i {
				case 0:
					c.push(u.Op.Position, g.LOAD_ZERO)
				case 1:
					c.push(u.Op.Position, g.LOAD_NEG_ONE)
				default:
					high, low := index(len(c.pool))
					c.pool = append(c.pool, g.MakeInt(-i))
					c.push(u.Op.Position, g.LOAD_CONST, high, low)
				}

			default:
				u.Operand.Traverse(c)
				u.Traverse(c)
				c.push(u.Op.Position, g.NEGATE)
			}
		default:
			u.Operand.Traverse(c)
			u.Traverse(c)
			c.push(u.Op.Position, g.NEGATE)
		}

	case ast.NOT:
		u.Traverse(c)
		c.push(u.Op.Position, g.NOT)

	case ast.TILDE:
		u.Traverse(c)
		c.push(u.Op.Position, g.COMPLEMENT)

	default:
		panic("unreachable")
	}
}

func (c *compiler) visitBasicExpr(basic *ast.BasicExpr) {

	high, low := index(len(c.pool))

	// TODO create pool hash map

	switch basic.Token.Kind {

	case ast.NULL:
		c.push(basic.Token.Position, g.LOAD_NULL)

	case ast.TRUE:
		c.push(basic.Token.Position, g.LOAD_TRUE)

	case ast.FALSE:
		c.push(basic.Token.Position, g.LOAD_FALSE)

	case ast.STR:
		c.pool = append(c.pool, g.MakeStr(basic.Token.Text))
		c.push(basic.Token.Position, g.LOAD_CONST, high, low)

	case ast.INT:
		i := parseInt(basic.Token.Text)
		switch i {
		case 0:
			c.push(basic.Token.Position, g.LOAD_ZERO)
		case 1:
			c.push(basic.Token.Position, g.LOAD_ONE)
		default:
			c.pool = append(c.pool, g.MakeInt(i))
			c.push(basic.Token.Position, g.LOAD_CONST, high, low)
		}

	case ast.FLOAT:
		f := parseFloat(basic.Token.Text)
		c.pool = append(c.pool, g.MakeFloat(f))
		c.push(basic.Token.Position, g.LOAD_CONST, high, low)

	default:
		panic("unreachable")
	}

}

func (c *compiler) visitIdentExpr(ident *ast.IdentExpr) {
	v := ident.Variable
	high, low := index(v.Index)
	if v.IsCapture {
		c.push(ident.Begin(), g.LOAD_CAPTURE, high, low)
	} else {
		c.push(ident.Begin(), g.LOAD_LOCAL, high, low)
	}
}

func (c *compiler) visitBuiltinExpr(blt *ast.BuiltinExpr) {

	switch blt.Fn.Kind {
	case ast.FN_PRINT:
		high, low := index(g.PRINT)
		c.push(blt.Fn.Position, g.LOAD_BUILTIN, high, low)
	case ast.FN_PRINTLN:
		high, low := index(g.PRINTLN)
		c.push(blt.Fn.Position, g.LOAD_BUILTIN, high, low)
	case ast.FN_STR:
		high, low := index(g.STR)
		c.push(blt.Fn.Position, g.LOAD_BUILTIN, high, low)
	case ast.FN_LEN:
		high, low := index(g.LEN)
		c.push(blt.Fn.Position, g.LOAD_BUILTIN, high, low)
	case ast.FN_RANGE:
		high, low := index(g.RANGE)
		c.push(blt.Fn.Position, g.LOAD_BUILTIN, high, low)
	case ast.FN_ASSERT:
		high, low := index(g.ASSERT)
		c.push(blt.Fn.Position, g.LOAD_BUILTIN, high, low)

	default:
		panic("unknown builtin function")
	}
}

func (c *compiler) visitFunc(fe *ast.FnExpr) {
	high, low := index(len(c.funcs))
	c.push(fe.Begin(), g.NEW_FUNC, high, low)

	for _, pc := range fe.ParentCaptures {
		high, low = index(pc.Index)
		if pc.IsCapture {
			c.push(fe.Begin(), g.FUNC_CAPTURE, high, low)
		} else {
			c.push(fe.Begin(), g.FUNC_LOCAL, high, low)
		}
	}

	c.funcs = append(c.funcs, fe)
}

func (c *compiler) visitInvoke(inv *ast.InvokeExpr) {

	inv.Traverse(c)
	high, low := index(len(inv.Params))
	c.push(inv.Begin(), g.INVOKE, high, low)
}

func (c *compiler) visitObjExpr(obj *ast.ObjExpr) {

	// create ObjDef for keys
	def := &g.ObjDef{make([]string, len(obj.Keys), len(obj.Keys))}
	for i, k := range obj.Keys {
		def.Keys[i] = k.Text
	}
	high, low := index(len(c.defs))
	c.defs = append(c.defs, def)

	// create un-initialized obj
	c.push(obj.Begin(), g.NEW_OBJ)

	// if the obj is referenced by a 'this', then store local
	if obj.LocalThisIndex != -1 {
		high, low := index(obj.LocalThisIndex)
		c.push(obj.Begin(), g.DUP)
		c.push(obj.Begin(), g.STORE_LOCAL, high, low)
	}

	// eval each value
	for _, v := range obj.Values {
		c.Visit(v)
	}

	// initialize the object
	c.push(obj.End(), g.INIT_OBJ, high, low)
}

func (c *compiler) visitThisExpr(this *ast.ThisExpr) {
	v := this.Variable
	high, low := index(v.Index)
	if v.IsCapture {
		c.push(this.Begin(), g.LOAD_CAPTURE, high, low)
	} else {
		c.push(this.Begin(), g.LOAD_LOCAL, high, low)
	}
}

func (c *compiler) visitFieldExpr(fe *ast.FieldExpr) {
	c.Visit(fe.Operand)
	high, low := index(len(c.pool))
	c.pool = append(c.pool, g.MakeStr(fe.Key.Text))
	c.push(fe.Key.Position, g.GET_FIELD, high, low)
}

func (c *compiler) visitIndexExpr(ie *ast.IndexExpr) {
	c.Visit(ie.Operand)
	c.Visit(ie.Index)
	c.push(ie.Index.Begin(), g.GET_INDEX)
}

func (c *compiler) visitSliceExpr(s *ast.SliceExpr) {
	c.Visit(s.Operand)
	c.Visit(s.From)
	c.Visit(s.To)
	c.push(s.From.Begin(), g.SLICE)
}

func (c *compiler) visitSliceFromExpr(s *ast.SliceFromExpr) {
	c.Visit(s.Operand)
	c.Visit(s.From)
	c.push(s.From.Begin(), g.SLICE_FROM)
}

func (c *compiler) visitSliceToExpr(s *ast.SliceToExpr) {
	c.Visit(s.Operand)
	c.Visit(s.To)
	c.push(s.To.Begin(), g.SLICE_TO)
}

func (c *compiler) visitListExpr(ls *ast.ListExpr) {

	// eval each element
	for _, v := range ls.Elems {
		c.Visit(v)
	}

	// create the list
	high, low := index(len(ls.Elems))
	c.push(ls.Begin(), g.NEW_LIST, high, low)
}

func (c *compiler) visitTupleExpr(tp *ast.TupleExpr) {

	// eval each element
	for _, v := range tp.Elems {
		c.Visit(v)
	}

	// create the list
	high, low := index(len(tp.Elems))
	c.push(tp.Begin(), g.NEW_TUPLE, high, low)
}

func (c *compiler) visitDictExpr(d *ast.DictExpr) {

	// eval each entry
	for _, v := range d.Entries {
		c.Visit(v)
	}

	// create the list
	high, low := index(len(d.Entries))
	c.push(d.Begin(), g.NEW_DICT, high, low)
}

func parseInt(text string) int64 {
	i, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		panic("unreachable")
	}
	if i < 0 {
		panic("unreachable")
	}
	return int64(i)
}

func parseFloat(text string) float64 {
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		panic("unreachable")
	}
	if f < 0 {
		panic("unreachable")
	}
	return float64(f)
}

// returns the length of opc *before* the bytes are pushed
func (c *compiler) push(pos ast.Pos, bytes ...byte) int {
	n := len(c.opc)
	for _, b := range bytes {
		c.opc = append(c.opc, b)
	}

	ln := len(c.opln)
	if (ln == 0) || (pos.Line != c.opln[ln-1].LineNum) {
		c.opln = append(c.opln, g.OpcLine{n, pos.Line})
	}

	return n
}

// replace a mocked-up jump value with the 'real' destination
func (c *compiler) setJump(jmp int, dest *instPtr) {
	c.opc[jmp+1] = dest.high
	c.opc[jmp+2] = dest.low
}

//--------------------------------------------------------------

type instPtr struct {
	ip   int
	high byte
	low  byte
}

func (c *compiler) opcLen() *instPtr {
	high, low := index(len(c.opc))
	return &instPtr{len(c.opc), high, low}
}

func index(n int) (byte, byte) {
	if n >= (2 << 16) {
		panic("TODO wide index")
	}
	return byte((n >> 8) & 0xFF), byte(n & 0xFF)
}

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

package ast

import (
	"bytes"
	"fmt"
	"strings"
)

//--------------------------------------------------------------
// Node

// interfaces
type (
	Node interface {
		fmt.Stringer
		Traverse(Visitor)
		Begin() Pos
		End() Pos
	}

	Stmt interface {
		Node
		stmtMarker()
	}

	Expr interface {
		Node
		exprMarker()
	}

	Assignable interface {
		Expr
		assignableMarker()
	}
)

// structs
type (

	//---------------------
	// statement

	Block struct {
		LBrace *Token
		Nodes  []Node
		RBrace *Token
	}

	Const struct {
		Token     *Token
		Ident     *IdentExpr
		Val       Expr
		Semicolon *Token
	}

	Let struct {
		Token     *Token
		Ident     *IdentExpr
		Val       Expr
		Semicolon *Token
	}

	If struct {
		Token *Token
		Cond  Expr
		Then  *Block
		Else  Stmt
	}

	While struct {
		Token *Token
		Cond  Expr
		Body  *Block
	}

	Break struct {
		Token     *Token
		Semicolon *Token
	}

	Continue struct {
		Token     *Token
		Semicolon *Token
	}

	Return struct {
		Token     *Token
		Val       Expr
		Semicolon *Token
	}

	//---------------------
	// expression

	Assignment struct {
		Assignee Assignable
		Eq       *Token
		Val      Expr
	}

	TernaryExpr struct {
		Cond Expr
		Then Expr
		Else Expr
	}

	BinaryExpr struct {
		Lhs Expr
		Op  *Token
		Rhs Expr
	}

	UnaryExpr struct {
		Op      *Token
		Operand Expr
	}

	PostfixExpr struct {
		Assignee Assignable
		Op       *Token
	}

	BasicExpr struct {
		Token *Token
	}

	IdentExpr struct {
		Symbol   *Token
		Variable *Variable
	}

	FnExpr struct {
		Token *Token

		FormalParams []*IdentExpr
		Body         *Block

		// set by analyzer
		NumLocals      int
		NumCaptures    int
		ParentCaptures []*Variable
	}

	InvokeExpr struct {
		Operand Expr
		LParen  *Token
		Params  []Expr
		RParen  *Token
	}

	ListExpr struct {
		LBracket *Token
		Elems    []Expr
		RBracket *Token
	}

	ObjExpr struct {
		ObjToken *Token
		LBrace   *Token
		Keys     []*Token
		Values   []Expr
		RBrace   *Token

		// The index of the obj expression in the local variable array.
		// '-1' means that the obj is not referenced by a 'this', and thus
		// is not stored in the local variable array
		LocalThisIndex int
	}

	ThisExpr struct {
		Token    *Token
		Variable *Variable
	}

	FieldExpr struct {
		Operand Expr
		Key     *Token
	}

	IndexExpr struct {
		Operand  Expr
		LBracket *Token
		Index    Expr
		RBracket *Token
	}

	SliceExpr struct {
		Operand  Expr
		LBracket *Token
		From     Expr
		To       Expr
		RBracket *Token
	}

	SliceFromExpr struct {
		Operand  Expr
		LBracket *Token
		From     Expr
		RBracket *Token
	}

	SliceToExpr struct {
		Operand  Expr
		LBracket *Token
		To       Expr
		RBracket *Token
	}
)

//--------------------------------------------------------------
// markers

func (*Block) stmtMarker()    {}
func (*Const) stmtMarker()    {}
func (*Let) stmtMarker()      {}
func (*If) stmtMarker()       {}
func (*While) stmtMarker()    {}
func (*Break) stmtMarker()    {}
func (*Continue) stmtMarker() {}
func (*Return) stmtMarker()   {}

func (*Assignment) exprMarker()    {}
func (*TernaryExpr) exprMarker()   {}
func (*BinaryExpr) exprMarker()    {}
func (*UnaryExpr) exprMarker()     {}
func (*PostfixExpr) exprMarker()   {}
func (*BasicExpr) exprMarker()     {}
func (*IdentExpr) exprMarker()     {}
func (*FnExpr) exprMarker()        {}
func (*InvokeExpr) exprMarker()    {}
func (*ListExpr) exprMarker()      {}
func (*ObjExpr) exprMarker()       {}
func (*ThisExpr) exprMarker()      {}
func (*FieldExpr) exprMarker()     {}
func (*IndexExpr) exprMarker()     {}
func (*SliceExpr) exprMarker()     {}
func (*SliceFromExpr) exprMarker() {}
func (*SliceToExpr) exprMarker()   {}

func (*IdentExpr) assignableMarker() {}
func (*FieldExpr) assignableMarker() {}
func (*IndexExpr) assignableMarker() {}

//--------------------------------------------------------------
// Begin, End

func (n *Block) Begin() Pos { return n.LBrace.Position }
func (n *Block) End() Pos   { return n.RBrace.Position }

func (n *Const) Begin() Pos { return n.Token.Position }
func (n *Const) End() Pos   { return n.Semicolon.Position }

func (n *Let) Begin() Pos { return n.Token.Position }
func (n *Let) End() Pos   { return n.Semicolon.Position }

func (n *If) Begin() Pos { return n.Token.Position }
func (n *If) End() Pos {
	if n.Else == nil {
		return n.Then.End()
	} else {
		return n.Else.End()
	}
}

func (n *While) Begin() Pos { return n.Token.Position }
func (n *While) End() Pos   { return n.Body.End() }

func (n *Break) Begin() Pos { return n.Token.Position }
func (n *Break) End() Pos   { return n.Semicolon.Position }

func (n *Continue) Begin() Pos { return n.Token.Position }
func (n *Continue) End() Pos   { return n.Semicolon.Position }

func (n *Return) Begin() Pos { return n.Token.Position }
func (n *Return) End() Pos   { return n.Semicolon.Position }

func (n *Assignment) Begin() Pos { return n.Assignee.Begin() }
func (n *Assignment) End() Pos   { return n.Val.End() }

func (n *TernaryExpr) Begin() Pos { return n.Cond.Begin() }
func (n *TernaryExpr) End() Pos   { return n.Else.End() }

func (n *BinaryExpr) Begin() Pos { return n.Lhs.Begin() }
func (n *BinaryExpr) End() Pos   { return n.Rhs.End() }

func (n *UnaryExpr) Begin() Pos { return n.Op.Position }
func (n *UnaryExpr) End() Pos   { return n.Operand.End() }

func (n *PostfixExpr) Begin() Pos { return n.Assignee.Begin() }
func (n *PostfixExpr) End() Pos   { return n.Op.Position }

func (n *BasicExpr) Begin() Pos { return n.Token.Position }
func (n *BasicExpr) End() Pos {
	return Pos{
		n.Token.Position.Line,
		n.Token.Position.Col + len(n.Token.Text) - 1}
}

func (n *IdentExpr) Begin() Pos { return n.Symbol.Position }
func (n *IdentExpr) End() Pos {
	return Pos{
		n.Symbol.Position.Line,
		n.Symbol.Position.Col + len(n.Symbol.Text) - 1}
}

func (n *FnExpr) Begin() Pos { return n.Token.Position }
func (n *FnExpr) End() Pos   { return n.Body.End() }

func (n *InvokeExpr) Begin() Pos { return n.Operand.Begin() }
func (n *InvokeExpr) End() Pos   { return n.RParen.Position }

func (n *ListExpr) Begin() Pos { return n.LBracket.Position }
func (n *ListExpr) End() Pos   { return n.RBracket.Position }

func (n *ObjExpr) Begin() Pos { return n.ObjToken.Position }
func (n *ObjExpr) End() Pos   { return n.RBrace.Position }

func (n *ThisExpr) Begin() Pos { return n.Token.Position }
func (n *ThisExpr) End() Pos {
	return Pos{
		n.Token.Position.Line,
		n.Token.Position.Col + len("this") - 1}
}

func (n *FieldExpr) Begin() Pos { return n.Operand.Begin() }
func (n *FieldExpr) End() Pos   { return n.Key.Position }

func (n *IndexExpr) Begin() Pos { return n.Operand.Begin() }
func (n *IndexExpr) End() Pos   { return n.RBracket.Position }

func (n *SliceExpr) Begin() Pos     { return n.Operand.Begin() }
func (n *SliceExpr) End() Pos       { return n.RBracket.Position }
func (n *SliceFromExpr) Begin() Pos { return n.Operand.Begin() }
func (n *SliceFromExpr) End() Pos   { return n.RBracket.Position }
func (n *SliceToExpr) Begin() Pos   { return n.Operand.Begin() }
func (n *SliceToExpr) End() Pos     { return n.RBracket.Position }

//--------------------------------------------------------------
// string

func (blk *Block) String() string {
	var buf bytes.Buffer
	buf.WriteString("{ ")
	writeNodes(blk.Nodes, &buf)
	buf.WriteString(" }")
	return buf.String()
}

func (cns *Const) String() string {
	return fmt.Sprintf("const %v = %v;", cns.Ident, cns.Val)
}

func (let *Let) String() string {
	return fmt.Sprintf("let %v = %v;", let.Ident, let.Val)
}

func (asn *Assignment) String() string {
	return fmt.Sprintf("(%v = %v)", asn.Assignee, asn.Val)
}

func (ifn *If) String() string {
	if ifn.Else == nil {
		return fmt.Sprintf("if %v %v", ifn.Cond, ifn.Then)
	} else {
		return fmt.Sprintf("if %v %v else %v", ifn.Cond, ifn.Then, ifn.Else)
	}
}

func (wh *While) String() string {
	return fmt.Sprintf("while %v %v", wh.Cond, wh.Body)
}

func (br *Break) String() string {
	return "break;"
}

func (cn *Continue) String() string {
	return "continue;"
}

func (rt *Return) String() string {
	if rt.Val == nil {
		return "return;"
	} else {
		return fmt.Sprintf("return %v;", rt.Val)
	}
}

func (trn *TernaryExpr) String() string {
	return fmt.Sprintf("(%v ? %v : %v)", trn.Cond, trn.Then, trn.Else)
}

func (bin *BinaryExpr) String() string {
	return fmt.Sprintf("(%v %s %v)", bin.Lhs, bin.Op.Text, bin.Rhs)
}

func (unary *UnaryExpr) String() string {
	return fmt.Sprintf("%s%v", unary.Op.Text, unary.Operand)
}

func (pf *PostfixExpr) String() string {
	return fmt.Sprintf("%v%s", pf.Assignee, pf.Op.Text)
}

func (basic *BasicExpr) String() string {
	if basic.Token.Kind == STR {
		// TODO escape embedded delim, \n, \r, \t, \u
		return strings.Join([]string{"'", basic.Token.Text, "'"}, "")
	} else {
		return basic.Token.Text
	}
}

func (ident *IdentExpr) String() string {
	return ident.Symbol.Text
}

func (fn *FnExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString("fn(")
	for idx, p := range fn.FormalParams {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(p.String())
	}
	buf.WriteString(") ")

	buf.WriteString(fn.Body.String())

	return buf.String()
}

func (inv *InvokeExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(inv.Operand.String())
	buf.WriteString("(")
	for idx, p := range inv.Params {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(p.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (ls *ListExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString("[ ")
	for idx, v := range ls.Elems {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.String())
	}
	buf.WriteString(" ]")
	return buf.String()
}

func (obj *ObjExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString("obj { ")
	for idx, k := range obj.Keys {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(k.Text)
		buf.WriteString(": ")
		buf.WriteString(obj.Values[idx].String())
	}
	buf.WriteString(" }")
	return buf.String()
}

func (this *ThisExpr) String() string {
	return "this"
}

func (f *FieldExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(f.Operand.String())
	buf.WriteString(".")
	buf.WriteString(f.Key.Text)
	return buf.String()
}

func (i *IndexExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(i.Operand.String())
	buf.WriteString("[")
	buf.WriteString(i.Index.String())
	buf.WriteString("]")
	return buf.String()
}

func (s *SliceExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(s.Operand.String())
	buf.WriteString("[")
	buf.WriteString(s.From.String())
	buf.WriteString(":")
	buf.WriteString(s.To.String())
	buf.WriteString("]")
	return buf.String()
}

func (s *SliceFromExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(s.Operand.String())
	buf.WriteString("[")
	buf.WriteString(s.From.String())
	buf.WriteString(":]")
	return buf.String()
}

func (s *SliceToExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(s.Operand.String())
	buf.WriteString("[:")
	buf.WriteString(s.To.String())
	buf.WriteString("]")
	return buf.String()
}

func writeNodes(nodes []Node, buf *bytes.Buffer) {
	for idx, n := range nodes {
		if idx > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(n.String())
		if _, ok := n.(Expr); ok {
			buf.WriteString(";")
		}
	}
}

//--------------------------------------------------------------
// A Variable points to a Ref.  Variables are defined either
// as formal params for a Function, or via Let or Const, or via
// the capture mechanism.

type Variable struct {
	Index     int
	IsConst   bool
	IsCapture bool
}

func (v *Variable) String() string {
	return fmt.Sprintf("(%d,%v,%v)", v.Index, v.IsConst, v.IsCapture)
}

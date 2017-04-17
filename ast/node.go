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
		//Begin() Pos
		//End() Pos
	}

	Stmt interface {
		Node
		stmtMarker()
	}

	Expr interface {
		Node
		exprMarker()
	}
)

// structs
type (

	//---------------------
	// statement

	Block struct {
		Nodes []Node
	}

	Const struct {
		Ident *IdentExpr
		Val   Expr
	}

	Let struct {
		Ident *IdentExpr
		Val   Expr
	}

	If struct {
		Cond Expr
		Then *Block
		Else Stmt
	}

	While struct {
		Cond Expr
		Body *Block
	}

	Break struct {
	}

	Continue struct {
	}

	Return struct {
		Val Expr
	}

	//---------------------
	// expression

	Assignment struct {
		Ident *IdentExpr
		Val   Expr
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

	BasicExpr struct {
		Token *Token
	}

	IdentExpr struct {
		Symbol   *Token
		Variable *Variable
	}

	FnExpr struct {
		FormalParams []*IdentExpr
		Body         *Block

		// set by analyzer
		NumLocals      int
		NumCaptures    int
		ParentCaptures []*Variable
	}

	InvokeExpr struct {
		Operand Expr
		Params  []Expr
	}

	ObjExpr struct {
		Keys   []*Token
		Values []Expr
		// The index of the obj expression in the local variable array.
		// '-1' means that the obj is not referenced by a 'this', and thus
		// is not stored in the local variable array
		LocalThisIndex int
	}

	ThisExpr struct {
		Variable *Variable
	}

	SelectExpr struct {
		Operand Expr
		Key     *Token
	}

	PutExpr struct {
		Operand Expr
		Key     *Token
		Value   Expr
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

func (*Assignment) exprMarker() {}
func (*BinaryExpr) exprMarker() {}
func (*UnaryExpr) exprMarker()  {}
func (*BasicExpr) exprMarker()  {}
func (*IdentExpr) exprMarker()  {}
func (*FnExpr) exprMarker()     {}
func (*InvokeExpr) exprMarker() {}
func (*ObjExpr) exprMarker()    {}
func (*ThisExpr) exprMarker()   {}
func (*SelectExpr) exprMarker() {}
func (*PutExpr) exprMarker()    {}

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
	return fmt.Sprintf("%v = %v", asn.Ident, asn.Val)
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

func (bin *BinaryExpr) String() string {
	return fmt.Sprintf("(%v %s %v)", bin.Lhs, bin.Op.Text, bin.Rhs)
}

func (unary *UnaryExpr) String() string {
	return fmt.Sprintf("%s%v", unary.Op.Text, unary.Operand)
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

func (s *SelectExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(s.Operand.String())
	buf.WriteString(".")
	buf.WriteString(s.Key.Text)
	return buf.String()
}

func (p *PutExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(p.Operand.String())
	buf.WriteString(".")
	buf.WriteString(p.Key.Text)
	buf.WriteString(" = ")
	buf.WriteString(p.Value.String())
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

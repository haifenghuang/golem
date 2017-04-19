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
)

//--------------------------------------------------------------
// Visitor

type Visitor interface {
	Visit(node Node)
}

func (cns *Const) Traverse(v Visitor) {
	// Do not traverse cns.Ident!!!
	// It will confuse the Analyzer.
	v.Visit(cns.Val)
}

func (let *Let) Traverse(v Visitor) {
	// Do not traverse let.Ident!!!
	// It will confuse the Analyzer.
	v.Visit(let.Val)
}

func (asn *Assignment) Traverse(v Visitor) {
	// Do not traverse asn.Ident!!!
	// It will confuse the Analyzer.
	v.Visit(asn.Val)
}

func (ifn *If) Traverse(v Visitor) {
	v.Visit(ifn.Cond)
	v.Visit(ifn.Then)
	if ifn.Else != nil {
		v.Visit(ifn.Else)
	}
}

func (wh *While) Traverse(v Visitor) {
	v.Visit(wh.Cond)
	v.Visit(wh.Body)
}

func (br *Break) Traverse(v Visitor) {
}

func (cn *Continue) Traverse(v Visitor) {
}

func (rt *Return) Traverse(v Visitor) {
	if rt.Val != nil {
		v.Visit(rt.Val)
	}
}

func (blk *Block) Traverse(v Visitor) {
	for _, n := range blk.Nodes {
		v.Visit(n)
	}
}

func (bin *BinaryExpr) Traverse(v Visitor) {
	v.Visit(bin.Lhs)
	v.Visit(bin.Rhs)
}

func (un *UnaryExpr) Traverse(v Visitor) {
	v.Visit(un.Operand)
}

func (pf *PostfixExpr) Traverse(v Visitor) {
	v.Visit(pf.Operand)
}

func (basic *BasicExpr) Traverse(v Visitor) {
}

func (ident *IdentExpr) Traverse(v Visitor) {
}

func (fn *FnExpr) Traverse(v Visitor) {
	for _, n := range fn.FormalParams {
		v.Visit(n)
	}
	v.Visit(fn.Body)
}

func (inv *InvokeExpr) Traverse(v Visitor) {
	v.Visit(inv.Operand)
	for _, n := range inv.Params {
		v.Visit(n)
	}
}

func (obj *ObjExpr) Traverse(v Visitor) {
	for _, val := range obj.Values {
		v.Visit(val)
	}
}

func (this *ThisExpr) Traverse(v Visitor) {
}

func (f *FieldExpr) Traverse(v Visitor) {
	v.Visit(f.Operand)
}

func (p *PutExpr) Traverse(v Visitor) {
	v.Visit(p.Operand)
	v.Visit(p.Value)
}

//--------------------------------------------------------------
// ast debug

type dump struct {
	buf    bytes.Buffer
	indent int
}

func Dump(node Node) string {
	p := &dump{}
	p.Visit(node)
	return p.buf.String()
}

func (p *dump) Visit(node Node) {

	for i := 0; i < p.indent; i++ {
		p.buf.WriteString(".   ")
	}

	switch t := node.(type) {

	case *Block:
		p.buf.WriteString("Block\n")

	case *Const:
		p.buf.WriteString("Const\n")
		p.indent++
		p.Visit(t.Ident)
		p.indent--
	case *Let:
		p.buf.WriteString("Let\n")
		p.indent++
		p.Visit(t.Ident)
		p.indent--
	case *Assignment:
		p.buf.WriteString("Assignment\n")
		p.indent++
		p.Visit(t.Ident)
		p.indent--

	case *If:
		p.buf.WriteString("If\n")
	case *While:
		p.buf.WriteString("While\n")
	case *Break:
		p.buf.WriteString("Break\n")
	case *Continue:
		p.buf.WriteString("Continue\n")
	case *Return:
		p.buf.WriteString("Return\n")

	case *BinaryExpr:
		p.buf.WriteString(fmt.Sprintf("BinaryExpr(%q)\n", t.Op.Text))
	case *UnaryExpr:
		p.buf.WriteString(fmt.Sprintf("UnaryExpr(%q)\n", t.Op.Text))
	case *PostfixExpr:
		p.buf.WriteString(fmt.Sprintf("PostfixExpr(%q)\n", t.Op.Text))
	case *BasicExpr:
		p.buf.WriteString(fmt.Sprintf("BasicExpr(%v,%q)\n", t.Token.Kind, t.Token.Text))
	case *IdentExpr:
		p.buf.WriteString(fmt.Sprintf("IdentExpr(%v,%v)\n", t.Symbol.Text, t.Variable))

	case *FnExpr:
		p.buf.WriteString(fmt.Sprintf("FnExpr(numLocals:%d", t.NumLocals))
		p.buf.WriteString(fmt.Sprintf(" numCaptures:%d", t.NumCaptures))
		p.buf.WriteString(" parentCaptures:")
		p.buf.WriteString(varsString(t.ParentCaptures))
		p.buf.WriteString(")\n")
	case *InvokeExpr:
		p.buf.WriteString("InvokeExpr\n")

	case *ObjExpr:
		p.buf.WriteString(fmt.Sprintf("ObjExpr(%v,%d)\n", tokensString(t.Keys), t.LocalThisIndex))
	case *ThisExpr:
		p.buf.WriteString(fmt.Sprintf("ThisExpr(%v)\n", t.Variable))

	case *FieldExpr:
		p.buf.WriteString(fmt.Sprintf("FieldExpr(%v)\n", t.Key.Text))
	case *PutExpr:
		p.buf.WriteString(fmt.Sprintf("PutExpr(%v)\n", t.Key.Text))

	default:
		panic("cannot visit")
	}

	p.indent++
	node.Traverse(p)
	p.indent--
}

func varsString(vars []*Variable) string {

	var buf bytes.Buffer
	buf.WriteString("[")
	n := 0
	for v := range vars {
		if n > 0 {
			buf.WriteString(", ")
		}
		n++
		buf.WriteString(fmt.Sprintf("%v", vars[v]))
	}
	buf.WriteString("]")
	return buf.String()
}

func tokensString(tokens []*Token) string {

	var buf bytes.Buffer
	buf.WriteString("[")
	n := 0
	for t := range tokens {
		if n > 0 {
			buf.WriteString(", ")
		}
		n++
		buf.WriteString(fmt.Sprintf("%v", tokens[t].Text))
	}
	buf.WriteString("]")
	return buf.String()
}

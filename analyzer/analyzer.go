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

package analyzer

import (
	"fmt"
	"golem/ast"
	"sort"
)

type Analyzer interface {
	ast.Visitor
	Module() *ast.FnExpr
	Analyze() []error
	scope() *scope
}

type analyzer struct {
	mod       *ast.FnExpr
	rootScope *scope
	curScope  *scope
	loops     []*ast.While
	objs      []*ast.ObjExpr
	errors    []error
}

func NewAnalyzer(mod *ast.FnExpr) Analyzer {

	rootScope := newFuncScope(nil)

	return &analyzer{mod, rootScope, rootScope, []*ast.While{}, []*ast.ObjExpr{}, nil}
}

func (a *analyzer) scope() *scope {
	return a.rootScope
}

func (a *analyzer) Analyze() []error {

	a.doVisitFunc(a.mod)

	return a.errors
}

func (a *analyzer) Module() *ast.FnExpr {
	return a.mod
}

func (a *analyzer) Visit(node ast.Node) {
	switch t := node.(type) {

	case *ast.Block:
		a.visitBlock(t)

	case *ast.FnExpr:
		a.visitFunc(t)

	case *ast.Const:
		a.visitDecls(t.Decls, true)

	case *ast.Let:
		a.visitDecls(t.Decls, false)

	case *ast.Assignment:
		a.visitAssignment(t)

	case *ast.PostfixExpr:
		a.visitPostfixExpr(t)

	case *ast.IdentExpr:
		a.visitIdentExpr(t)

	case *ast.While:
		a.loops = append(a.loops, t)
		t.Traverse(a)
		a.loops = a.loops[:len(a.loops)-1]

	case *ast.Break:
		if len(a.loops) == 0 {
			a.errors = append(a.errors, &aerror{"'break' outside of loop"})
		}

	case *ast.Continue:
		if len(a.loops) == 0 {
			a.errors = append(a.errors, &aerror{"'continue' outside of loop"})
		}

	case *ast.ObjExpr:
		a.visitObjExpr(t)

	case *ast.ThisExpr:
		a.visitThisExpr(t)

	default:
		t.Traverse(a)

	}
}

func (a *analyzer) visitBlock(blk *ast.Block) {

	a.curScope = newBlockScope(a.curScope)
	blk.Traverse(a)
	a.curScope = a.curScope.parent
}

func (a *analyzer) visitFunc(fn *ast.FnExpr) {

	a.curScope = newFuncScope(a.curScope)
	a.doVisitFunc(fn)
	a.curScope = a.curScope.parent
}

func (a *analyzer) doVisitFunc(fn *ast.FnExpr) {

	for _, f := range fn.FormalParams {
		f.Variable = a.curScope.put(f.Symbol.Text, false)
	}
	a.visitBlock(fn.Body)

	af := a.curScope.funcScope

	m := af.parentCaptures

	// Sort the keys so that the list comes out the same every time.
	keys := make([]string, len(m))
	i := 0
	for k, _ := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	// make an array out of the values
	pc := make([]*ast.Variable, 0, len(m))
	for _, k := range keys {
		pc = append(pc, m[k])
	}

	fn.NumLocals = af.numLocals
	fn.NumCaptures = len(af.captures)
	fn.ParentCaptures = pc
}

func (a *analyzer) visitDecls(decls []*ast.Decl, isConst bool) {

	for _, d := range decls {
		if d.Val != nil {
			a.Visit(d.Val)
		}

		sym := d.Ident.Symbol.Text
		if _, ok := a.curScope.get(sym); ok {
			a.errors = append(a.errors,
				&aerror{fmt.Sprintf("Symbol '%s' is already defined", sym)})
		} else {
			d.Ident.Variable = a.curScope.put(sym, isConst)
		}
	}
}

func (a *analyzer) visitAssignment(asn *ast.Assignment) {

	switch t := asn.Assignee.(type) {

	case *ast.IdentExpr:
		a.Visit(asn.Val)
		a.doVisitAssignIdent(t)

	case *ast.FieldExpr:
		a.Visit(t.Operand)
		a.Visit(asn.Val)

	case *ast.IndexExpr:
		a.Visit(t.Operand)
		a.Visit(t.Index)
		a.Visit(asn.Val)

	default:
		panic("invalid assignee type")
	}
}

func (a *analyzer) visitPostfixExpr(ps *ast.PostfixExpr) {

	switch t := ps.Assignee.(type) {

	case *ast.IdentExpr:
		a.doVisitAssignIdent(t)

	case *ast.FieldExpr:
		a.Visit(t.Operand)

	case *ast.IndexExpr:
		a.Visit(t.Operand)
		a.Visit(t.Index)

	default:
		panic("invalid assignee type")
	}
}

// visit an Ident that is part of an assignment
func (a *analyzer) doVisitAssignIdent(ident *ast.IdentExpr) {
	sym := ident.Symbol.Text
	if v, ok := a.curScope.get(sym); ok {
		if v.IsConst {
			a.errors = append(a.errors,
				&aerror{fmt.Sprintf("Symbol '%s' is constant", sym)})
		}
		ident.Variable = v
	} else {
		a.errors = append(a.errors,
			&aerror{fmt.Sprintf("Symbol '%s' is not defined", sym)})
	}
}

func (a *analyzer) visitIdentExpr(ident *ast.IdentExpr) {

	sym := ident.Symbol.Text

	if v, ok := a.curScope.get(sym); ok {
		ident.Variable = v
	} else {
		a.errors = append(a.errors,
			&aerror{fmt.Sprintf("Symbol '%s' is not defined", sym)})
	}
}

func (a *analyzer) visitObjExpr(obj *ast.ObjExpr) {
	a.objs = append(a.objs, obj)

	a.curScope = newObjScope(a.curScope, obj)
	obj.Traverse(a)
	a.curScope = a.curScope.parent

	a.objs = a.objs[:len(a.objs)-1]
}

func (a *analyzer) visitThisExpr(this *ast.ThisExpr) {

	n := len(a.objs)
	if n == 0 {
		a.errors = append(a.errors, &aerror{"'this' outside of loop"})
	} else {
		this.Variable = a.curScope.this()
	}
}

//--------------------------------------------------------------
// aerror

type aerror struct {
	msg string
}

func (e *aerror) Error() string {
	return e.msg
}

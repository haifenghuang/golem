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
	//"fmt"
	"golem/ast"
	"golem/scanner"
	"testing"
)

func ok(t *testing.T, p *Parser, expect string) {

	mod, err := p.ParseModule()
	if err != nil {
		t.Error(err, " != nil")
	}

	if mod.String() != expect {
		t.Error(mod, " != ", expect)
	}
}

func fail(t *testing.T, p *Parser, expect string) {

	mod, err := p.ParseModule()
	if mod != nil {
		t.Error(mod, " != nil")
	}

	if err.Error() != expect {
		t.Error(err, " != ", expect)
	}
}

func ok_expr(t *testing.T, p *Parser, expect string) {

	expr, err := p.parseExpression()
	if err != nil {
		t.Error(err, " != nil")
	}

	if expr.String() != expect {
		t.Error(expr, " != ", expect)
	}
}

func fail_expr(t *testing.T, p *Parser, expect string) {

	expr, err := p.parseExpression()
	if expr != nil {
		t.Error(expr, " != nil")
	}

	if err.Error() != expect {
		t.Error(err, " != ", expect)
	}
}

func newParser(source string) *Parser {
	return NewParser(scanner.NewScanner(source))
}

func TestPrimary(t *testing.T) {

	p := newParser("")
	fail_expr(t, p, "Unexpected EOF at (1, 1)")

	p = newParser("#")
	fail_expr(t, p, "Unexpected Character '#' at (1, 1)")

	p = newParser("'")
	fail_expr(t, p, "Unexpected EOF at (1, 2)")

	p = newParser("1 2")
	fail_expr(t, p, "Unexpected Token '2' at (1, 3)")

	p = newParser("1 #")
	fail_expr(t, p, "Unexpected Character '#' at (1, 3)")

	p = newParser("1")
	ok_expr(t, p, "1")

	p = newParser("0xa")
	ok_expr(t, p, "0xa")

	p = newParser("1.2")
	ok_expr(t, p, "1.2")

	p = newParser("null")
	ok_expr(t, p, "null")

	p = newParser("true")
	ok_expr(t, p, "true")

	p = newParser("false")
	ok_expr(t, p, "false")

	p = newParser("'a'")
	ok_expr(t, p, "'a'")

	p = newParser("('a')")
	ok_expr(t, p, "'a'")

	p = newParser("bar")
	ok_expr(t, p, "bar")
}

func TestUnary(t *testing.T) {
	p := newParser("-1")
	ok_expr(t, p, "-1")

	p = newParser("- - 2")
	ok_expr(t, p, "--2")

	p = newParser("--2")
	ok_expr(t, p, "--2")

	p = newParser("!a")
	ok_expr(t, p, "!a")
}

func TestMultiplicative(t *testing.T) {
	p := newParser("1*2")
	ok_expr(t, p, "(1 * 2)")

	p = newParser("-1*-2")
	ok_expr(t, p, "(-1 * -2)")

	p = newParser("1*2*3")
	ok_expr(t, p, "((1 * 2) * 3)")

	p = newParser("1*2/3*4/5")
	ok_expr(t, p, "((((1 * 2) / 3) * 4) / 5)")
}

func TestAdditive(t *testing.T) {
	p := newParser("1*2+3")
	ok_expr(t, p, "((1 * 2) + 3)")

	p = newParser("1+2*3")
	ok_expr(t, p, "(1 + (2 * 3))")

	p = newParser("1+2*-3")
	ok_expr(t, p, "(1 + (2 * -3))")

	p = newParser("1+2+-3")
	ok_expr(t, p, "((1 + 2) + -3)")

	p = newParser("1+2*3+4")
	ok_expr(t, p, "((1 + (2 * 3)) + 4)")

	p = newParser("(1+2) * 3")
	ok_expr(t, p, "((1 + 2) * 3)")

	p = newParser("(1*2) * 3")
	ok_expr(t, p, "((1 * 2) * 3)")

	p = newParser("1 * (2 + 3)")
	ok_expr(t, p, "(1 * (2 + 3))")

	p = newParser("1 +")
	fail_expr(t, p, "Unexpected EOF at (1, 4)")
}

func TestComparitive(t *testing.T) {
	p := newParser("1==3")
	ok_expr(t, p, "(1 == 3)")

	p = newParser("1 ==2 +3 * - 4")
	ok_expr(t, p, "(1 == (2 + (3 * -4)))")

	p = newParser("(1== 2)+ 3")
	ok_expr(t, p, "((1 == 2) + 3)")

	p = newParser("1!=3")
	ok_expr(t, p, "(1 != 3)")

	ok_expr(t, newParser("1 < 3"), "(1 < 3)")
	ok_expr(t, newParser("1 > 3"), "(1 > 3)")
	ok_expr(t, newParser("1 <= 3"), "(1 <= 3)")
	ok_expr(t, newParser("1 >= 3"), "(1 >= 3)")
	ok_expr(t, newParser("1 <=> 3"), "(1 <=> 3)")

	ok_expr(t, newParser("1 <=> 2 + 3 * 4"), "(1 <=> (2 + (3 * 4)))")
}

func TestAndOr(t *testing.T) {

	ok_expr(t, newParser("1 || 2"), "(1 || 2)")
	ok_expr(t, newParser("1 || 2 || 3"), "((1 || 2) || 3)")

	ok_expr(t, newParser("1 || 2 && 3"), "(1 || (2 && 3))")
	ok_expr(t, newParser("1 || 2 && 3 < 4"), "(1 || (2 && (3 < 4)))")
}

func TestModule(t *testing.T) {
	p := newParser("let a =1==3; 2+ true; z =27;const a = 3;")
	ok(t, p, "fn() { let a = (1 == 3); (2 + true); z = 27; const a = 3; }")
}

func TestStatement(t *testing.T) {
	p := newParser("if a { b;let c=12; }")
	ok(t, p, "fn() { if a { b; let c = 12; } }")

	p = newParser("if a { b; } else { c; }")
	ok(t, p, "fn() { if a { b; } else { c; } }")

	p = newParser("if a { b; } else { if(12 == 3) { z+5; }}")
	ok(t, p, "fn() { if a { b; } else { if (12 == 3) { (z + 5); } } }")

	p = newParser("if a {} else if b {} else {}")
	ok(t, p, "fn() { if a {  } else if b {  } else {  } }")

	p = newParser("while a { b; }")
	ok(t, p, "fn() { while a { b; } }")

	p = newParser("break; continue; while a { b; continue; break; }")
	ok(t, p, "fn() { break; continue; while a { b; continue; break; } }")

	p = newParser("a = b;")
	ok(t, p, "fn() { a = b; }")
}

func TestFn(t *testing.T) {
	p := newParser("fn() { }")
	ok_expr(t, p, "fn() {  }")

	p = newParser("fn() { a = 3; }")
	ok_expr(t, p, "fn() { a = 3; }")

	p = newParser("fn(x) { a = 3; }")
	ok_expr(t, p, "fn(x) { a = 3; }")

	p = newParser("fn(x,y) { a = 3; }")
	ok_expr(t, p, "fn(x, y) { a = 3; }")

	p = newParser("fn(x,y,z) { a = 3; }")
	ok_expr(t, p, "fn(x, y, z) { a = 3; }")

	p = newParser("fn(x) { let a = fn(y) { return x + y; }; }")
	ok_expr(t, p, "fn(x) { let a = fn(y) { return (x + y); }; }")

	p = newParser("return;")
	ok(t, p, "fn() { return; }")

	p = newParser("z = fn(x) { a = 2; return b; c = 3; };")
	ok(t, p, "fn() { z = fn(x) { a = 2; return b; c = 3; }; }")
}

func TestInvoke(t *testing.T) {
	p := newParser("a()")
	ok_expr(t, p, "a()")

	p = newParser("a(1)")
	ok_expr(t, p, "a(1)")

	p = newParser("a(1, 2, 3)")
	ok_expr(t, p, "a(1, 2, 3)")
}

func TestObj(t *testing.T) {
	p := newParser("obj{}")
	ok_expr(t, p, "obj {  }")

	p = newParser("obj{a:1}")
	ok_expr(t, p, "obj { a: 1 }")

	p = newParser("obj{a:1,b:2}")
	ok_expr(t, p, "obj { a: 1, b: 2 }")

	p = newParser("obj{a:1,b:2,c:3}")
	ok_expr(t, p, "obj { a: 1, b: 2, c: 3 }")

	p = newParser("obj{a:1,b:2,c:obj{d:3}}")
	ok_expr(t, p, "obj { a: 1, b: 2, c: obj { d: 3 } }")

	p = newParser("obj{a:1, b: fn(x) { y + x;} }")
	ok_expr(t, p, "obj { a: 1, b: fn(x) { (y + x); } }")

	p = newParser("obj{a:1, b: fn(x) { y + x;}, c: obj {d:3} }")
	ok_expr(t, p, "obj { a: 1, b: fn(x) { (y + x); }, c: obj { d: 3 } }")

	p = newParser("a.b")
	ok_expr(t, p, "a.b")

	p = newParser("a.b = 3")
	ok_expr(t, p, "a.b = 3")

	p = newParser("let a.b = 3;")
	fail(t, p, "Unexpected Token '.' at (1, 6)")

	p = newParser("this")
	ok_expr(t, p, "this")

	p = newParser("obj{a:this + true, b: this}")
	ok_expr(t, p, "obj { a: (this + true), b: this }")

	p = newParser("a = this")
	ok_expr(t, p, "a = this")

	p = newParser("obj{ a: this }")
	ok_expr(t, p, "obj { a: this }")

	p = newParser("obj{ a: this == 2 }")
	ok_expr(t, p, "obj { a: (this == 2) }")

	p = newParser("this.b = 3")
	ok_expr(t, p, "this.b = 3")

	p = newParser("obj { a: this.b = 3 }")
	ok_expr(t, p, "obj { a: this.b = 3 }")

	p = newParser("b = this")
	ok_expr(t, p, "b = this")

	p = newParser("obj { a: b = this }")
	ok_expr(t, p, "obj { a: b = this }")

	p = newParser("this = b")
	fail(t, p, "Unexpected Token '=' at (1, 6)")
}

func TestPrimarySuffix(t *testing.T) {
	p := newParser("a.b()")
	ok_expr(t, p, "a.b()")

	p = newParser("a.b.c")
	ok_expr(t, p, "a.b.c")

	p = newParser("a.b().c")
	ok_expr(t, p, "a.b().c")
}

func okExprPos(t *testing.T, p *Parser, expectBegin ast.Pos, expectEnd ast.Pos) {

	expr, err := p.parseExpression()
	if err != nil {
		t.Error(err, " != nil")
	}

	if expr.Begin() != expectBegin {
		t.Error(expr.Begin(), " != ", expectBegin)
	}

	if expr.End() != expectEnd {
		t.Error(expr.End(), " != ", expectEnd)
	}
}

func okPos(t *testing.T, p *Parser, expectBegin ast.Pos, expectEnd ast.Pos) {

	mod, err := p.ParseModule()
	if err != nil {
		t.Error(err, " != nil")
	}

	if len(mod.Body.Nodes) != 1 {
		t.Error("node count", len(mod.Body.Nodes))
	}

	n := mod.Body.Nodes[0]
	if n.Begin() != expectBegin {
		t.Error(n.Begin(), " != ", expectBegin)
	}

	if n.End() != expectEnd {
		t.Error(n.End(), " != ", expectEnd)
	}
}

func TestPos(t *testing.T) {
	p := newParser("1.23")
	okExprPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 4})

	p = newParser("-1")
	okExprPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 2})

	p = newParser("null + true")
	okExprPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 11})

	p = newParser("a1")
	okExprPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 2})

	p = newParser("a = \n3")
	okExprPos(t, p, ast.Pos{1, 1}, ast.Pos{2, 1})

	p = newParser("a(b,c)")
	okExprPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 6})

	p = newParser("obj{}")
	okExprPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 5})

	p = newParser("obj { a: 1 }")
	okExprPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 12})

	p = newParser("   this")
	okExprPos(t, p, ast.Pos{1, 4}, ast.Pos{1, 7})

	p = newParser("a.b")
	okExprPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 3})

	p = newParser("a.b = 2")
	okExprPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 7})

	p = newParser(`
fn() { 
    return x; 
}`)
	okExprPos(t, p, ast.Pos{2, 1}, ast.Pos{4, 1})

	p = newParser("const a = 1;")
	okPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 12})

	p = newParser("let a = 1\n;")
	okPos(t, p, ast.Pos{1, 1}, ast.Pos{2, 1})

	p = newParser("break;")
	okPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 6})

	p = newParser("\n  continue;")
	okPos(t, p, ast.Pos{2, 3}, ast.Pos{2, 11})

	p = newParser("return;")
	okPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 7})

	p = newParser("while true { 42; \n}")
	okPos(t, p, ast.Pos{1, 1}, ast.Pos{2, 1})

	p = newParser("if 0 {}")
	okPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 7})

	p = newParser("if 0 {} else {}")
	okPos(t, p, ast.Pos{1, 1}, ast.Pos{1, 15})
}

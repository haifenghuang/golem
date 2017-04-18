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

package scanner

import (
	//"fmt"
	"golem/ast"
	"reflect"
	"testing"
)

func ok(t *testing.T, s *Scanner, tokenKind ast.TokenKind, text string, line int, col int) {

	token := &ast.Token{tokenKind, text, ast.Pos{line, col}}

	nextToken := s.Next()

	if !reflect.DeepEqual(*nextToken, *token) {
		t.Error(nextToken, " != ", token)
	}
}

func TestDelimiter(t *testing.T) {

	s := NewScanner("")
	ok(t, s, ast.EOF, "", 1, 1)
	ok(t, s, ast.EOF, "", 1, 1)
	ok(t, s, ast.EOF, "", 1, 1)
	ok(t, s, ast.EOF, "", 1, 1)

	s = NewScanner("#")
	ok(t, s, ast.UNEXPECTED_CHAR, "#", 1, 1)
	ok(t, s, ast.UNEXPECTED_CHAR, "#", 1, 1)
	ok(t, s, ast.UNEXPECTED_CHAR, "#", 1, 1)

	s = NewScanner("+")
	ok(t, s, ast.PLUS, "+", 1, 1)
	ok(t, s, ast.EOF, "", 1, 2)

	s = NewScanner("-\n/")
	ok(t, s, ast.MINUS, "-", 1, 1)
	ok(t, s, ast.DIV, "/", 2, 1)
	ok(t, s, ast.EOF, "", 2, 2)

	s = NewScanner("+-*/)(")
	ok(t, s, ast.PLUS, "+", 1, 1)
	ok(t, s, ast.MINUS, "-", 1, 2)
	ok(t, s, ast.MULT, "*", 1, 3)
	ok(t, s, ast.DIV, "/", 1, 4)
	ok(t, s, ast.RPAREN, ")", 1, 5)
	ok(t, s, ast.LPAREN, "(", 1, 6)
	ok(t, s, ast.EOF, "", 1, 7)

	s = NewScanner("}{==;=+=,:.")
	ok(t, s, ast.RBRACE, "}", 1, 1)
	ok(t, s, ast.LBRACE, "{", 1, 2)
	ok(t, s, ast.DBL_EQ, "==", 1, 3)
	ok(t, s, ast.SEMICOLON, ";", 1, 5)
	ok(t, s, ast.EQ, "=", 1, 6)
	ok(t, s, ast.PLUS, "+", 1, 7)
	ok(t, s, ast.EQ, "=", 1, 8)
	ok(t, s, ast.COMMA, ",", 1, 9)
	ok(t, s, ast.COLON, ":", 1, 10)
	ok(t, s, ast.DOT, ".", 1, 11)
	ok(t, s, ast.EOF, "", 1, 12)

	s = NewScanner("! !=")
	ok(t, s, ast.NOT, "!", 1, 1)
	ok(t, s, ast.NOT_EQ, "!=", 1, 3)
	ok(t, s, ast.EOF, "", 1, 5)

	s = NewScanner("> >=")
	ok(t, s, ast.GT, ">", 1, 1)
	ok(t, s, ast.GT_EQ, ">=", 1, 3)
	ok(t, s, ast.EOF, "", 1, 5)

	s = NewScanner("< <= <=>")
	ok(t, s, ast.LT, "<", 1, 1)
	ok(t, s, ast.LT_EQ, "<=", 1, 3)
	ok(t, s, ast.CMP, "<=>", 1, 6)
	ok(t, s, ast.EOF, "", 1, 9)

	s = NewScanner("& && | ||")
	ok(t, s, ast.AMP, "&", 1, 1)
	ok(t, s, ast.DBL_AMP, "&&", 1, 3)
	ok(t, s, ast.PIPE, "|", 1, 6)
	ok(t, s, ast.DBL_PIPE, "||", 1, 8)
	ok(t, s, ast.EOF, "", 1, 10)
}

func TestInt(t *testing.T) {

	s := NewScanner("0")
	ok(t, s, ast.INT, "0", 1, 1)
	ok(t, s, ast.EOF, "", 1, 2)

	s = NewScanner("12+34 - 5 ")
	ok(t, s, ast.INT, "12", 1, 1)
	ok(t, s, ast.PLUS, "+", 1, 3)
	ok(t, s, ast.INT, "34", 1, 4)
	ok(t, s, ast.MINUS, "-", 1, 7)
	ok(t, s, ast.INT, "5", 1, 9)
	ok(t, s, ast.EOF, "", 1, 11)

	s = NewScanner("678")
	ok(t, s, ast.INT, "678", 1, 1)
	ok(t, s, ast.EOF, "", 1, 4)

	s = NewScanner("0 00")
	ok(t, s, ast.INT, "0", 1, 1)
	ok(t, s, ast.UNEXPECTED_CHAR, "0", 1, 4)

	s = NewScanner("00 1")
	ok(t, s, ast.UNEXPECTED_CHAR, "0", 1, 2)

	s = NewScanner("0xabcdef123456789")
	ok(t, s, ast.INT, "0xabcdef123456789", 1, 1)
	ok(t, s, ast.EOF, "", 1, 18)

	s = NewScanner("0xABCDEF")
	ok(t, s, ast.INT, "0xABCDEF", 1, 1)
	ok(t, s, ast.EOF, "", 1, 9)

	s = NewScanner("0x")
	ok(t, s, ast.UNEXPECTED_EOF, "", 1, 3)

	s = NewScanner("0xg")
	ok(t, s, ast.UNEXPECTED_CHAR, "g", 1, 3)
}

func TestFloat(t *testing.T) {
	s := NewScanner("0.12 0.34")
	ok(t, s, ast.FLOAT, "0.12", 1, 1)
	ok(t, s, ast.FLOAT, "0.34", 1, 6)
	ok(t, s, ast.EOF, "", 1, 10)

	s = NewScanner("12.34 56.78")
	ok(t, s, ast.FLOAT, "12.34", 1, 1)
	ok(t, s, ast.FLOAT, "56.78", 1, 7)
	ok(t, s, ast.EOF, "", 1, 12)

	s = NewScanner("0.34E1")
	ok(t, s, ast.FLOAT, "0.34E1", 1, 1)
	ok(t, s, ast.EOF, "", 1, 7)

	s = NewScanner("0.34E-1")
	ok(t, s, ast.FLOAT, "0.34E-1", 1, 1)
	ok(t, s, ast.EOF, "", 1, 8)

	s = NewScanner("0.34E+1")
	ok(t, s, ast.FLOAT, "0.34E+1", 1, 1)
	ok(t, s, ast.EOF, "", 1, 8)

	s = NewScanner("0.34e2")
	ok(t, s, ast.FLOAT, "0.34e2", 1, 1)
	ok(t, s, ast.EOF, "", 1, 7)

	s = NewScanner("0e6")
	ok(t, s, ast.FLOAT, "0e6", 1, 1)
	ok(t, s, ast.EOF, "", 1, 4)

	s = NewScanner("1e6")
	ok(t, s, ast.FLOAT, "1e6", 1, 1)
	ok(t, s, ast.EOF, "", 1, 4)

	s = NewScanner("0.")
	ok(t, s, ast.UNEXPECTED_EOF, "", 1, 3)
	s = NewScanner("0. ")
	ok(t, s, ast.UNEXPECTED_CHAR, " ", 1, 3)

	s = NewScanner("0.1e")
	ok(t, s, ast.UNEXPECTED_EOF, "", 1, 5)
	s = NewScanner("0.1e ")
	ok(t, s, ast.UNEXPECTED_CHAR, " ", 1, 5)
}

func TestStr(t *testing.T) {
	s := NewScanner("''")
	ok(t, s, ast.STR, "", 1, 1)
	ok(t, s, ast.EOF, "", 1, 3)

	s = NewScanner("'a'")
	ok(t, s, ast.STR, "a", 1, 1)
	ok(t, s, ast.EOF, "", 1, 4)

	s = NewScanner("\"\"")
	ok(t, s, ast.STR, "", 1, 1)
	ok(t, s, ast.EOF, "", 1, 3)

	s = NewScanner("\"a\"")
	ok(t, s, ast.STR, "a", 1, 1)
	ok(t, s, ast.EOF, "", 1, 4)

	s = NewScanner("'ab' 'c'")
	ok(t, s, ast.STR, "ab", 1, 1)
	ok(t, s, ast.STR, "c", 1, 6)
	ok(t, s, ast.EOF, "", 1, 9)

	s = NewScanner("'ab")
	ok(t, s, ast.UNEXPECTED_EOF, "", 1, 4)

	s = NewScanner("'\n'")
	ok(t, s, ast.UNEXPECTED_CHAR, "\n", 2, 0)

	s = NewScanner("'\\'\\n\\r\\t\\\\'")
	ok(t, s, ast.STR, "'\n\r\t\\", 1, 1)
	ok(t, s, ast.EOF, "", 1, 13)
}

func TestIdentOrKeyword(t *testing.T) {
	s := NewScanner("a bar")
	ok(t, s, ast.IDENT, "a", 1, 1)
	ok(t, s, ast.IDENT, "bar", 1, 3)
	ok(t, s, ast.EOF, "", 1, 6)

	s = NewScanner("_ zork")
	ok(t, s, ast.BLANK_IDENT, "_", 1, 1)
	ok(t, s, ast.IDENT, "zork", 1, 3)
	ok(t, s, ast.EOF, "", 1, 7)

	s = NewScanner("null true false")
	ok(t, s, ast.NULL, "null", 1, 1)
	ok(t, s, ast.TRUE, "true", 1, 6)
	ok(t, s, ast.FALSE, "false", 1, 11)
	ok(t, s, ast.EOF, "", 1, 16)

	s = NewScanner("if else")
	ok(t, s, ast.IF, "if", 1, 1)
	ok(t, s, ast.ELSE, "else", 1, 4)
	ok(t, s, ast.EOF, "", 1, 8)

	s = NewScanner("while break continue")
	ok(t, s, ast.WHILE, "while", 1, 1)
	ok(t, s, ast.BREAK, "break", 1, 7)
	ok(t, s, ast.CONTINUE, "continue", 1, 13)
	ok(t, s, ast.EOF, "", 1, 21)

	s = NewScanner("fn return const let ")
	ok(t, s, ast.FN, "fn", 1, 1)
	ok(t, s, ast.RETURN, "return", 1, 4)
	ok(t, s, ast.CONST, "const", 1, 11)
	ok(t, s, ast.LET, "let", 1, 17)
	ok(t, s, ast.EOF, "", 1, 21)

	s = NewScanner("obj this")
	ok(t, s, ast.OBJ, "obj", 1, 1)
	ok(t, s, ast.THIS, "this", 1, 5)
	ok(t, s, ast.EOF, "", 1, 9)
}

func TestComments(t *testing.T) {

	s := NewScanner("1 //foo\n2")
	ok(t, s, ast.INT, "1", 1, 1)
	ok(t, s, ast.INT, "2", 2, 1)
	ok(t, s, ast.EOF, "", 2, 2)

	s = NewScanner("1 2 //foo")
	ok(t, s, ast.INT, "1", 1, 1)
	ok(t, s, ast.INT, "2", 1, 3)
	ok(t, s, ast.EOF, "", 1, 10)

	s = NewScanner("1 /*foo*/2")
	ok(t, s, ast.INT, "1", 1, 1)
	ok(t, s, ast.INT, "2", 1, 10)
	ok(t, s, ast.EOF, "", 1, 11)

	s = NewScanner("1 2/**/")
	ok(t, s, ast.INT, "1", 1, 1)
	ok(t, s, ast.INT, "2", 1, 3)
	ok(t, s, ast.EOF, "", 1, 8)

	s = NewScanner("1 /*")
	ok(t, s, ast.INT, "1", 1, 1)
	ok(t, s, ast.UNEXPECTED_EOF, "", 1, 5)

	s = NewScanner("1 /* *")
	ok(t, s, ast.INT, "1", 1, 1)
	ok(t, s, ast.UNEXPECTED_EOF, "", 1, 7)

}

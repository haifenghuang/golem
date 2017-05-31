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

package core

import (
	"bytes"
	"fmt"
)

//---------------------------------------------------------------
// An StructDef contains the information needed to instantiate an Struct
// instance.  StructDefs are created at compile time, and
// are immutable at run time.

type StructDef []string

//---------------------------------------------------------------
// Ref

type Module struct {
	Pool       []Basic
	Refs       []*Ref
	StructDefs []StructDef
	Templates  []*Template

	Symbols map[string]*Symbol
}

func (m *Module) String() string {
	var buf bytes.Buffer
	buf.WriteString("----------------------------\n")
	buf.WriteString("Module:\n")

	buf.WriteString("    Pool:\n")
	for i, val := range m.Pool {
		typeOf := val.TypeOf()
		buf.WriteString("        ")
		buf.WriteString(fmt.Sprintf("%d: %v(%v)\n", i, typeOf, val))
	}

	buf.WriteString("    Refs:\n")
	for i, ref := range m.Refs {
		buf.WriteString("        ")
		buf.WriteString(fmt.Sprintf("%d: %v\n", i, ref))
	}

	buf.WriteString("    StructDefs:\n")
	for i, def := range m.StructDefs {
		buf.WriteString("        ")
		buf.WriteString(fmt.Sprintf("%d: %v\n", i, def))
	}

	for i, t := range m.Templates {

		buf.WriteString(fmt.Sprintf(
			"    Template(%d): Arity: %d, NumCaptures: %d, NumLocals: %d\n",
			i, t.Arity, t.NumCaptures, t.NumLocals))

		buf.WriteString("        OpCodes:\n")
		for i := 0; i < len(t.OpCodes); {
			text := FmtOpcode(t.OpCodes, i)
			buf.WriteString("            ")
			buf.WriteString(text)
			i += OpCodeSize(t.OpCodes[i])
		}

		buf.WriteString("        LineNumberTable:\n")
		for _, ln := range t.LineNumberTable {
			buf.WriteString("            ")
			buf.WriteString(fmt.Sprintf("%v\n", ln))
		}

		buf.WriteString("        ExceptionHandlers:\n")
		for _, eh := range t.ExceptionHandlers {
			buf.WriteString("            ")
			buf.WriteString(fmt.Sprintf("%v\n", eh))
		}
	}

	return buf.String()
}

//---------------------------------------------------------------
// A Ref is a reference, i.e. a container for a value

type Ref struct {
	Val Value
}

func NewRef(val Value) *Ref {
	return &Ref{val}
}

func (r *Ref) String() string {
	return fmt.Sprintf("Ref(%v)", r.Val)
}

//---------------------------------------------------------------
// A Symbol is a Ref that is visible outside the module

type Symbol struct {
	RefIndex int
	IsConst  bool
}

func NewSymbol(refIndex int, isConst bool) *Symbol {
	return &Symbol{refIndex, isConst}
}

func (s *Symbol) String() string {
	return fmt.Sprintf("Symbol(%d,%v)", s.RefIndex, s.IsConst)
}

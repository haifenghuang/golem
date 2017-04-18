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

type Module struct {
	Pool      []Value
	Locals    []*Ref
	ObjDefs   []*ObjDef
	Templates []*Template
}

func (m *Module) String() string {
	var buf bytes.Buffer
	buf.WriteString("----------------------------\n")
	buf.WriteString("Module:\n")

	buf.WriteString("    Pool:\n")
	for i, val := range m.Pool {
		tp, _ := val.TypeOf()
		buf.WriteString("        ")
		buf.WriteString(fmt.Sprintf("%d: %v(%v)\n", i, tp, val))
	}

	buf.WriteString("    Locals:\n")
	for i, ref := range m.Locals {
		buf.WriteString("        ")
		buf.WriteString(fmt.Sprintf("%d: %v\n", i, ref))
	}

	buf.WriteString("    ObjDefs:\n")
	for i, def := range m.ObjDefs {
		buf.WriteString("        ")
		buf.WriteString(fmt.Sprintf("%d: %v\n", i, def.Keys))
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

		buf.WriteString("        OpcLines:\n")
		for _, opl := range t.OpcLines {
			buf.WriteString("            ")
			buf.WriteString(fmt.Sprintf("%v\n", opl))
		}
	}

	return buf.String()
}

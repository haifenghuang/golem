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
	"fmt"
)

//---------------------------------------------------------------
// Template represents the information needed to invoke a function
// instance.  Templates are created at compile time, and
// are immutable at run time.

type Template struct {
	Arity       int
	NumCaptures int
	NumLocals   int
	OpCodes     []byte
	OpcLines    []OpcLine
}

// Return the line number for the opcode at the given instruction pointer
func (t *Template) LineNumber(instPtr int) int {

	oln := t.OpcLines
	n := len(oln) - 1

	for i := 0; i < n; i++ {
		if (instPtr >= oln[i].Index) && (instPtr < oln[i+1].Index) {
			return oln[i].LineNum
		}
	}
	return oln[n].LineNum
}

//---------------------------------------------------------------
// OpcLine tracks which sequence of opcodes are on a given line

type OpcLine struct {
	Index   int
	LineNum int
}

//---------------------------------------------------------------
// _func represents an instance of a function

type _func struct {
	template *Template
	captures []*Ref
}

// Called via NEW_FUNC opcode at runtime
func NewFunc(template *Template) Func {
	captures := make([]*Ref, 0, template.NumCaptures)
	return &_func{template, captures}
}

func (f *_func) TypeOf() (Type, Error) { return TFUNC, nil }

func (f *_func) String() (Str, Error) {
	return MakeStr(f.doStr()), nil
}

func (f *_func) doStr() string {
	return fmt.Sprintf("func(%p)", f)
}

func (f *_func) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *_func:
		if f.doStr() == t.doStr() {
			return TRUE, nil
		} else {
			return FALSE, nil
		}
	default:
		return FALSE, nil
	}
}

func (f *_func) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (f *_func) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat([]Value{f, t})

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (f *_func) Template() *Template {
	return f.template
}

func (f *_func) GetCapture(idx int) *Ref {
	return f.captures[idx]
}

func (f *_func) PushCapture(ref *Ref) {
	f.captures = append(f.captures, ref)
}

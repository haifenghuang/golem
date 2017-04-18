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

// OpcLine tracks which sequence of opcodes are on a given line
type OpcLine struct {
	OpcIndex int
	LineNum  int
}

//---------------------------------------------------------------
// Func represents an instance of a function

type Func struct {
	Template *Template
	Captures []*Ref
}

// Called via NEW_FUNC opcode at runtime
func NewFunc(template *Template) *Func {
	captures := make([]*Ref, 0, template.NumCaptures)
	return &Func{template, captures}
}

func (f *Func) TypeOf() (Type, Error) { return TFUNC, nil }

func (f *Func) String() (Str, Error) {
	return Str(f.doStr()), nil
}

func (f *Func) doStr() string {
	return fmt.Sprintf("Func(%p)", f)
}

func (f *Func) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *Func:
		if f.doStr() == t.doStr() {
			return Bool(true), nil
		} else {
			return Bool(false), nil
		}
	default:
		return Bool(false), nil
	}
}

func (f *Func) Cmp(v Value) (Int, Error) { return Int(0), ExpectedCmpError() }

func (f *Func) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat([]Value{f, t})

	default:
		return nil, ExpectedNumberError()
	}
}

func (f *Func) Sub(v Value) (Number, Error) { return nil, ExpectedNumberError() }
func (f *Func) Mul(v Value) (Number, Error) { return nil, ExpectedNumberError() }
func (f *Func) Div(v Value) (Number, Error) { return nil, ExpectedNumberError() }

func (f *Func) Negate() (Number, Error) { return Int(0), ExpectedNumberError() }
func (f *Func) Not() (Bool, Error)      { return false, ExpectedBoolError() }

func (f *Func) Select(key string) (Value, Error) { return nil, ExpectedObjError() }
func (f *Func) Put(key string, val Value) Error  { return ExpectedObjError() }

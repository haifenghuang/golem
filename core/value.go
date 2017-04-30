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
// Value

type Value interface {
	TypeOf() (Type, Error)

	HashCode() (Int, Error)
	Eq(Value) (Bool, Error)
	ToStr() (Str, Error)

	Cmp(Value) (Int, Error)
	Add(Value) (Value, Error)
}

//---------------------------------------------------------------
// Shared Interfaces

type (
	Getable interface {
		Get(Value) (Value, Error)
	}

	Indexable interface {
		Getable
		Set(Value, Value) Error
	}

	Lenable interface {
		Len() (Int, Error)
	}

	Sliceable interface {
		Slice(Value, Value) (Value, Error)
		SliceFrom(Value) (Value, Error)
		SliceTo(Value) (Value, Error)
	}

	Iterable interface {
		NewIterator() Iterator
	}
)

//---------------------------------------------------------------
// Basic

type (
	Basic interface {
		Value
		basicMarker()
	}

	Null interface {
		Basic
	}

	Bool interface {
		Basic
		BoolVal() bool

		Not() Bool
	}

	Str interface {
		Basic
		fmt.Stringer

		Getable
		Lenable
		Sliceable
		Iterable
	}

	Number interface {
		Basic
		FloatVal() float64
		IntVal() int64

		Sub(Value) (Number, Error)
		Mul(Value) (Number, Error)
		Div(Value) (Number, Error)
		Negate() (Number, Error)
	}

	Float interface {
		Number
	}

	Int interface {
		Number

		Rem(Value) (Int, Error)
		BitAnd(Value) (Int, Error)
		BitOr(Value) (Int, Error)
		BitXOr(Value) (Int, Error)
		LeftShift(Value) (Int, Error)
		RightShift(Value) (Int, Error)
		Complement() (Int, Error)
	}
)

//---------------------------------------------------------------
// Composite

type (
	Composite interface {
		Value
		compositeMarker()
	}

	List interface {
		Composite
		Indexable
		Lenable
		Sliceable
		Iterable

		Append(Value) Error
	}

	Range interface {
		Composite
		Getable
		Lenable
		Sliceable
		Iterable

		From() Int
		To() Int
		Step() Int
	}

	Tuple interface {
		Composite
		Getable
		Lenable
	}

	Dict interface {
		Composite
		Indexable
		Lenable
		Iterable
	}

	Obj interface {
		Composite
		Indexable

		Init(*ObjDef, []Value)

		GetField(Str) (Value, Error)
		PutField(Str, Value) Error
		Has(Value) (Bool, Error)
	}

	Iterator interface {
		Obj
		IterNext() Bool
		IterGet() (Value, Error)
	}
)

//---------------------------------------------------------------
// Func

type (

	// Func represents an instance of a function
	Func interface {
		Value
	}

	// BytecodeFunc represents a function that is defined
	// via Golem source code
	BytecodeFunc interface {
		Func

		Template() *Template
		GetCapture(int) *Ref
		PushCapture(*Ref)
	}

	// NativeFunc represents a function that is defined
	// natively within Go.
	NativeFunc interface {
		Func

		Invoke([]Value) (Value, Error)
	}
)

type (
	// Template represents the information needed to invoke a function
	// instance.  Templates are created at compile time, and
	// are immutable at run time.
	Template struct {
		Arity       int
		NumCaptures int
		NumLocals   int
		OpCodes     []byte
		OpcLines    []OpcLine
	}

	// OpcLine tracks which sequence of opcodes are on a ven line
	OpcLine struct {
		Index   int
		LineNum int
	}
)

//---------------------------------------------------------------
// Type

type Type int

const (
	TNULL Type = iota
	TBOOL
	TSTR
	TINT
	TFLOAT
	TFUNC
	TLIST
	TRANGE
	TTUPLE
	TDICT
	TOBJ
)

func (t Type) String() string {
	switch t {
	case TNULL:
		return "Null"
	case TBOOL:
		return "Bool"
	case TSTR:
		return "Str"
	case TINT:
		return "Int"
	case TFLOAT:
		return "Float"
	case TFUNC:
		return "Func"
	case TLIST:
		return "List"
	case TRANGE:
		return "Range"
	case TTUPLE:
		return "Tuple"
	case TDICT:
		return "Dict"
	case TOBJ:
		return "Obj"

	default:
		panic("unreachable")
	}
}

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
	case TOBJ:
		return "Obj"

	default:
		panic("unreachable")
	}
}

//---------------------------------------------------------------
// Shared Functionality

// Getable: Str, List, Obj
// Lenable: Str, List
// Setable: List, Obj
// Sliceable: List, Str

type (
	Getable interface {
		Get(Value) (Value, Error)
	}

	//Setable interface {
	//	Set(Value, Value) Error
	//}

	//Lenable interface {
	//	Len() (Int, Error)
	//}

	//Sliceable interface {
	//	Slice(Value, Value) (Value, Error)
	//	SliceFrom(Value) (Value, Error)
	//	SliceTo(Value) (Value, Error)
	//}
)

//---------------------------------------------------------------
// Value

type (
	Value interface {
		TypeOf() (Type, Error)

		Eq(Value) (Bool, Error)
		String() (Str, Error)
		Cmp(Value) (Int, Error)
		Add(Value) (Value, Error)
	}

	Null interface {
		Value
	}

	Bool interface {
		Value
		BoolVal() bool

		Not() Bool
	}

	Str interface {
		Value
		StrVal() string

		Getable
	}

	Number interface {
		Value
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

	List interface {
		Value

		Getable
		Set(Value, Value) Error
		Len() (Int, Error)

		Append(Value) Error
	}

	Obj interface {
		Value
		Init(*ObjDef, []Value)

		Getable

		GetField(Str) (Value, Error)
		PutField(Str, Value) Error
	}

	Func interface {
		Value

		Template() *Template
		GetCapture(int) *Ref
		PushCapture(*Ref)
	}
)

//---------------------------------------------------------------
// Ref

type Ref struct {
	Val Value
}

func NewRef(val Value) *Ref {
	return &Ref{val}
}

func (r *Ref) String() string {
	return fmt.Sprintf("Ref(%v)", r.Val)
}

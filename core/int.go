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
	//	"strings"
)

type Int int64

var ZERO Int = Int(0)
var ONE Int = Int(1)
var NEG_ONE Int = Int(-1)

func (i Int) TypeOf() (Type, Error) { return TINT, nil }

func (i Int) String() (Str, Error) {
	return MakeStr(fmt.Sprintf("%d", i)), nil
}

func (i Int) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {

	case Int:
		return MakeBool(i == t), nil

	case Float:
		a := float64(i)
		b := t.FloatVal()
		return MakeBool(a == b), nil

	default:
		return FALSE, nil
	}
}

func (i Int) Cmp(v Value) (Int, Error) {
	switch t := v.(type) {

	case Int:
		if i < t {
			return Int(-1), nil
		} else if i > t {
			return Int(1), nil
		} else {
			return Int(0), nil
		}

	case Float:
		a := float64(i)
		b := t.FloatVal()
		if a < b {
			return -1, nil
		} else if a > b {
			return 1, nil
		} else {
			return 0, nil
		}

	default:
		return Int(0), TypeMismatchError("Expected Comparable Type")
	}
}

func (i Int) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat([]Value{i, t})

	case Int:
		return i + t, nil

	case Float:
		a := float64(i)
		b := t.FloatVal()
		return MakeFloat(a + b), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (i Int) Sub(v Value) (Number, Error) {
	switch t := v.(type) {

	case Int:
		return i - t, nil

	case Float:
		a := float64(i)
		b := t.FloatVal()
		return MakeFloat(a - b), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (i Int) Mul(v Value) (Number, Error) {
	switch t := v.(type) {

	case Int:
		return i * t, nil

	case Float:
		a := float64(i)
		b := t.FloatVal()
		return MakeFloat(a * b), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (i Int) Div(v Value) (Number, Error) {
	switch t := v.(type) {

	case Int:
		if t == 0 {
			return nil, DivideByZeroError()
		} else {
			return i / t, nil
		}

	case Float:
		a := float64(i)
		b := t.FloatVal()
		if b == 0.0 {
			return nil, DivideByZeroError()
		} else {
			return MakeFloat(a / b), nil
		}

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (i Int) Negate() (Number, Error) {
	return 0 - i, nil
}

func (i Int) Rem(v Value) (Int, Error) {
	switch t := v.(type) {
	case Int:
		return i % t, nil
	default:
		return Int(0), TypeMismatchError("Expected 'Int'")
	}
}
func (i Int) BitAnd(v Value) (Int, Error) {
	switch t := v.(type) {
	case Int:
		return i & t, nil
	default:
		return Int(0), TypeMismatchError("Expected 'Int'")
	}
}
func (i Int) BitOr(v Value) (Int, Error) {
	switch t := v.(type) {
	case Int:
		return i | t, nil
	default:
		return Int(0), TypeMismatchError("Expected 'Int'")
	}
}
func (i Int) BitXOr(v Value) (Int, Error) {
	switch t := v.(type) {
	case Int:
		return i ^ t, nil
	default:
		return Int(0), TypeMismatchError("Expected 'Int'")
	}
}
func (i Int) LeftShift(v Value) (Int, Error) {
	switch t := v.(type) {
	case Int:
		if t < 0 {
			return Int(0), InvalidArgumentError("Shift count cannot be less than zero")
		} else {
			return i << uint(t), nil
		}
	default:
		return Int(0), TypeMismatchError("Expected 'Int'")
	}
}

func (i Int) RightShift(v Value) (Int, Error) {
	switch t := v.(type) {
	case Int:
		if t < 0 {
			return Int(0), InvalidArgumentError("Shift count cannot be less than zero")
		} else {
			return i >> uint(t), nil
		}
	default:
		return Int(0), TypeMismatchError("Expected 'Int'")
	}
}

func (i Int) Complement() (Int, Error) {
	return ^i, nil
}

func (i Int) GetField(key string) (Value, Error)   { return nil, TypeMismatchError("Expected 'Obj'") }
func (i Int) PutField(key string, val Value) Error { return TypeMismatchError("Expected 'Obj'") }

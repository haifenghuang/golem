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

type _float float64

func (f _float) FloatVal() float64 {
	return float64(f)
}

func MakeFloat(f float64) Float {
	return _float(f)
}

func (f _float) TypeOf() (Type, Error) { return TFLOAT, nil }

func (f _float) String() (Str, Error) {
	return MakeStr(fmt.Sprintf("%g", f)), nil
}

func (f _float) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {

	case _float:
		return MakeBool(f == t), nil

	case Int:
		return MakeBool(f == _float(t)), nil

	default:
		return FALSE, nil
	}
}

func (f _float) Cmp(v Value) (Int, Error) {
	switch t := v.(type) {

	case _float:
		if f < t {
			return -1, nil
		} else if f > t {
			return 1, nil
		} else {
			return 0, nil
		}

	case Int:
		g := _float(t)
		if f < g {
			return -1, nil
		} else if f > g {
			return 1, nil
		} else {
			return 0, nil
		}

	default:
		return 0, TypeMismatchError("Expected Comparable Type")
	}
}

func (f _float) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat([]Value{f, t})

	case Int:
		return f + _float(t), nil

	case _float:
		return f + t, nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (f _float) Sub(v Value) (Number, Error) {
	switch t := v.(type) {

	case Int:
		return f - _float(t), nil

	case _float:
		return f - t, nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (f _float) Mul(v Value) (Number, Error) {
	switch t := v.(type) {

	case Int:
		return f * _float(t), nil

	case _float:
		return f * t, nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (f _float) Div(v Value) (Number, Error) {
	switch t := v.(type) {

	case Int:
		if t == 0 {
			return nil, DivideByZeroError()
		} else {
			return f / _float(t), nil
		}

	case _float:
		if t == 0.0 {
			return nil, DivideByZeroError()
		} else {
			return f / t, nil
		}

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (f _float) Negate() (Number, Error) {
	return 0 - f, nil
}

func (f _float) Rem(v Value) (Int, Error)       { return Int(0), TypeMismatchError("Expected 'Int'") }
func (f _float) BitAnd(v Value) (Int, Error)    { return Int(0), TypeMismatchError("Expected 'Int'") }
func (f _float) BitOr(v Value) (Int, Error)     { return Int(0), TypeMismatchError("Expected 'Int'") }
func (f _float) BitXOr(v Value) (Int, Error)    { return Int(0), TypeMismatchError("Expected 'Int'") }
func (f _float) LeftShift(v Value) (Int, Error) { return Int(0), TypeMismatchError("Expected 'Int'") }
func (f _float) RightShift(Value) (Int, Error)  { return Int(0), TypeMismatchError("Expected 'Int'") }

func (f _float) Complement() (Int, Error) { return Int(0), TypeMismatchError("Expected 'Int'") }

func (f _float) GetField(key string) (Value, Error)   { return nil, TypeMismatchError("Expected 'Obj'") }
func (f _float) PutField(key string, val Value) Error { return TypeMismatchError("Expected 'Obj'") }

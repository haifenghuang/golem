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

type Float float64

func (f Float) TypeOf() (Type, Error) { return TFLOAT, nil }

func (f Float) number() {}

func (f Float) String() (Str, Error) {
	return Str(fmt.Sprintf("%g", f)), nil
}

func (f Float) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {

	case Float:
		return f == t, nil

	case Int:
		return f == Float(t), nil

	default:
		return false, nil
	}
}

func (f Float) Cmp(v Value) (Int, Error) {
	switch t := v.(type) {

	case Float:
		if f < t {
			return -1, nil
		} else if f > t {
			return 1, nil
		} else {
			return 0, nil
		}

	case Int:
		g := Float(t)
		if f < g {
			return -1, nil
		} else if f > g {
			return 1, nil
		} else {
			return 0, nil
		}

	default:
		return 0, ExpectedCmpError()
	}
}

func (f Float) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat([]Value{f, t})

	case Int:
		return f + Float(t), nil

	case Float:
		return f + t, nil

	case *Null:
		return nil, NullValueError()

	default:
		return nil, ExpectedNumberError()
	}
}

func (f Float) Sub(v Value) (Number, Error) {
	switch t := v.(type) {

	case Int:
		return f - Float(t), nil

	case Float:
		return f - t, nil

	case *Null:
		return nil, NullValueError()

	default:
		return nil, ExpectedNumberError()
	}
}

func (f Float) Mul(v Value) (Number, Error) {
	switch t := v.(type) {

	case Int:
		return f * Float(t), nil

	case Float:
		return f * t, nil

	case *Null:
		return nil, NullValueError()

	default:
		return nil, ExpectedNumberError()
	}
}

func (f Float) Div(v Value) (Number, Error) {
	switch t := v.(type) {

	case Int:
		if t == 0 {
			return nil, DivideByZeroError()
		} else {
			return f / Float(t), nil
		}

	case Float:
		if t == 0.0 {
			return nil, DivideByZeroError()
		} else {
			return f / t, nil
		}

	case *Null:
		return nil, NullValueError()

	default:
		return nil, ExpectedNumberError()
	}
}

func (f Float) Rem(v Value) (Int, Error)       { return Int(0), ExpectedIntError() }
func (f Float) BitAnd(v Value) (Int, Error)    { return Int(0), ExpectedIntError() }
func (f Float) BitOr(v Value) (Int, Error)     { return Int(0), ExpectedIntError() }
func (f Float) BitXOr(v Value) (Int, Error)    { return Int(0), ExpectedIntError() }
func (f Float) LeftShift(v Value) (Int, Error) { return Int(0), ExpectedIntError() }
func (f Float) RightShift(Value) (Int, Error)  { return Int(0), ExpectedIntError() }

func (f Float) Negate() (Number, Error) {
	return 0 - f, nil
}
func (f Float) Not() (Bool, Error)       { return false, ExpectedBoolError() }
func (f Float) Complement() (Int, Error) { return Int(0), ExpectedIntError() }

func (f Float) Select(key string) (Value, Error) { return nil, ExpectedObjError() }
func (f Float) Put(key string, val Value) Error  { return ExpectedObjError() }

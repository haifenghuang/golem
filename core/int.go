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

func (i Int) TypeOf() (Type, Error) { return TINT, nil }

func (i Int) number() {}

func (i Int) String() (Str, Error) {
	return Str(fmt.Sprintf("%d", i)), nil
}

func (i Int) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {

	case Int:
		return i == t, nil

	case Float:
		j := Int(t)
		return i == j, nil

	default:
		return false, nil
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
		j := Int(t)
		if i < j {
			return -1, nil
		} else if i > j {
			return 1, nil
		} else {
			return 0, nil
		}

	default:
		return Int(0), ExpectedCmpError()
	}
}

func (i Int) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat([]Value{i, t})

	case Int:
		return i + t, nil

	case Float:
		return Float(i) + t, nil

	case *Null:
		return nil, NullValueError()

	default:
		return nil, ExpectedNumberError()
	}
}

func (i Int) Sub(v Value) (Number, Error) {
	switch t := v.(type) {

	case Int:
		return i - t, nil

	case Float:
		return Float(i) - t, nil

	case *Null:
		return nil, NullValueError()

	default:
		return nil, ExpectedNumberError()
	}
}

func (i Int) Mul(v Value) (Number, Error) {
	switch t := v.(type) {

	case Int:
		return i * t, nil

	case Float:
		return Float(i) * t, nil

	case *Null:
		return nil, NullValueError()

	default:
		return nil, ExpectedNumberError()
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
		if t == 0.0 {
			return nil, DivideByZeroError()
		} else {
			return Float(i) / t, nil
		}

	case *Null:
		return nil, NullValueError()

	default:
		return nil, ExpectedNumberError()
	}
}

func (i Int) Rem(v Value) (Int, Error)       { return Int(0), ExpectedIntError() }
func (i Int) BitAnd(v Value) (Int, Error)    { return Int(0), ExpectedIntError() }
func (i Int) BitOr(v Value) (Int, Error)     { return Int(0), ExpectedIntError() }
func (i Int) BitXOr(v Value) (Int, Error)    { return Int(0), ExpectedIntError() }
func (i Int) LeftShift(v Value) (Int, Error) { return Int(0), ExpectedIntError() }
func (i Int) RightShift(Value) (Int, Error)  { return Int(0), ExpectedIntError() }

func (i Int) Negate() (Number, Error) {
	return 0 - i, nil
}

func (i Int) Not() (Bool, Error) { return false, ExpectedBoolError() }

func (i Int) Complement() (Int, Error) { return Int(0), ExpectedIntError() }

func (i Int) Select(key string) (Value, Error) { return nil, ExpectedObjError() }
func (i Int) Put(key string, val Value) Error  { return ExpectedObjError() }

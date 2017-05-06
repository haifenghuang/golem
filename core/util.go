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

import ()

// Parse an index value.
// The value must be between 0 (inclusive) and max (exclusive).
func ParseIndex(val Value, max int) (Int, Error) {
	if i, ok := val.(Int); ok {
		n := int(i.IntVal())
		if (n < 0) || (n >= max) {
			return nil, IndexOutOfBoundsError()
		} else {
			return i, nil
		}
	} else {
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

func Strcat(a Value, b Value) (Str, Error) {

	sa := fromValue(a)
	sb := fromValue(b)

	// copy to avoid memory leaks
	ca := make([]rune, len(sa))
	copy(ca, sa)

	cb := make([]rune, len(sb))
	copy(cb, sb)

	result := make(str, 0, len(ca)+len(cb))
	result = append(result, ca...)
	result = append(result, cb...)
	return result, nil
}

func fromValue(v Value) str {
	if sv, ok := v.(str); ok {
		return sv
	} else {
		return toRunes(v.ToStr().String())
	}
}

func valuesEq(as []Value, bs []Value) Bool {

	if len(as) != len(bs) {
		return FALSE
	}

	for i, a := range as {
		if a.Eq(bs[i]) == FALSE {
			return FALSE
		}
	}

	return TRUE
}

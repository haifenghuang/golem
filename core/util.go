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
)

// Parse an index value.
// The value must be between 0 (inclusive) and max (exclusive)
func parseIndex(val Value, max int) (Int, Error) {
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

// Concatenate two values into a string
func strcat(a Value, b Value) (Str, Error) {
	as, err := a.String()
	if err != nil {
		return nil, err
	}

	bs, err := b.String()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString(as.StrVal())
	buf.WriteString(bs.StrVal())
	return MakeStr(buf.String()), nil
}

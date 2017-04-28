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
	//"fmt"
	"testing"
)

func TestFunc(t *testing.T) {

	a := NewBytecodeFunc(&Template{})
	b := NewBytecodeFunc(&Template{})

	okType(t, a, TFUNC)

	z, err := a.Eq(b)
	ok(t, z, err, FALSE)
	z, err = b.Eq(a)
	ok(t, z, err, FALSE)
	z, err = a.Eq(a)
	ok(t, z, err, TRUE)
}

func TestLineNumber(t *testing.T) {

	tp := &Template{0, 0, 0, nil,
		[]OpcLine{
			OpcLine{0, 0},
			OpcLine{1, 2},
			OpcLine{11, 3},
			OpcLine{20, 4},
			OpcLine{29, 0}}}

	assert(t, tp.LineNumber(0) == 0)
	assert(t, tp.LineNumber(1) == 2)
	assert(t, tp.LineNumber(10) == 2)
	assert(t, tp.LineNumber(11) == 3)
	assert(t, tp.LineNumber(19) == 3)
	assert(t, tp.LineNumber(20) == 4)
	assert(t, tp.LineNumber(28) == 4)
	assert(t, tp.LineNumber(29) == 0)
}

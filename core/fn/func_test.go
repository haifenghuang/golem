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

package fn

import (
	//"fmt"
	g "golem/core"
	"reflect"
	"testing"
)

func assert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
	}
}

func ok(t *testing.T, val g.Value, err g.Error, expect g.Value) {

	if err != nil {
		t.Error(err, " != ", nil)
	}

	if !reflect.DeepEqual(val, expect) {
		panic("asdfasfad")
		t.Error(val, " != ", expect)
	}
}

func fail(t *testing.T, val g.Value, err g.Error, expect string) {

	if val != nil {
		t.Error(val, " != ", nil)
	}

	if err.Error() != expect {
		t.Error(err.Error(), " != ", expect)
	}
}

func okType(t *testing.T, val g.Value, expected g.Type) {
	tp, err := val.TypeOf()
	assert(t, tp == expected)
	assert(t, err == nil)
}

func TestFunc(t *testing.T) {

	a := NewBytecodeFunc(&Template{})
	b := NewBytecodeFunc(&Template{})

	okType(t, a, g.TFUNC)

	z, err := a.Eq(b)
	ok(t, z, err, g.FALSE)
	z, err = b.Eq(a)
	ok(t, z, err, g.FALSE)
	z, err = a.Eq(a)
	ok(t, z, err, g.TRUE)
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

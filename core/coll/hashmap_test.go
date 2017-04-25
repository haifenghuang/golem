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

package coll

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

func debug(hm *HashMap) {
	//fmt.Println("--------------------------")
	//fmt.Printf("size: %d\n", hm.size)
	//for i, b := range hm.buckets {
	//	fmt.Print(i, ": [")
	//	for j, e := range b {
	//		if j > 0 {
	//			fmt.Print(", ")
	//		}
	//		fmt.Printf("(%v:%v)", e.Key, e.Value)
	//	}
	//	fmt.Println("]")
	//}
	//fmt.Println("--------------------------")
}

func TestHashMap(t *testing.T) {
	hm := NewHashMap(nil)
	debug(hm)

	ok(t, hm.Len(), nil, g.ZERO)
	v, err := hm.Get(g.MakeInt(3))
	ok(t, v, err, g.NULL)

	err = hm.Put(g.MakeInt(3), g.MakeInt(33))
	ok(t, nil, err, nil)
	debug(hm)

	ok(t, hm.Len(), nil, g.ONE)
	v, err = hm.Get(g.MakeInt(3))
	ok(t, v, err, g.MakeInt(33))
	v, err = hm.Get(g.MakeInt(5))
	ok(t, v, err, g.NULL)

	err = hm.Put(g.MakeInt(3), g.MakeInt(33))
	ok(t, nil, err, nil)
	debug(hm)

	ok(t, hm.Len(), nil, g.ONE)
	v, err = hm.Get(g.MakeInt(3))
	ok(t, v, err, g.MakeInt(33))
	v, err = hm.Get(g.MakeInt(5))
	ok(t, v, err, g.NULL)

	err = hm.Put(g.MakeInt(int64(2)), g.MakeInt(int64(22)))
	ok(t, nil, err, nil)
	debug(hm)
	ok(t, hm.Len(), nil, g.MakeInt(2))

	err = hm.Put(g.MakeInt(int64(1)), g.MakeInt(int64(11)))
	ok(t, nil, err, nil)
	debug(hm)
	ok(t, hm.Len(), nil, g.MakeInt(3))

	for i := 1; i <= 20; i++ {
		err = hm.Put(g.MakeInt(int64(i)), g.MakeInt(int64(i*10+i)))
		ok(t, nil, err, nil)
	}
	debug(hm)

	for i := 1; i <= 40; i++ {
		v, err = hm.Get(g.MakeInt(int64(i)))
		if i <= 20 {
			ok(t, v, err, g.MakeInt(int64(i*10+i)))
		} else {
			ok(t, v, err, g.NULL)
		}
	}
}

func TestStrHashMap(t *testing.T) {

	hm := NewHashMap(nil)

	err := hm.Put(g.MakeStr("abc"), g.MakeStr("xyz"))
	ok(t, nil, err, nil)

	v, err := hm.Get(g.MakeStr("abc"))
	ok(t, v, err, g.MakeStr("xyz"))
}

func testIteratorEntries(t *testing.T, initial []*HEntry, expect []*HEntry) {

	hm := NewHashMap(initial)

	entries := []*HEntry{}
	itr := hm.Iterator()
	for itr.Next() {
		entries = append(entries, itr.Get())
	}

	if !reflect.DeepEqual(entries, expect) {
		t.Error("iterator failed")
	}
}

func TestHashMapIterator(t *testing.T) {

	testIteratorEntries(t,
		[]*HEntry{},
		[]*HEntry{})

	testIteratorEntries(t,
		[]*HEntry{
			&HEntry{g.MakeStr("a"), g.MakeInt(1)}},
		[]*HEntry{
			&HEntry{g.MakeStr("a"), g.MakeInt(1)}})

	testIteratorEntries(t,
		[]*HEntry{
			&HEntry{g.MakeStr("a"), g.MakeInt(1)},
			&HEntry{g.MakeStr("b"), g.MakeInt(2)}},
		[]*HEntry{
			&HEntry{g.MakeStr("b"), g.MakeInt(2)},
			&HEntry{g.MakeStr("a"), g.MakeInt(1)}})

	testIteratorEntries(t,
		[]*HEntry{
			&HEntry{g.MakeStr("a"), g.MakeInt(1)},
			&HEntry{g.MakeStr("b"), g.MakeInt(2)},
			&HEntry{g.MakeStr("c"), g.MakeInt(3)}},
		[]*HEntry{
			&HEntry{g.MakeStr("b"), g.MakeInt(2)},
			&HEntry{g.MakeStr("a"), g.MakeInt(1)},
			&HEntry{g.MakeStr("c"), g.MakeInt(3)}})
}

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
	"reflect"
	"testing"
)

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

	ok(t, hm.Len(), nil, ZERO)
	v, err := hm.Get(MakeInt(3))
	ok(t, v, err, NULL)

	err = hm.Put(MakeInt(3), MakeInt(33))
	ok(t, nil, err, nil)
	debug(hm)

	ok(t, hm.Len(), nil, ONE)
	v, err = hm.Get(MakeInt(3))
	ok(t, v, err, MakeInt(33))
	v, err = hm.Get(MakeInt(5))
	ok(t, v, err, NULL)

	err = hm.Put(MakeInt(3), MakeInt(33))
	ok(t, nil, err, nil)
	debug(hm)

	ok(t, hm.Len(), nil, ONE)
	v, err = hm.Get(MakeInt(3))
	ok(t, v, err, MakeInt(33))
	v, err = hm.Get(MakeInt(5))
	ok(t, v, err, NULL)

	err = hm.Put(MakeInt(int64(2)), MakeInt(int64(22)))
	ok(t, nil, err, nil)
	debug(hm)
	ok(t, hm.Len(), nil, MakeInt(2))

	err = hm.Put(MakeInt(int64(1)), MakeInt(int64(11)))
	ok(t, nil, err, nil)
	debug(hm)
	ok(t, hm.Len(), nil, MakeInt(3))

	for i := 1; i <= 20; i++ {
		err = hm.Put(MakeInt(int64(i)), MakeInt(int64(i*10+i)))
		ok(t, nil, err, nil)
	}
	debug(hm)

	for i := 1; i <= 40; i++ {
		v, err = hm.Get(MakeInt(int64(i)))
		if i <= 20 {
			ok(t, v, err, MakeInt(int64(i*10+i)))
		} else {
			ok(t, v, err, NULL)
		}
	}
}

func TestStrHashMap(t *testing.T) {

	hm := NewHashMap(nil)

	err := hm.Put(MakeStr("abc"), MakeStr("xyz"))
	ok(t, nil, err, nil)

	v, err := hm.Get(MakeStr("abc"))
	ok(t, v, err, MakeStr("xyz"))

	v, err = hm.ContainsKey(MakeStr("abc"))
	ok(t, v, err, TRUE)

	v, err = hm.ContainsKey(MakeStr("bogus"))
	ok(t, v, err, FALSE)
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
			{MakeStr("a"), MakeInt(1)}},
		[]*HEntry{
			{MakeStr("a"), MakeInt(1)}})

	testIteratorEntries(t,
		[]*HEntry{
			{MakeStr("a"), MakeInt(1)},
			{MakeStr("b"), MakeInt(2)}},
		[]*HEntry{
			{MakeStr("b"), MakeInt(2)},
			{MakeStr("a"), MakeInt(1)}})

	testIteratorEntries(t,
		[]*HEntry{
			{MakeStr("a"), MakeInt(1)},
			{MakeStr("b"), MakeInt(2)},
			{MakeStr("c"), MakeInt(3)}},
		[]*HEntry{
			{MakeStr("b"), MakeInt(2)},
			{MakeStr("a"), MakeInt(1)},
			{MakeStr("c"), MakeInt(3)}})
}

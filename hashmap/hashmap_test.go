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

package hashmap

import (
	//"fmt"
	g "golem/core"
	"testing"
)

func assert(t *testing.T, flag bool) {
	if !flag {
		panic("assertion failure")
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

	assert(t, hm.Len() == 0)
	assert(t, hm.Get(g.MakeInt(3)) == nil)

	hm.Put(g.MakeInt(3), g.MakeInt(33))
	debug(hm)

	assert(t, hm.Len() == 1)
	assert(t, hm.Get(g.MakeInt(3)) == g.MakeInt(33))
	assert(t, hm.Get(g.MakeInt(5)) == nil)

	hm.Put(g.MakeInt(3), g.MakeInt(33))
	debug(hm)

	assert(t, hm.Len() == 1)
	assert(t, hm.Get(g.MakeInt(3)) == g.MakeInt(33))
	assert(t, hm.Get(g.MakeInt(5)) == nil)

	hm.Put(g.MakeInt(int64(2)), g.MakeInt(int64(22)))
	debug(hm)
	assert(t, hm.Len() == 2)

	hm.Put(g.MakeInt(int64(1)), g.MakeInt(int64(11)))
	debug(hm)
	assert(t, hm.Len() == 3)

	for i := 1; i <= 20; i++ {
		hm.Put(g.MakeInt(int64(i)), g.MakeInt(int64(i*10+i)))
	}
	debug(hm)

	for i := 1; i <= 40; i++ {
		if i <= 20 {
			assert(t, hm.Get(g.MakeInt(int64(i))) != nil)
		} else {
			assert(t, hm.Get(g.MakeInt(int64(i))) == nil)
		}
	}
}

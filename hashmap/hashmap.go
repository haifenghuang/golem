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
)

// A Custom HashMap implementation.  This allows us
// to use things like []rune as a key into a hash map.

type (
	HashMap struct {
		buckets []bucket
		size    int
	}

	Entry struct {
		Key   g.Value
		Value g.Value
	}

	bucket []*Entry
)

func NewHashMap(entries []Entry) *HashMap {
	capacity := 5
	buckets := make([]bucket, capacity, capacity)
	hm := &HashMap{buckets, 0}

	for _, e := range entries {
		hm.Put(e.Key, e.Value)
	}
	return hm
}

func (hm *HashMap) Get(key g.Value) (value g.Value, err g.Error) {

	// panic-recover is the cleanest approach
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(g.Error); ok {
				value = nil
				err = e
			}
			panic(r)
		}
	}()

	b := hm.buckets[hm.hashBucket(key)]
	n := indexOf(b, key)
	if n == -1 {
		return nil, nil
	} else {
		return b[n].Value, nil
	}
}

func (hm *HashMap) Put(key g.Value, value g.Value) (err g.Error) {

	// panic-recover is the cleanest approach
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(g.Error); ok {
				err = e
			}
			panic(r)
		}
	}()

	h := hm.hashBucket(key)
	n := indexOf(hm.buckets[h], key)
	if n == -1 {
		if hm.tooFull() {
			hm.rehash()
			h = hm.hashBucket(key)
		}
		hm.buckets[h] = append(hm.buckets[h], &Entry{key, value})
		hm.size++

	} else {
		hm.buckets[h][n].Value = value
	}

	return nil
}

func (hm *HashMap) Len() g.Int {
	return g.MakeInt(int64(hm.size))
}

func indexOf(b bucket, key g.Value) int {
	for i, e := range b {

		// panic-recover is the cleanest approach
		eq, err := e.Key.Eq(key)
		if err != nil {
			panic(err)
		}

		if eq.BoolVal() {
			return i
		}
	}
	return -1
}

func (hm *HashMap) tooFull() bool {
	headroom := (hm.size + 1) << 1
	return headroom > len(hm.buckets)
}

func (hm *HashMap) rehash() {
	oldBuckets := hm.buckets

	capacity := len(hm.buckets)<<1 + 1
	hm.buckets = make([]bucket, capacity, capacity)
	for _, b := range oldBuckets {
		for _, e := range b {
			h := hm.hashBucket(e.Key)
			hm.buckets[h] = append(hm.buckets[h], e)
		}
	}
}

func (hm *HashMap) hashBucket(key g.Value) int {

	// panic-recover is the cleanest approach
	hc, err := key.HashCode()
	if err != nil {
		panic(err)
	}

	hv := int(hc.IntVal())
	if hv < 0 {
		hv = 0 - hv
	}

	return hv % len(hm.buckets)
}

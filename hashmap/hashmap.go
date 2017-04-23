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

type (
	HashMap struct {
		buckets []bucket
		size    int
	}

	Entry struct {
		Key   g.HashKey
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

func (hm *HashMap) Get(key g.HashKey) g.Value {
	h := key.HashCode() % len(hm.buckets)
	b := hm.buckets[h]
	n := indexOf(b, key)
	if n == -1 {
		return nil
	} else {
		return b[n].Value
	}
}

func (hm *HashMap) Put(key g.HashKey, value g.Value) {
	h := key.HashCode() % len(hm.buckets)
	b := hm.buckets[h]
	n := indexOf(b, key)
	if n == -1 {
		if hm.tooFull() {
			hm.rehash()
			h = key.HashCode() % len(hm.buckets)
			b = hm.buckets[h]
			n = indexOf(b, key)
		}
		hm.buckets[h] = append(b, &Entry{key, value})
		hm.size++

	} else {
		b[n].Value = value
	}
}

func (hm *HashMap) Len() int {
	return hm.size
}

func indexOf(b bucket, key g.HashKey) int {
	for i, e := range b {
		if e.Key == key {
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
			h := e.Key.HashCode() % len(hm.buckets)
			b := hm.buckets[h]
			hm.buckets[h] = append(b, e)
		}
	}
}

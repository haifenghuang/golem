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

package comp

import (
	g "golem/core"
)

type (
	Composite interface {
		g.Value
		compositeMarker()
	}

	List interface {
		Composite
		g.Indexable
		g.Lenable
		g.Sliceable

		Append(g.Value) g.Error
	}

	Tuple interface {
		Composite
		g.Getable
		g.Lenable
	}

	Obj interface {
		Composite
		g.Indexable

		Init(*ObjDef, []g.Value)

		GetField(g.Str) (g.Value, g.Error)
		PutField(g.Str, g.Value) g.Error
		Has(g.Value) (g.Bool, g.Error)
	}

	Dict interface {
		Composite
		g.Indexable
		g.Lenable
	}
)

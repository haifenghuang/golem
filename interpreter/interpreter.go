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

package interpreter

import (
	"fmt"
	g "golem/core"
)

//---------------------------------------------------------------
// The Golem Interpreter

type Interpreter struct {
	mod    *g.Module
	frames []*frame
}

func NewInterpreter(mod *g.Module) *Interpreter {
	tpl := mod.Templates[0]
	if tpl.Arity != 0 || tpl.NumCaptures != 0 {
		panic("TODO")
	}

	return &Interpreter{mod, []*frame{}}
}

func (i *Interpreter) Init() (g.Value, *errorTrace) {

	// use the zeroth template
	tpl := i.mod.Templates[0]

	// create empty locals
	i.mod.Locals = newLocals(tpl.NumLocals, nil)

	// make func
	fn := g.NewBytecodeFunc(tpl)

	// push a frame
	i.frames = append(i.frames, &frame{fn, i.mod.Locals, []g.Value{}, 0})

	// go
	return i.run()
}

func (i *Interpreter) run() (result g.Value, errTrace *errorTrace) {

	var err g.Error
	for result == nil {
		result, err = i.advance(0)
		if err != nil {

			instPtr := i.frames[len(i.frames)-1].ip
			errTrace = &errorTrace{err, i.stackTrace()}

			// look for catch and finally clauses
			for j := len(i.frames) - 1; j >= 0; j-- {
				f := i.frames[j]
				tpl := f.fn.Template()

				// check each excpeptionHandler in backwards order so that
				// nested try blocks works properly
				for k := len(tpl.ExceptionHandlers) - 1; k >= 0; k-- {
					eh := tpl.ExceptionHandlers[k]

					// found an active finally clause
					if eh.Finally != -1 && (instPtr >= eh.Begin && instPtr < eh.End) {
						// The error that we end up returning could be an error
						// thrown from a finally, rather than the original error.
						_, ferr := i.runFinally(f, j, eh)
						if ferr != nil {
							errTrace = &errorTrace{ferr, i.stackTrace()}
						}
					}
				}
			}

			// TODO: catch recovery
			return nil, errTrace
		}
	}

	return result, nil
}

func (i *Interpreter) runFinally(
	f *frame, frameIndex int, eh g.ExceptionHandler) (result g.Value, err g.Error) {

	f.ip = eh.Finally
	for result == nil {
		result, err = i.advance(frameIndex)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (i *Interpreter) stackTrace() []string {

	n := len(i.frames)
	stack := []string{}

	for j := n - 1; j >= 0; j-- {
		tpl := i.frames[j].fn.Template()
		lineNum := tpl.LineNumber(i.frames[j].ip)
		stack = append(stack, fmt.Sprintf("    at line %d", lineNum))
	}

	return stack
}

func newLocals(numLocals int, params []g.Value) []*g.Ref {
	p := len(params)
	locals := make([]*g.Ref, numLocals, numLocals)
	for j := 0; j < numLocals; j++ {
		if j < p {
			locals[j] = &g.Ref{params[j]}
		} else {
			locals[j] = &g.Ref{g.NULL}
		}
	}
	return locals
}

func (i *Interpreter) dump() {

	println("-----------------------------------------")

	f := i.frames[len(i.frames)-1]
	opc := f.fn.Template().OpCodes
	print(g.FmtOpcode(opc, f.ip))

	for j, f := range i.frames {
		fmt.Printf("frame %d\n", j)
		f.dump()
	}
}

//---------------------------------------------------------------
// An execution environment, a.k.a 'stack frame'.

type frame struct {
	fn     g.BytecodeFunc
	locals []*g.Ref
	stack  []g.Value
	ip     int
}

func (f *frame) dump() {
	fmt.Printf("    locals:\n")
	for j, r := range f.locals {
		fmt.Printf("        %d: %s\n", j, r.Val.ToStr())
	}
	fmt.Printf("    stack:\n")
	for j, v := range f.stack {
		fmt.Printf("        %d: %s\n", j, v.ToStr())
	}
	fmt.Printf("    ip: %d\n", f.ip)
}

//---------------------------------------------------------------
// A combination of an error, and a stack trace

type errorTrace struct {
	err        g.Error
	stackTrace []string
}

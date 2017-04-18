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
// An execution environment, a.k.a 'stack frame'.

type frame struct {
	function *g.Func
	locals   []*g.Ref
	stack    []g.Value
	instPtr  int
}

//---------------------------------------------------------------
// A ErrorStack is returned when there is an unrecoverable error

type ErrorStack struct {
	Err   g.Error
	Stack []string
}

//---------------------------------------------------------------
// The Interpreter

type Interpreter struct {
	mod *g.Module
}

func NewInterpreter(mod *g.Module) *Interpreter {
	tpl := mod.Templates[0]
	if tpl.Arity != 0 || tpl.NumCaptures != 0 {
		panic("TODO")
	}

	return &Interpreter{mod}
}

func (inp *Interpreter) Init() (g.Value, *ErrorStack) {

	// use the zeroth template
	tpl := inp.mod.Templates[0]

	// create empty locals
	locals := newLocals(tpl.NumLocals, nil)
	inp.mod.Locals = locals

	// make func
	fn := g.NewFunc(tpl)

	// go
	return inp.invoke(fn, locals)
}

func (inp *Interpreter) invoke(fn *g.Func, locals []*g.Ref) (g.Value, *ErrorStack) {

	pool := inp.mod.Pool
	defs := inp.mod.ObjDefs
	frames := []*frame{}
	opc := fn.Template.OpCodes

	// stack and instruction pointer
	s := []g.Value{}
	ip := 0

	// loop over giant switch
	for {
		n := len(s) - 1

		switch opc[ip] {

		case g.INVOKE:

			idx := index(opc, ip)

			// get function from stack
			operand, ok := s[n-idx].(*g.Func)
			if !ok {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(g.ExpectedFuncError(), frames)
			}
			if operand.Template.Arity != idx {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(g.ArityMismatchError(operand.Template.Arity, idx), frames)
			}

			// get params from stack
			params := s[n-idx+1:]

			// remove function and params from stack
			s = s[:n-idx]

			// move the instruction pointer
			ip += 3

			// save the execution environment
			frames = append(frames, &frame{fn, locals, s, ip})

			// create a new execution environment
			fn = operand
			locals = newLocals(fn.Template.NumLocals, params)
			opc = fn.Template.OpCodes
			s = []g.Value{}
			ip = 0

		case g.RETURN:

			// get result from top of stack
			result := s[n]

			// if there is no frame to pop, then we are done
			if len(frames) == 0 {
				return result, nil
			}

			// pop the frame
			lf := len(frames) - 1
			fr := frames[lf]
			frames = frames[:lf]

			// restore the execution environment
			fn = fr.function
			locals = fr.locals
			opc = fn.Template.OpCodes
			s = fr.stack
			ip = fr.instPtr

			// push the result
			s = append(s, result)

		case g.NEW_FUNC:

			// push a function
			idx := index(opc, ip)
			tpl := inp.mod.Templates[idx]
			nf := g.NewFunc(tpl)
			s = append(s, nf)
			ip += 3

		case g.FUNC_LOCAL:

			// get function from stack
			operand, ok := s[n].(*g.Func)
			if !ok {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(g.ExpectedFuncError(), frames)
			}

			// push a local onto the captures of the function
			idx := index(opc, ip)
			operand.Captures = append(operand.Captures, locals[idx])
			ip += 3

		case g.FUNC_CAPTURE:

			// get function from stack
			operand, ok := s[n].(*g.Func)
			if !ok {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(g.ExpectedFuncError(), frames)
			}

			// push a capture onto the captures of the function
			idx := index(opc, ip)
			operand.Captures = append(operand.Captures, fn.Captures[idx])
			ip += 3

		case g.NEW_OBJ:
			s = append(s, g.NewObj())
			ip++

		case g.INIT_OBJ:

			// look up ObjDef
			def := defs[index(opc, ip)]
			size := len(def.Keys)

			// get obj and values
			obj, ok := s[n-size].(*g.Obj)
			if !ok {
				panic("Invalid INIT_OBJ")
			}
			vals := s[n-size+1:]

			// initialize object
			obj.Init(def, vals)

			// pop values
			s = s[:n-size+1]

			// done
			ip += 3

		case g.SELECT:

			idx := index(opc, ip)
			key, ok := pool[idx].(g.Str)
			if !ok {
				panic("Invalid SELECT Key")
			}

			ks, err := key.String()
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}

			result, err := s[n].Select(string(ks))
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}

			s[n] = result
			ip += 3

		case g.PUT:

			idx := index(opc, ip)
			key, ok := pool[idx].(g.Str)
			if !ok {
				panic("Invalid PUT Key")
			}

			operand := s[n-1]
			value := s[n]

			ks, err := key.String()
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}

			err = operand.Put(string(ks), value)
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}

			s[n-1] = value
			s = s[:n]
			ip += 3

		case g.LOAD_NULL:
			s = append(s, g.NULL)
			ip++
		case g.LOAD_TRUE:
			s = append(s, g.TRUE)
			ip++
		case g.LOAD_FALSE:
			s = append(s, g.FALSE)
			ip++

		case g.LOAD_CONST:
			idx := index(opc, ip)
			s = append(s, pool[idx])
			ip += 3

		case g.LOAD_LOCAL:
			idx := index(opc, ip)
			s = append(s, locals[idx].Val)
			ip += 3

		case g.LOAD_CAPTURE:
			idx := index(opc, ip)
			s = append(s, fn.Captures[idx].Val)
			ip += 3

		case g.STORE_LOCAL:
			idx := index(opc, ip)
			locals[idx].Val = s[n]
			s = s[:n]
			ip += 3

		case g.STORE_CAPTURE:
			idx := index(opc, ip)
			fn.Captures[idx].Val = s[n]
			s = s[:n]
			ip += 3

		case g.JUMP:
			ip = index(opc, ip)

		case g.JUMP_TRUE:
			b, ok := s[n].(g.Bool)
			if !ok {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(g.ExpectedBoolError(), frames)
			}

			s = s[:n]
			if b {
				ip = index(opc, ip)
			} else {
				ip += 3
			}

		case g.JUMP_FALSE:
			b, ok := s[n].(g.Bool)
			if !ok {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(g.ExpectedBoolError(), frames)
			}

			s = s[:n]
			if b {
				ip += 3
			} else {
				ip = index(opc, ip)
			}

		case g.EQ:
			val, err := s[n-1].Eq(s[n])
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.NE:
			val, err := s[n-1].Eq(s[n])
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s = s[:n]
			s[n-1] = !val
			ip++

		case g.LT:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s = s[:n]
			s[n-1] = g.Bool(val < 0)
			ip++

		case g.LTE:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s = s[:n]
			s[n-1] = g.Bool(val <= 0)
			ip++

		case g.GT:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s = s[:n]
			s[n-1] = g.Bool(val > 0)
			ip++

		case g.GTE:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s = s[:n]
			s[n-1] = g.Bool(val >= 0)
			ip++

		case g.CMP:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.ADD:
			val, err := s[n-1].Add(s[n])
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.SUB:
			val, err := s[n-1].Sub(s[n])
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.MUL:
			val, err := s[n-1].Mul(s[n])
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.DIV:
			val, err := s[n-1].Div(s[n])
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.NEGATE:
			val, err := s[n].Negate()
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s[n] = val
			ip++

		case g.NOT:
			val, err := s[n].Not()
			if err != nil {
				frames = append(frames, &frame{fn, locals, s, ip})
				return nil, errorStack(err, frames)
			}
			s[n] = val
			ip++

		case g.DUP:
			s = append(s, s[n])
			ip++

		default:
			panic("Invalid opcode")
		}
	}
}

func errorStack(err g.Error, frames []*frame) *ErrorStack {
	n := len(frames)
	stack := make([]string, n)
	for i := n - 1; i >= 0; i-- {
		tp := frames[i].function.Template
		lineNum := tp.LineNumber(frames[i].instPtr)
		stack = append(stack, fmt.Sprintf("    at line %d", lineNum))
	}
	return &ErrorStack{err, stack}
}

func index(opcodes []byte, ip int) int {
	high := opcodes[ip+1]
	low := opcodes[ip+2]
	return int(high)<<8 + int(low)
}

func newLocals(numLocals int, params []g.Value) []*g.Ref {
	p := len(params)
	locals := make([]*g.Ref, numLocals, numLocals)
	for i := 0; i < numLocals; i++ {
		if i < p {
			locals[i] = &g.Ref{params[i]}
		} else {
			locals[i] = &g.Ref{g.NULL}
		}
	}
	return locals
}

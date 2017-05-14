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

// Advance the interpreter forwards by one opcode.
func (i *Interpreter) advance() g.Error {

	pool := i.mod.Pool
	f := i.frames[len(i.frames)-1]
	n := len(f.stack) - 1
	opc := f.fn.Template().OpCodes

	switch opc[f.ip] {

	case g.INVOKE:

		idx := index(opc, f.ip)
		params := f.stack[n-idx+1:]

		switch fn := f.stack[n-idx].(type) {
		case g.BytecodeFunc:

			/////////////////////////////////
			// invoke a bytecode-defined func

			// check arity
			arity := len(params)
			if arity != fn.Template().Arity {
				err := g.ArityMismatchError(
					fmt.Sprintf("%d", fn.Template().Arity), arity)
				return err
			}

			// pop from stack, and advance instruction pointer
			f.stack = f.stack[:n-idx]
			f.ip += 3

			// push a new frame
			locals := newLocals(fn.Template().NumLocals, params)
			i.frames = append(i.frames, &frame{fn, locals, []g.Value{}, 0})

		case g.NativeFunc:

			/////////////////////////////////
			// invoke a natively-defined func

			val, err := fn.Invoke(params) // Invoke() is responsible for arity check
			if err != nil {
				return err
			}

			f.stack = f.stack[:n-idx]
			f.stack = append(f.stack, val)
			f.ip += 3

		default:
			return g.TypeMismatchError("Expected 'Func'")
		}

	case g.RETURN:

		// TODO once we've written a Control Flow Graph
		// turn this sanity check on to make sure we are managing
		// the stack properly

		//if len(f.stack) < 1 || len(f.stack) > 2 {
		//	for j, v := range f.stack {
		//		fmt.Printf("stack %d: %s\n", j, v.ToStr())
		//	}
		//	panic("invalid stack")
		//}

		// If we would pop the last frame, then we are done.
		// (Don't pop the last frame! We still need it.)
		if len(i.frames) == 1 {
			i.done = true
		} else {
			// get result from top of stack
			result := f.stack[n]

			// pop the old frame
			i.frames = i.frames[:len(i.frames)-1]

			// push the result onto the new top frame
			f = i.frames[len(i.frames)-1]
			f.stack = append(f.stack, result)
		}

	case g.NEW_FUNC:

		// push a function
		idx := index(opc, f.ip)
		tpl := i.mod.Templates[idx]
		nf := g.NewBytecodeFunc(tpl)
		f.stack = append(f.stack, nf)
		f.ip += 3

	case g.FUNC_LOCAL:

		// get function from stack
		fn, ok := f.stack[n].(g.BytecodeFunc)
		if !ok {
			return g.TypeMismatchError("Expected 'BytecodeFunc'")
		}

		// push a local onto the captures of the function
		idx := index(opc, f.ip)
		fn.PushCapture(f.locals[idx])
		f.ip += 3

	case g.FUNC_CAPTURE:

		// get function from stack
		fn, ok := f.stack[n].(g.BytecodeFunc)
		if !ok {
			return g.TypeMismatchError("Expected 'BytecodeFunc'")
		}

		// push a capture onto the captures of the function
		idx := index(opc, f.ip)
		fn.PushCapture(f.fn.GetCapture(idx))
		f.ip += 3

	case g.NEW_OBJ:
		def := i.mod.StructDefs[index(opc, f.ip)]
		f.stack = append(f.stack, g.BlankStruct(def))
		f.ip += 3

	case g.NEW_LIST:

		size := index(opc, f.ip)
		vals := make([]g.Value, size)
		copy(vals, f.stack[n-size+1:])

		f.stack = f.stack[:n-size+1]
		f.stack = append(f.stack, g.NewList(vals))
		f.ip += 3

	case g.NEW_SET:

		size := index(opc, f.ip)
		vals := make([]g.Value, size)
		copy(vals, f.stack[n-size+1:])

		f.stack = f.stack[:n-size+1]
		f.stack = append(f.stack, g.NewSet(vals))
		f.ip += 3

	case g.NEW_TUPLE:

		size := index(opc, f.ip)
		vals := make([]g.Value, size)
		copy(vals, f.stack[n-size+1:])

		f.stack = f.stack[:n-size+1]
		f.stack = append(f.stack, g.NewTuple(vals))
		f.ip += 3

	case g.CHECK_TUPLE:

		// make sure the top of the stack is really a tuple
		tp, ok := f.stack[n].(g.Tuple)
		if !ok {
			return g.TypeMismatchError("Expected 'Tuple'")
		}

		// and make sure its of the expected length
		expectedLen := index(opc, f.ip)
		tpLen := tp.Len()
		if expectedLen != int(tpLen.IntVal()) {
			return g.InvalidArgumentError(
				fmt.Sprintf("Expected Tuple of length %d", expectedLen))
		}

		// do not alter stack
		f.ip += 3

	case g.NEW_DICT:

		size := index(opc, f.ip)
		entries := make([]*g.HEntry, 0, size)

		numVals := size * 2
		for j := n - numVals + 1; j <= n; j += 2 {
			entries = append(entries, &g.HEntry{f.stack[j], f.stack[j+1]})
		}

		f.stack = f.stack[:n-numVals+1]
		f.stack = append(f.stack, g.NewDict(entries))
		f.ip += 3

	case g.GET_FIELD:

		idx := index(opc, f.ip)
		key, ok := pool[idx].(g.Str)
		if !ok {
			panic("Invalid GET_FIELD Key")
		}

		result, err := f.stack[n].GetField(key)
		if err != nil {
			return err
		}

		f.stack[n] = result
		f.ip += 3

	case g.PUT_FIELD:

		idx := index(opc, f.ip)
		key, ok := pool[idx].(g.Str)
		if !ok {
			panic("Invalid PUT_FIELD Key")
		}

		// get struct from stack
		stc, ok := f.stack[n-1].(g.Struct)
		if !ok {
			return g.TypeMismatchError("Expected 'Struct'")
		}

		// get value from stack
		value := f.stack[n]

		err := stc.PutField(key, value)
		if err != nil {
			return err
		}

		f.stack[n-1] = value
		f.stack = f.stack[:n]
		f.ip += 3

	case g.INC_FIELD:

		idx := index(opc, f.ip)
		key, ok := pool[idx].(g.Str)
		if !ok {
			panic("Invalid GET_FIELD Key")
		}

		// get struct from stack
		stc, ok := f.stack[n-1].(g.Struct)
		if !ok {
			return g.TypeMismatchError("Expected 'Struct'")
		}

		// get value from stack
		value := f.stack[n]

		before, err := stc.GetField(key)
		if err != nil {
			return err
		}

		after, err := before.Plus(value)
		if err != nil {
			return err
		}

		err = stc.PutField(key, after)
		if err != nil {
			return err
		}

		f.stack[n-1] = before
		f.stack = f.stack[:n]
		f.ip += 3

	case g.GET_INDEX:

		// get Getable from stack
		gtb, ok := f.stack[n-1].(g.Getable)
		if !ok {
			return g.TypeMismatchError("Expected 'Getable'")
		}

		// get index from stack
		idx := f.stack[n]

		result, err := gtb.Get(idx)
		if err != nil {
			return err
		}

		f.stack[n-1] = result
		f.stack = f.stack[:n]
		f.ip++

	case g.SET_INDEX:

		// get Indexable from stack
		ibl, ok := f.stack[n-2].(g.Indexable)
		if !ok {
			return g.TypeMismatchError("Expected 'Getable'")
		}

		// get index from stack
		idx := f.stack[n-1]

		// get value from stack
		val := f.stack[n]

		err := ibl.Set(idx, val)
		if err != nil {
			return err
		}

		f.stack[n-2] = val
		f.stack = f.stack[:n-1]
		f.ip++

	case g.INC_INDEX:

		// get Indexable from stack
		ibl, ok := f.stack[n-2].(g.Indexable)
		if !ok {
			return g.TypeMismatchError("Expected 'Getable'")
		}

		// get index from stack
		idx := f.stack[n-1]

		// get value from stack
		val := f.stack[n]

		before, err := ibl.Get(idx)
		if err != nil {
			return err
		}

		after, err := before.Plus(val)
		if err != nil {
			return err
		}

		err = ibl.Set(idx, after)
		if err != nil {
			return err
		}

		f.stack[n-2] = before
		f.stack = f.stack[:n-1]
		f.ip++

	case g.SLICE:

		// get Sliceable from stack
		slb, ok := f.stack[n-2].(g.Sliceable)
		if !ok {
			return g.TypeMismatchError("Expected 'Sliceable'")
		}

		// get indices from stack
		from := f.stack[n-1]
		to := f.stack[n]

		result, err := slb.Slice(from, to)
		if err != nil {
			return err
		}

		f.stack[n-2] = result
		f.stack = f.stack[:n-1]
		f.ip++

	case g.SLICE_FROM:

		// get Sliceable from stack
		slb, ok := f.stack[n-1].(g.Sliceable)
		if !ok {
			return g.TypeMismatchError("Expected 'Sliceable'")
		}

		// get index from stack
		from := f.stack[n]

		result, err := slb.SliceFrom(from)
		if err != nil {
			return err
		}

		f.stack[n-1] = result
		f.stack = f.stack[:n]
		f.ip++

	case g.SLICE_TO:

		// get Sliceable from stack
		slb, ok := f.stack[n-1].(g.Sliceable)
		if !ok {
			return g.TypeMismatchError("Expected 'Sliceable'")
		}

		// get index from stack
		to := f.stack[n]

		result, err := slb.SliceTo(to)
		if err != nil {
			return err
		}

		f.stack[n-1] = result
		f.stack = f.stack[:n]
		f.ip++

	case g.LOAD_NULL:
		f.stack = append(f.stack, g.NULL)
		f.ip++
	case g.LOAD_TRUE:
		f.stack = append(f.stack, g.TRUE)
		f.ip++
	case g.LOAD_FALSE:
		f.stack = append(f.stack, g.FALSE)
		f.ip++
	case g.LOAD_ZERO:
		f.stack = append(f.stack, g.ZERO)
		f.ip++
	case g.LOAD_ONE:
		f.stack = append(f.stack, g.ONE)
		f.ip++
	case g.LOAD_NEG_ONE:
		f.stack = append(f.stack, g.NEG_ONE)
		f.ip++

	case g.LOAD_BUILTIN:
		idx := index(opc, f.ip)
		f.stack = append(f.stack, g.Builtins[idx])
		f.ip += 3

	case g.LOAD_CONST:
		idx := index(opc, f.ip)
		f.stack = append(f.stack, pool[idx])
		f.ip += 3

	case g.LOAD_LOCAL:
		idx := index(opc, f.ip)
		f.stack = append(f.stack, f.locals[idx].Val)
		f.ip += 3

	case g.LOAD_CAPTURE:
		idx := index(opc, f.ip)
		f.stack = append(f.stack, f.fn.GetCapture(idx).Val)
		f.ip += 3

	case g.STORE_LOCAL:
		idx := index(opc, f.ip)
		f.locals[idx].Val = f.stack[n]
		f.stack = f.stack[:n]
		f.ip += 3

	case g.STORE_CAPTURE:
		idx := index(opc, f.ip)
		f.fn.GetCapture(idx).Val = f.stack[n]
		f.stack = f.stack[:n]
		f.ip += 3

	case g.JUMP:
		f.ip = index(opc, f.ip)

	case g.JUMP_TRUE:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return g.TypeMismatchError("Expected 'Bool'")
		}

		f.stack = f.stack[:n]
		if b.BoolVal() {
			f.ip = index(opc, f.ip)
		} else {
			f.ip += 3
		}

	case g.JUMP_FALSE:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return g.TypeMismatchError("Expected 'Bool'")
		}

		f.stack = f.stack[:n]
		if b.BoolVal() {
			f.ip += 3
		} else {
			f.ip = index(opc, f.ip)
		}

	case g.EQ:
		b := f.stack[n-1].Eq(f.stack[n])
		f.stack = f.stack[:n]
		f.stack[n-1] = b
		f.ip++

	case g.NE:
		b := f.stack[n-1].Eq(f.stack[n])
		f.stack = f.stack[:n]
		f.stack[n-1] = b.Not()
		f.ip++

	case g.LT:
		val, err := f.stack[n-1].Cmp(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.MakeBool(val.IntVal() < 0)
		f.ip++

	case g.LTE:
		val, err := f.stack[n-1].Cmp(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.MakeBool(val.IntVal() <= 0)
		f.ip++

	case g.GT:
		val, err := f.stack[n-1].Cmp(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.MakeBool(val.IntVal() > 0)
		f.ip++

	case g.GTE:
		val, err := f.stack[n-1].Cmp(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.MakeBool(val.IntVal() >= 0)
		f.ip++

	case g.CMP:
		val, err := f.stack[n-1].Cmp(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.HAS:

		// get struct from stack
		stc, ok := f.stack[n-1].(g.Struct)
		if !ok {
			return g.TypeMismatchError("Expected 'Struct'")
		}

		val, err := stc.Has(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.PLUS:
		val, err := f.stack[n-1].Plus(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.NOT:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return g.TypeMismatchError("Expected 'Bool'")
		}

		f.stack[n] = b.Not()
		f.ip++

	case g.SUB:
		z, ok := f.stack[n-1].(g.Number)
		if !ok {
			return g.TypeMismatchError("Expected Number Type")
		}

		val, err := z.Sub(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.MUL:
		z, ok := f.stack[n-1].(g.Number)
		if !ok {
			return g.TypeMismatchError("Expected Number Type")
		}

		val, err := z.Mul(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.DIV:
		z, ok := f.stack[n-1].(g.Number)
		if !ok {
			return g.TypeMismatchError("Expected Number Type")
		}

		val, err := z.Div(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.NEGATE:
		z, ok := f.stack[n].(g.Number)
		if !ok {
			return g.TypeMismatchError("Expected Number Type")
		}

		val := z.Negate()
		f.stack[n] = val
		f.ip++

	case g.REM:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return g.TypeMismatchError("Expected 'Int'")
		}

		val, err := z.Rem(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.BIT_AND:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return g.TypeMismatchError("Expected 'Int'")
		}

		val, err := z.BitAnd(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.BIT_OR:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return g.TypeMismatchError("Expected 'Int'")
		}

		val, err := z.BitOr(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.BIT_XOR:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return g.TypeMismatchError("Expected 'Int'")
		}

		val, err := z.BitXOr(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.LEFT_SHIFT:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return g.TypeMismatchError("Expected 'Int'")
		}

		val, err := z.LeftShift(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.RIGHT_SHIFT:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return g.TypeMismatchError("Expected 'Int'")
		}

		val, err := z.RightShift(f.stack[n])
		if err != nil {
			return err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.COMPLEMENT:
		z, ok := f.stack[n].(g.Int)
		if !ok {
			return g.TypeMismatchError("Expected 'Int'")
		}

		val := z.Complement()
		f.stack[n] = val
		f.ip++

	case g.ITER:

		ibl, ok := f.stack[n].(g.Iterable)
		if !ok {
			panic("Expected Iterable")
		}

		f.stack[n] = ibl.NewIterator()
		f.ip++

	case g.ITER_NEXT:

		itr, ok := f.stack[n].(g.Iterator)
		if !ok {
			panic("Expected Iterator")
		}

		f.stack[n] = itr.IterNext()
		f.ip++

	case g.ITER_GET:

		itr, ok := f.stack[n].(g.Iterator)
		if !ok {
			panic("Expected Iterator")
		}

		val, err := itr.IterGet()
		if err != nil {
			return err
		}

		f.stack[n] = val
		f.ip++

	case g.DUP:
		f.stack = append(f.stack, f.stack[n])
		f.ip++

	case g.POP:
		f.stack = f.stack[:n]
		f.ip++

	default:
		panic("Invalid opcode")
	}

	return nil
}

func index(opcodes []byte, ip int) int {
	high := opcodes[ip+1]
	low := opcodes[ip+2]
	return int(high)<<8 + int(low)
}

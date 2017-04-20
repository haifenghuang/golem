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
	"fmt"
)

const (
	LOAD_NULL byte = iota
	LOAD_TRUE
	LOAD_FALSE
	LOAD_ZERO
	LOAD_ONE
	LOAD_NEG_ONE

	LOAD_CONST
	LOAD_LOCAL
	STORE_LOCAL
	LOAD_CAPTURE
	STORE_CAPTURE

	JUMP
	JUMP_TRUE
	JUMP_FALSE

	EQ
	NE
	GT
	GTE
	LT
	LTE
	CMP

	ADD
	SUB
	MUL
	DIV

	REM
	BIT_AND
	BIT_OR
	BIT_XOR
	LEFT_SHIFT
	RIGHT_SHIFT

	NEGATE
	NOT
	COMPLEMENT

	NEW_FUNC
	FUNC_CAPTURE
	FUNC_LOCAL

	INVOKE
	RETURN

	NEW_OBJ
	INIT_OBJ
	GET_FIELD
	PUT_FIELD

	DUP

	// These are temporary values created during compilation.
	// The interpreter will panic if it encounters them.
	BREAK    = 0xFD
	CONTINUE = 0xFE
)

func OpCodeSize(opc byte) int {

	switch opc {

	case LOAD_CONST, LOAD_LOCAL, LOAD_CAPTURE, STORE_LOCAL, STORE_CAPTURE,
		JUMP, JUMP_TRUE, JUMP_FALSE, BREAK, CONTINUE,
		NEW_FUNC, FUNC_CAPTURE, FUNC_LOCAL, INVOKE,
		INIT_OBJ, GET_FIELD, PUT_FIELD:

		return 3

	default:
		return 1
	}
}

func FmtOpcode(opcodes []byte, i int) string {

	switch opcodes[i] {

	case LOAD_NULL:
		return fmt.Sprintf("%d: LOAD_NULL\n", i)
	case LOAD_TRUE:
		return fmt.Sprintf("%d: LOAD_TRUE\n", i)
	case LOAD_FALSE:
		return fmt.Sprintf("%d: LOAD_FALSE\n", i)
	case LOAD_ZERO:
		return fmt.Sprintf("%d: LOAD_ZERO\n", i)
	case LOAD_ONE:
		return fmt.Sprintf("%d: LOAD_ONE\n", i)
	case LOAD_NEG_ONE:
		return fmt.Sprintf("%d: LOAD_NEG_ONE\n", i)

	case LOAD_CONST:
		return fmtIndex(opcodes, i, "LOAD_CONST")
	case LOAD_LOCAL:
		return fmtIndex(opcodes, i, "LOAD_LOCAL")
	case STORE_LOCAL:
		return fmtIndex(opcodes, i, "STORE_LOCAL")
	case LOAD_CAPTURE:
		return fmtIndex(opcodes, i, "LOAD_CAPTURE")
	case STORE_CAPTURE:
		return fmtIndex(opcodes, i, "STORE_CAPTURE")

	case JUMP:
		return fmtIndex(opcodes, i, "JUMP")
	case JUMP_TRUE:
		return fmtIndex(opcodes, i, "JUMP_TRUE")
	case JUMP_FALSE:
		return fmtIndex(opcodes, i, "JUMP_FALSE")

	case EQ:
		return fmt.Sprintf("%d: EQ\n", i)
	case NE:
		return fmt.Sprintf("%d: NE\n", i)
	case GT:
		return fmt.Sprintf("%d: GT\n", i)
	case GTE:
		return fmt.Sprintf("%d: GTE\n", i)
	case LT:
		return fmt.Sprintf("%d: LT\n", i)
	case LTE:
		return fmt.Sprintf("%d: LTE\n", i)
	case CMP:
		return fmt.Sprintf("%d: CMP\n", i)

	case ADD:
		return fmt.Sprintf("%d: ADD\n", i)
	case SUB:
		return fmt.Sprintf("%d: SUB\n", i)
	case MUL:
		return fmt.Sprintf("%d: MUL\n", i)
	case DIV:
		return fmt.Sprintf("%d: DIV\n", i)

	case NEGATE:
		return fmt.Sprintf("%d: NEGATE\n", i)
	case NOT:
		return fmt.Sprintf("%d: NOT\n", i)

	case NEW_FUNC:
		return fmtIndex(opcodes, i, "NEW_FUNC")
	case FUNC_CAPTURE:
		return fmtIndex(opcodes, i, "FUNC_CAPTURE")
	case FUNC_LOCAL:
		return fmtIndex(opcodes, i, "FUNC_LOCAL")

	case INVOKE:
		return fmtIndex(opcodes, i, "INVOKE")
	case RETURN:
		return fmt.Sprintf("%d: RETURN\n", i)

	case NEW_OBJ:
		return fmt.Sprintf("%d: NEW_OBJ\n", i)
	case INIT_OBJ:
		return fmtIndex(opcodes, i, "INIT_OBJ")
	case GET_FIELD:
		return fmtIndex(opcodes, i, "GET_FIELD")
	case PUT_FIELD:
		return fmtIndex(opcodes, i, "PUT_FIELD")

	case BREAK:
		return fmtIndex(opcodes, i, "BREAK")
	case CONTINUE:
		return fmtIndex(opcodes, i, "CONTINUE")

	case DUP:
		return fmt.Sprintf("%d: DUP\n", i)

	default:
		panic("unreachable")
	}
}

func fmtIndex(opcodes []byte, i int, tag string) string {
	high := opcodes[i+1]
	low := opcodes[i+2]
	index := int(high)<<8 + int(low)
	return fmt.Sprintf("%d: %s %d %d (%d)\n", i, tag, high, low, index)
}

// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package interpreter

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	script "github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/script/interpreter/errs"
)

// TestOpcodeDisabled tests the opcodeDisabled function manually because all
// disabled opcodes result in a script execution failure when executed normally,
// so the function is not called under normal circumstances.
func TestOpcodeDisabled(t *testing.T) {
	t.Parallel()

	tests := []byte{script.Op2MUL, script.Op2DIV}
	for _, opcodeVal := range tests {
		pop := ParsedOpcode{op: opcodeArray[opcodeVal], Data: nil}
		err := opcodeDisabled(&pop, nil)
		if !errs.IsErrorCode(err, errs.ErrDisabledOpcode) {
			t.Errorf("opcodeDisabled: unexpected error - got %v, "+
				"want %v", err, errs.ErrDisabledOpcode)
			continue
		}
	}
}

func TestParse(t *testing.T) {
	tt := []struct {
		name            string
		scriptHexString string

		expectedParsedScript ParsedScript
	}{
		{
			name:            "1 op return",
			scriptHexString: "0168776a0024dc",

			expectedParsedScript: ParsedScript{
				ParsedOpcode{
					op: opcode{
						val:    script.OpDATA1,
						name:   "OP_DATA_1",
						length: 2,
						exec:   opcodePushData,
					},
					Data: []byte{script.OpENDIF},
				},
				ParsedOpcode{
					op: opcode{
						val:    script.OpNIP,
						name:   "OP_NIP",
						length: 1,
						exec:   opcodeNip,
					},
					Data: nil,
				},
				ParsedOpcode{
					op: opcode{
						val:    script.OpRETURN,
						name:   "OP_RETURN",
						length: 1,
						exec:   opcodeReturn,
					},
					Data: nil,
				},
				ParsedOpcode{
					op: opcode{
						val:    0,
						name:   "Unformatted Data",
						length: 3,
					},
					Data: []byte{36, 220},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := script.NewFromHex(tc.scriptHexString)
			require.NoError(t, err)

			codeParser := DefaultOpcodeParser{}
			p, err := codeParser.Parse(s)
			require.NoError(t, err)

			for i := range p {
				require.Equal(t, tc.expectedParsedScript[i].Data, p[i].Data)
				require.Equal(t, tc.expectedParsedScript[i].op.length, p[i].op.length)
				require.Equal(t, tc.expectedParsedScript[i].op.name, p[i].op.name)
				require.Equal(t, tc.expectedParsedScript[i].op.val, p[i].op.val)
			}
		})
	}
}

func TestOpShift(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		op       byte
		initial  []byte
		shift    int64
		expected []byte
	}{
		{"RSHIFT 8 bits", script.OpRSHIFT, []byte{0x01, 0x02}, 8, []byte{0x02}},
		{"RSHIFT 1 bit", script.OpRSHIFT, []byte{0x01, 0x02}, 1, []byte{0x00, 0x01}},
		{"RSHIFT 9 bits", script.OpRSHIFT, []byte{0x01, 0x02}, 9, []byte{0x01}},
		{"RSHIFT 16 bits", script.OpRSHIFT, []byte{0x01, 0x02}, 16, []byte{}},
		{"RSHIFT 0 bits", script.OpRSHIFT, []byte{0x01, 0x02}, 0, []byte{0x01, 0x02}},
		{"RSHIFT large shift", script.OpRSHIFT, []byte{0x01, 0x02}, 100, []byte{}},
		{"RSHIFT 8 byte value by 1", script.OpRSHIFT, []byte{0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01}, 1, []byte{0xF7, 0xE6, 0xD5, 0xC4, 0xB3, 0xA2, 0x91, 0x00}},
		{"LSHIFT 0 bits", script.OpLSHIFT, []byte{0x01, 0x02}, 0, []byte{0x01, 0x02}},
		{"LSHIFT 1 bit", script.OpLSHIFT, []byte{0x01, 0x02}, 1, []byte{0x02, 0x04}},
		{"LSHIFT 8 bits", script.OpLSHIFT, []byte{0x01, 0x02}, 8, []byte{0x00, 0x01, 0x02}},
		{"LSHIFT 9 bits", script.OpLSHIFT, []byte{0x01, 0x02}, 9, []byte{0x00, 0x02, 0x04}},
		{"LSHIFT 15 bits", script.OpLSHIFT, []byte{0x01, 0x02}, 15, []byte{0x00, 0x80, 0x00, 0x01}},
		{"LSHIFT 16 bits", script.OpLSHIFT, []byte{0x01, 0x02}, 16, []byte{0x00, 0x00, 0x01, 0x02}},
		{"LSHIFT large shift", script.OpLSHIFT, []byte{0x01, 0x02}, 100, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x20}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new thread with necessary configuration
			thread := &thread{
				dstack:       newStack(&afterGenesisConfig{}, false),
				afterGenesis: true,                  // Set this according to your needs
				cfg:          &afterGenesisConfig{}, // Use appropriate config
				scriptParser: &DefaultOpcodeParser{},
			}

			// Push the initial value and shift amount onto the stack
			thread.dstack.PushByteArray(tt.initial)

			thread.dstack.PushInt(&scriptNumber{
				val:          big.NewInt(tt.shift),
				afterGenesis: thread.afterGenesis,
			})

			// Create a ParsedOpcode for the shift operation
			pop := ParsedOpcode{op: opcodeArray[tt.op], Data: nil}

			// Execute the opcode
			var opErr error
			if tt.op == script.OpLSHIFT {
				opErr = opcodeLShift(&pop, thread)
			} else {
				opErr = opcodeRShift(&pop, thread)
			}

			// Check for errors
			require.NoError(t, opErr)

			// Check the result
			result, err := thread.dstack.PopByteArray()
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

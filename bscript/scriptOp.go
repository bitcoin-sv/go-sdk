package bscript

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type ScriptOp struct {
	OpCode byte
	Data   []byte
}

func (op *ScriptOp) String() string {
	if op.OpCode > Op0 && op.OpCode <= OpPUSHDATA4 {
		return hex.EncodeToString(op.Data)
	}
	return OpCodeValues[op.OpCode]
}

// ReadOp reads the next script operation from the Script starting at the given position.
// It returns the parsed ScriptOp and any error encountered during parsing.
// The position is updated to point to the next operation in the Script.
func (s *Script) ReadOp(pos *int) (op *ScriptOp, err error) {
	b := *s
	if len(b) <= *pos {
		err = fmt.Errorf("ReadOp: %d %d", len(b), *pos)
		return
	}
	switch b[*pos] {
	case OpPUSHDATA1:
		if len(b) < *pos+2 {
			err = ErrDataTooSmall
			return
		}

		l := int(b[*pos+1])
		*pos += 2

		if len(b) < *pos+l {
			err = ErrDataTooSmall
			return
		}

		op = &ScriptOp{OpCode: OpPUSHDATA1, Data: b[*pos : *pos+l]}
		*pos += l

	case OpPUSHDATA2:
		if len(b) < *pos+3 {
			err = ErrDataTooSmall
			return
		}

		l := int(binary.LittleEndian.Uint16(b[*pos+1:]))
		*pos += 3

		if len(b) < *pos+l {
			err = ErrDataTooSmall
			return
		}

		op = &ScriptOp{OpCode: OpPUSHDATA2, Data: b[*pos : *pos+l]}
		*pos += l

	case OpPUSHDATA4:
		if len(b) < *pos+5 {
			err = ErrDataTooSmall
			return
		}

		l := int(binary.LittleEndian.Uint32(b[*pos+1:]))
		*pos += 5

		if len(b) < *pos+l {
			err = ErrDataTooSmall
			return
		}

		op = &ScriptOp{OpCode: OpPUSHDATA4, Data: b[*pos : *pos+l]}
		*pos += l

	default:
		if b[*pos] >= OpDATA1 && b[*pos] < OpPUSHDATA1 {
			l := b[*pos]
			if len(b) < *pos+int(1+l) {
				err = ErrDataTooSmall
				return
			}
			op = &ScriptOp{OpCode: b[*pos], Data: b[*pos+1 : *pos+int(l+1)]}
			*pos += int(1 + l)
		} else {
			op = &ScriptOp{OpCode: b[*pos]}
			*pos++
		}
	}

	return
}

// ParseOps parses the script and returns a slice of ScriptOp objects.
// It returns the parsed ScriptOps and any error encountered during parsing.
func (s *Script) ParseOps() (ops []*ScriptOp, err error) {
	pos := 0
	for pos < len(*s) {
		op, err := s.ReadOp(&pos)
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}

	return
}

// NewScriptFromScriptOps creates a new Script from a slice of ScriptOps.
// It returns the new Script and any error encountered during parsing.
func NewScriptFromScriptOps(parts []*ScriptOp) (*Script, error) {
	length := 0
	for _, p := range parts {
		// op code
		length++

		switch p.OpCode {
		case OpPUSHDATA1:
			// length of data + 1 byte for length
			length += 1 + len(p.Data)
		case OpPUSHDATA2:
			// length of data + 2 bytes for length
			length += 2 + len(p.Data)
		case OpPUSHDATA4:
			// length of data + 4 bytes for length
			length += 4 + len(p.Data)
		default:
			if p.OpCode >= OpDATA1 && p.OpCode < OpPUSHDATA1 {
				// length of data
				length += len(p.Data)
			}
		}
	}

	s := make(Script, 0, length)
	for _, p := range parts {
		if p.OpCode >= OpDATA1 && p.OpCode <= OpPUSHDATA4 {
			s.AppendPushData(p.Data)
		} else {
			s.AppendOpcodes(p.OpCode)
		}
	}
	return &s, nil
}

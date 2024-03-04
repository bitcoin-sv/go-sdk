package script

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/bits"
	"strings"

	"github.com/bitcoin-sv/go-sdk/util"
)

// ScriptKey types.
const (
	// TODO: change to p2pk/p2pkh
	ScriptTypePubKey                = "pubkey"
	ScriptTypePubKeyHash            = "pubkeyhash"
	ScriptTypeNonStandard           = "nonstandard"
	ScriptTypeEmpty                 = "empty"
	ScriptTypeMultiSig              = "multisig"
	ScriptTypeNullData              = "nulldata"
	ScriptTypePubKeyHashInscription = "pubkeyhashinscription"
)

// Script type
type Script []byte

// NewFromHexString creates a new script from a hex encoded string.
func NewFromHexString(s string) (*Script, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return NewFromBytes(b), nil
}

// NewFromBytes wraps a byte slice with the Script type.
func NewFromBytes(b []byte) *Script {
	s := Script(b)
	return &s
}

// NewFromASM creates a new script from a BitCoin ASM formatted string.
func NewFromASM(str string) (*Script, error) {
	s := Script{}

	if len(str) == 0 {
		return &s, nil
	}
	for _, section := range strings.Split(str, " ") {
		if val, ok := OpCodeStrings[section]; ok {
			_ = s.AppendOpcodes(val)
		} else {
			if err := s.AppendPushDataHexString(section); err != nil {
				return nil, ErrInvalidOpCode
			}
		}
	}

	return &s, nil
}

// AppendPushData takes data bytes and appends them to the script
// with proper PUSHDATA prefixes
func (s *Script) AppendPushData(d []byte) error {
	p, err := EncodeParts([][]byte{d})
	if err != nil {
		return err
	}

	*s = append(*s, p...)
	return nil
}

// AppendPushDataHexString takes a hex string and appends them to the
// script with proper PUSHDATA prefixes
func (s *Script) AppendPushDataHexString(str string) error {
	h, err := hex.DecodeString(str)
	if err != nil {
		return err
	}

	return s.AppendPushData(h)
}

// AppendPushDataString takes a string and appends its UTF-8 encoding
// to the script with proper PUSHDATA prefixes
func (s *Script) AppendPushDataString(str string) error {
	return s.AppendPushData([]byte(str))
}

// AppendPushDataArray takes an array of data bytes and appends them
// to the script with proper PUSHDATA prefixes
func (s *Script) AppendPushDataArray(d [][]byte) error {
	for _, part := range d {
		err := s.AppendPushData(part)
		if err != nil {
			return err
		}
	}
	// p, err := EncodeParts(d)
	// if err != nil {
	// 	return err
	// }

	// *s = append(*s, p...)
	return nil
}

// AppendPushDataStrings takes an array of strings and appends their
// UTF-8 encoding to the script with proper PUSHDATA prefixes
func (s *Script) AppendPushDataStrings(pushDataStrings []string) error {
	dataBytes := make([][]byte, 0)
	for _, str := range pushDataStrings {
		strBytes := []byte(str)
		dataBytes = append(dataBytes, strBytes)
	}
	return s.AppendPushDataArray(dataBytes)
}

// AppendOpcodes appends opcodes type to the script.
// This does not support appending OP_PUSHDATA opcodes, so use `Script.AppendPushData` instead.
func (s *Script) AppendOpcodes(oo ...uint8) error {
	for _, o := range oo {
		if OpDATA1 <= o && o <= OpPUSHDATA4 {
			return fmt.Errorf("%w: %s", ErrInvalidOpcodeType, OpCodeValues[o])
		}
	}
	*s = append(*s, oo...)
	return nil
}

func (s *Script) AppendBigInt(bigInt *big.Int) error {
	if bigInt.Cmp(big.NewInt(0)) == 0 {
		_ = s.AppendOpcodes(OpZERO)
		return nil
	}

	if bigInt.Cmp(big.NewInt(-1)) == 0 {
		_ = s.AppendOpcodes(Op1NEGATE)
		return nil
	}

	if bigInt.Cmp(big.NewInt(1)) >= 0 && bigInt.Cmp(big.NewInt(16)) <= 0 {
		_ = s.AppendOpcodes(uint8(bigInt.Int64()) + OpONE - 1)
		return nil
	}

	var num []byte
	if bigInt.Cmp(big.NewInt(0)) == -1 {
		num = bigInt.Neg(bigInt).Bytes()
		if (num[0] & 0x80) != 0 {
			num = append([]byte{0x80}, num...)
		} else {
			num[0] = num[0] | 0x80
		}
	} else {
		num = bigInt.Bytes()
		if (num[0] & 0x80) != 0 {
			num = append([]byte{0x00}, num...)
		}
	}

	if len(num) == 1 && num[0] == 0 {
		num = []byte{}
	}

	return s.AppendPushData(util.ReverseBytes(num))
}

func (s *Script) AppendInt(i int64) error {
	bigInt := big.NewInt(i)
	return s.AppendBigInt(bigInt)
}

// String implements the stringer interface and returns the hex string of script.
func (s *Script) String() string {
	return hex.EncodeToString(*s)
}

// ToASM returns the string ASM opcodes of the script.
func (s *Script) ToASM() (string, error) {
	if s == nil || len(*s) == 0 {
		return "", nil
	}

	var asm strings.Builder
	pos := 0
	for pos < len(*s) {
		op, err := s.ReadOp(&pos)
		if err != nil {
			return "", err
		}
		asm.WriteRune(' ')

		if len(op.Data) > 0 {
			asm.WriteString(hex.EncodeToString(op.Data))
		} else {
			asm.WriteString(OpCodeValues[op.OpCode])
		}
	}

	return asm.String()[1:], nil
}

// IsP2PKH returns true if this is a pay to pubkey hash output script.
func (s *Script) IsP2PKH() bool {
	b := []byte(*s)
	return len(b) == 25 &&
		b[0] == OpDUP &&
		b[1] == OpHASH160 &&
		b[2] == OpDATA20 &&
		b[23] == OpEQUALVERIFY &&
		b[24] == OpCHECKSIG
}

// IsP2PK returns true if this is a public key output script.
func (s *Script) IsP2PK() bool {
	parts, err := DecodeParts(*s)
	if err != nil {
		return false
	}

	if len(parts) == 2 && len(parts[0]) > 0 && parts[1][0] == OpCHECKSIG {
		pubkey := parts[0]
		version := pubkey[0]

		if (version == 0x04 || version == 0x06 || version == 0x07) && len(pubkey) == 65 {
			return true
		} else if (version == 0x03 || version == 0x02) && len(pubkey) == 33 {
			return true
		}
	}
	return false
}

// IsP2SH returns true if this is a p2sh output script.
// TODO: remove all p2sh stuff from repo
func (s *Script) IsP2SH() bool {
	b := []byte(*s)

	return len(b) == 23 &&
		b[0] == OpHASH160 &&
		b[1] == OpDATA20 &&
		b[22] == OpEQUAL
}

// IsData returns true if this is a data output script. This
// means the script starts with OP_RETURN or OP_FALSE OP_RETURN.
func (s *Script) IsData() bool {
	b := []byte(*s)

	return (len(b) > 0 && b[0] == OpRETURN) ||
		(len(b) > 1 && b[0] == OpFALSE && b[1] == OpRETURN)
}

// IsInscribed returns true if this script includes an
// inscription with any prepended script (not just p2pkh).
func (s *Script) IsInscribed() bool {
	isncPattern, _ := hex.DecodeString("0063036f7264")
	return bytes.Contains(*s, isncPattern)
}

// IsP2PKHInscription checks if it's a standard
// inscription with a P2PKH prefix script.
func (s *Script) IsP2PKHInscription() bool {
	p, err := DecodeParts(*s)
	if err != nil {
		return false
	}

	return isP2PKHInscriptionHelper(p)
}

// isP2PKHInscriptionHelper helper so that we don't need to call
// `DecodeParts()` multiple times, such as in `ParseInscription()`
func isP2PKHInscriptionHelper(parts [][]byte) bool {
	if len(parts) < 13 {
		return false
	}
	valid := parts[0][0] == OpDUP &&
		parts[1][0] == OpHASH160 &&
		parts[3][0] == OpEQUALVERIFY &&
		parts[4][0] == OpCHECKSIG &&
		parts[5][0] == OpFALSE &&
		parts[6][0] == OpIF &&
		parts[7][0] == 0x6f && parts[7][1] == 0x72 && parts[7][2] == 0x64 && // op_push "ord"
		parts[8][0] == OpTRUE &&
		parts[10][0] == OpFALSE &&
		parts[12][0] == OpENDIF

	if len(parts) > 13 {
		return parts[13][0] == OpRETURN && valid
	}
	return valid
}

// Slice a script to get back a subset of that script.
func (s *Script) Slice(start, end uint64) *Script {
	ss := *s
	sss := ss[start:end]
	return &sss
}

// IsMultiSigOut returns true if this is a multisig output script.
func (s *Script) IsMultiSigOut() bool {
	parts, err := DecodeParts(*s)
	if err != nil {
		return false
	}

	if len(parts) < 3 {
		return false
	}

	if !isSmallIntOp(parts[0][0]) {
		return false
	}

	for i := 1; i < len(parts)-2; i++ {
		if len(parts[i]) < 1 {
			return false
		}
	}

	return len(parts[len(parts)-2]) > 0 && isSmallIntOp(parts[len(parts)-2][0]) && len(parts[len(parts)-1]) > 0 &&
		parts[len(parts)-1][0] == OpCHECKMULTISIG
}

func isSmallIntOp(opcode byte) bool {
	return opcode == OpZERO || (opcode >= OpONE && opcode <= Op16)
}

// PublicKeyHash returns a public key hash byte array if the script is a P2PKH script.
func (s *Script) PublicKeyHash() ([]byte, error) {
	if s == nil || len(*s) == 0 {
		return nil, ErrEmptyScript
	}

	if (*s)[0] != OpDUP || len(*s) <= 2 || (*s)[1] != OpHASH160 {
		return nil, ErrNotP2PKH
	}

	parts, err := DecodeParts((*s)[2:])
	if err != nil {
		return nil, err
	}

	return parts[0], nil
}

// ScriptType returns the type of script this is as a string.
func (s *Script) ScriptType() string {
	if len(*s) == 0 {
		return ScriptTypeEmpty
	}
	if s.IsP2PKH() {
		return ScriptTypePubKeyHash
	}
	if s.IsP2PK() {
		return ScriptTypePubKey
	}
	if s.IsMultiSigOut() {
		return ScriptTypeMultiSig
	}
	if s.IsData() {
		return ScriptTypeNullData
	}
	if s.IsP2PKHInscription() {
		return ScriptTypePubKeyHashInscription
	}
	return ScriptTypeNonStandard
}

// Equals will compare the script to b and return true if they match.
func (s *Script) Equals(b *Script) bool {
	return bytes.Equal(*s, *b)
}

// EqualsBytes will compare the script to a byte representation of a
// script, b, and return true if they match.
func (s *Script) EqualsBytes(b []byte) bool {
	return bytes.Equal(*s, b)
}

// EqualsHex will compare the script to a hex string h,
// if they match then true is returned otherwise false.
func (s *Script) EqualsHex(h string) bool {
	return s.String() == h
}

// MinPushSize returns the minimum size of a push operation of the given data.
func MinPushSize(bb []byte) int {
	l := len(bb)

	// data length is larger than max supported by the bitcoin protocol
	if bits.UintSize == 64 && int64(l) > 0xffffffff {
		return 0
	}

	if l == 0 {
		return 1
	}

	if l == 1 {
		// data can be represented as Op1 to Op16, or OpNegate
		if bb[0] <= 16 || bb[0] == 0x81 {
			// OpX
			return 1
		}
		// OP_DATA_1 + data
		return 2
	}

	// OP_DATA_X + data
	if l <= 75 {
		return l + 1
	}
	// OP_PUSHDATA1 + length byte + data
	if l <= 0xff {
		return l + 2
	}
	// OP_PUSHDATA2 + two length bytes + data
	if l <= 0xffff {
		return l + 3
	}

	// OP_PUSHDATA4 + four length bytes + data
	return l + 5
}

// MarshalJSON convert script into json.
func (s *Script) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}

// UnmarshalJSON covert from json into *script.Script.
func (s *Script) UnmarshalJSON(bb []byte) error {
	ss, err := NewFromHexString(string(bytes.Trim(bb, `"`)))
	if err != nil {
		return err
	}

	*s = *ss
	return nil
}

package bscript

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestNumber(t *testing.T) {
	n := Number(128)
	x := n.Bytes()
	if !(x[0] == 128 && x[1] == 0) {
		t.Error("expect [128, 0]")
	}

	n = Number(127)
	x = n.Bytes()
	if !(x[0] == 127) {
		t.Error("expect [127]")
	}

	n = Number(-10)
	x = n.Bytes()
	if !(x[0] == 138) {
		t.Error("expect [138]")
	}
}

func hexToBytes(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic("invalid hex in source file: " + s)
	}
	return b
}

// TestScriptNumBytes ensures that converting from integral script numbers to
// byte representations works as expected.
func TestScriptNumBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		num        Number
		serialized []byte
	}{
		{0, nil},
		{1, hexToBytes("01")},
		{-1, hexToBytes("81")},
		{127, hexToBytes("7f")},
		{-127, hexToBytes("ff")},
		{128, hexToBytes("8000")},
		{-128, hexToBytes("8080")},
		{129, hexToBytes("8100")},
		{-129, hexToBytes("8180")},
		{256, hexToBytes("0001")},
		{-256, hexToBytes("0081")},
		{32767, hexToBytes("ff7f")},
		{-32767, hexToBytes("ffff")},
		{32768, hexToBytes("008000")},
		{-32768, hexToBytes("008080")},
		{65535, hexToBytes("ffff00")},
		{-65535, hexToBytes("ffff80")},
		{524288, hexToBytes("000008")},
		{-524288, hexToBytes("000088")},
		{7340032, hexToBytes("000070")},
		{-7340032, hexToBytes("0000f0")},
		{8388608, hexToBytes("00008000")},
		{-8388608, hexToBytes("00008080")},
		{2147483647, hexToBytes("ffffff7f")},
		{-2147483647, hexToBytes("ffffffff")},

		// Values that are out of range for data that is interpreted as
		// numbers, but are allowed as the result of numeric operations.
		{2147483648, hexToBytes("0000008000")},
		{-2147483648, hexToBytes("0000008080")},
		{2415919104, hexToBytes("0000009000")},
		{-2415919104, hexToBytes("0000009080")},
		{4294967295, hexToBytes("ffffffff00")},
		{-4294967295, hexToBytes("ffffffff80")},
		{4294967296, hexToBytes("0000000001")},
		{-4294967296, hexToBytes("0000000081")},
		{281474976710655, hexToBytes("ffffffffffff00")},
		{-281474976710655, hexToBytes("ffffffffffff80")},
		{72057594037927935, hexToBytes("ffffffffffffff00")},
		{-72057594037927935, hexToBytes("ffffffffffffff80")},
		{9223372036854775807, hexToBytes("ffffffffffffff7f")},
		{-9223372036854775807, hexToBytes("ffffffffffffffff")},
	}

	for _, test := range tests {
		gotBytes := test.num.Bytes()
		if !bytes.Equal(gotBytes, test.serialized) {
			t.Errorf("Bytes: did not get expected bytes for %d - "+
				"got %x, want %x", test.num, gotBytes,
				test.serialized)
			continue
		}
	}
}

func TestNumberint32(t *testing.T) {
	tests := []struct {
		in   Number
		want int32
	}{
		// Values inside the valid int32 range are just the values
		// themselves cscript to an int32.
		{0, 0},
		{1, 1},
		{-1, -1},
		{127, 127},
		{-127, -127},
		{128, 128},
		{-128, -128},
		{129, 129},
		{-129, -129},
		{256, 256},
		{-256, -256},
		{32767, 32767},
		{-32767, -32767},
		{32768, 32768},
		{-32768, -32768},
		{65535, 65535},
		{-65535, -65535},
		{524288, 524288},
		{-524288, -524288},
		{7340032, 7340032},
		{-7340032, -7340032},
		{8388608, 8388608},
		{-8388608, -8388608},
		{2147483647, 2147483647},
		{-2147483647, -2147483647},
		{-2147483648, -2147483648},

		// Values outside of the valid int32 range are limited to int32.
		{2147483648, 2147483647},
		{-2147483649, -2147483648},
		{1152921504606846975, 2147483647},
		{-1152921504606846975, -2147483648},
		{2305843009213693951, 2147483647},
		{-2305843009213693951, -2147483648},
		{4611686018427387903, 2147483647},
		{-4611686018427387903, -2147483648},
		{9223372036854775807, 2147483647},
		{-9223372036854775808, -2147483648},
	}

	for _, test := range tests {
		got := test.in.Int32()
		if int32(got) != test.want {
			t.Errorf("Int32: did not get expected value for %d - "+
				"got %d, want %d", test.in, got, test.want)
			continue
		}
	}
}
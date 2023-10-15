package int128

import (
	"fmt"
	"testing"
)

func TestInt128Format(t *testing.T) {
	tests := []struct {
		format string
		value  Int128
		want   string
	}{
		{"%d", Int128{0, 0}, "0"},
		{"%d", Int128{0, 1}.Neg(), "-1"},
		{"%+d", Int128{0, 0}, "+0"},
		{"%+d", Int128{0, 1}.Neg(), "-1"},
		{"% d", Int128{0, 0}, " 0"},
		{"% d", Int128{0, 1}.Neg(), "-1"},
		{"%8d", Int128{0, 0}, "       0"},
		{"%08d", Int128{0, 0}, "00000000"},
		{"%-8d", Int128{0, 0}, "0       "},
		{"%+08d", Int128{0, 0}, "+0000000"},
		{"%+8d", Int128{0, 0}, "      +0"},
		{"% 08d", Int128{0, 0}, " 0000000"},

		{"%b", Int128{0, 0xAA}, "10101010"},
		{"%#b", Int128{0, 0xAA}, "0b10101010"},

		{"%o", Int128{0, 0}, "0"},
		{"%o", Int128{0, 1}, "1"},
		{"%#o", Int128{0, 0}, "0"},
		{"%#o", Int128{0, 1}, "01"},

		{"%O", Int128{0, 0}, "0o0"},
		{"%O", Int128{0, 1}, "0o1"},
		{"%#O", Int128{0, 0}, "0o0"},
		{"%#O", Int128{0, 1}, "0o01"},

		{"%x", Int128{0, 0xabcd}, "abcd"},
		{"%#x", Int128{0, 0xabcd}, "0xabcd"},

		{"%X", Int128{0, 0xabcd}, "ABCD"},
		{"%#X", Int128{0, 0xabcd}, "0XABCD"},
	}

	for _, tt := range tests {
		got := fmt.Sprintf(tt.format, tt.value)
		if got != tt.want {
			t.Errorf("%#v: want %q, got %q", tt, tt.want, got)
		}
	}
}

func TestUint128Format(t *testing.T) {
	tests := []struct {
		format string
		value  Uint128
		want   string
	}{
		{"%d", Uint128{0, 0}, "0"},
		{"%+d", Uint128{0, 0}, "+0"},
		{"% d", Uint128{0, 0}, " 0"},
		{"%8d", Uint128{0, 0}, "       0"},
		{"%08d", Uint128{0, 0}, "00000000"},
		{"%-8d", Uint128{0, 0}, "0       "},
		{"%+08d", Uint128{0, 0}, "+0000000"},
		{"%+8d", Uint128{0, 0}, "      +0"},
		{"% 08d", Uint128{0, 0}, " 0000000"},

		{"%b", Uint128{0, 0xAA}, "10101010"},
		{"%#b", Uint128{0, 0xAA}, "0b10101010"},

		{"%o", Uint128{0, 0}, "0"},
		{"%o", Uint128{0, 1}, "1"},
		{"%#o", Uint128{0, 0}, "0"},
		{"%#o", Uint128{0, 1}, "01"},

		{"%O", Uint128{0, 0}, "0o0"},
		{"%O", Uint128{0, 1}, "0o1"},
		{"%#O", Uint128{0, 0}, "0o0"},
		{"%#O", Uint128{0, 1}, "0o01"},

		{"%x", Uint128{0, 0xabcd}, "abcd"},
		{"%#x", Uint128{0, 0xabcd}, "0xabcd"},

		{"%X", Uint128{0, 0xabcd}, "ABCD"},
		{"%#X", Uint128{0, 0xabcd}, "0XABCD"},
	}

	for _, tt := range tests {
		got := fmt.Sprintf(tt.format, tt.value)
		if got != tt.want {
			t.Errorf("%#v: want %q, got %q", tt, tt.want, got)
		}
	}
}

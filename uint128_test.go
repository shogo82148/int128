package int128

import (
	"math/big"
	"runtime"
	"strconv"
	"testing"
)

func TestUint128_Add(t *testing.T) {
	testCases := []struct {
		a, b, want Uint128
	}{
		{
			Uint128{0, 0},
			Uint128{0, 0},
			Uint128{0, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0, 1},
			Uint128{1, 0},
		},
		{
			Uint128{0, 1},
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{1, 0},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			Uint128{0, 1},
			Uint128{0, 0},
		},
		{
			Uint128{0, 1},
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Add(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v + %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func BenchmarkUint128_Add(b *testing.B) {
	x := Uint128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0}
	y := Uint128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0}
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(x.Add(y))
	}
}

func BenchmarkBigUint128_Add(b *testing.B) {
	x, _ := new(big.Int).SetString("0x1234_5678_9abc_def0_1234_5678_9abc_def0", 0)
	y, _ := new(big.Int).SetString("0x1234_5678_9abc_def0_1234_5678_9abc_def0", 0)
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(new(big.Int).Add(x, y))
	}
}

func TestUint128_Sub(t *testing.T) {
	testCases := []struct {
		a, b, want Uint128
	}{
		{
			Uint128{0, 0},
			Uint128{0, 0},
			Uint128{0, 0},
		},
		{
			Uint128{1, 0},
			Uint128{0, 1},
			Uint128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Uint128{0, 0},
			Uint128{0, 1},
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Sub(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v - %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func BenchmarkUint128_Sub(b *testing.B) {
	x := Uint128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0}
	y := Uint128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0}
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(x.Sub(y))
	}
}

func BenchmarkBigUint128_Sub(b *testing.B) {
	x, _ := new(big.Int).SetString("0x1234_5678_9abc_def0_1234_5678_9abc_def0", 0)
	y, _ := new(big.Int).SetString("0x1234_5678_9abc_def0_1234_5678_9abc_def0", 0)
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(new(big.Int).Sub(x, y))
	}
}

func TestUint128_Mul(t *testing.T) {
	testCases := []struct {
		a, b, want Uint128
	}{
		{
			Uint128{0, 0},
			Uint128{0, 0},
			Uint128{0, 0},
		},
		{
			Uint128{0, 1},
			Uint128{0, 1},
			Uint128{0, 1},
		},
		{
			Uint128{0, 1},
			Uint128{1, 0},
			Uint128{1, 0},
		},
		{
			Uint128{1, 0},
			Uint128{0, 1},
			Uint128{1, 0},
		},
		{
			Uint128{1, 0},
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0xffff_ffff_ffff_ffff, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{1, 0},
			Uint128{0xffff_ffff_ffff_ffff, 0},
		},
		{
			Uint128{1, 0},
			Uint128{1, 0},
			Uint128{0, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Mul(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v * %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint128_Div(t *testing.T) {
	testCases := []struct {
		a, b, want Uint128
	}{
		{
			Uint128{0x101, 0x100},
			Uint128{0, 0x100},
			Uint128{0x1, 0x100_0000_0000_0001},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			Uint128{1, 0},
			Uint128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{1, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0x7fff_ffff_ffff_ffff},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Div(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v / %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint128_Cmp(t *testing.T) {
	testCases := []struct {
		a, b Uint128
		want int
	}{
		{
			Uint128{0, 0},
			Uint128{0, 0},
			0,
		},
		{
			Uint128{0, 1},
			Uint128{0, 0},
			1,
		},
		{
			Uint128{0, 0},
			Uint128{0, 1},
			-1,
		},
		{
			Uint128{1, 0},
			Uint128{0, 0},
			1,
		},
		{
			Uint128{0, 0},
			Uint128{1, 0},
			-1,
		},
	}

	for i, tc := range testCases {
		got := tc.a.Cmp(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v * %#v should %d, but %d", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint128_And(t *testing.T) {
	testCases := []struct {
		a, b, want Uint128
	}{
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0, 0},
			Uint128{0, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0, 0},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			Uint128{0xffff_ffff_ffff_ffff, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0xffff_ffff_ffff_ffff},
		},
	}

	for i, tc := range testCases {
		got := tc.a.And(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v & %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint128_Or(t *testing.T) {
	testCases := []struct {
		a, b, want Uint128
	}{
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0, 0},
			Uint128{0xffff_ffff_ffff_ffff, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0xffff_ffff_ffff_ffff},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Or(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v | %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint128_Xor(t *testing.T) {
	testCases := []struct {
		a, b, want Uint128
	}{
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0, 0},
			Uint128{0xffff_ffff_ffff_ffff, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Xor(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v ^ %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint128_AndNot(t *testing.T) {
	testCases := []struct {
		a, b, want Uint128
	}{
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0, 0},
			Uint128{0xffff_ffff_ffff_ffff, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.AndNot(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v &^ %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint128_Not(t *testing.T) {
	testCases := []struct {
		a, want Uint128
	}{
		{
			Uint128{0, 0},
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			Uint128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			Uint128{0xffff_ffff_ffff_ffff, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Not()
		if got != tc.want {
			t.Errorf("%d: ^%#v should %#v, but %#v", i, tc.a, tc.want, got)
		}
	}
}

func TestUint128_Neg(t *testing.T) {
	testCases := []struct {
		a, want Uint128
	}{
		{
			Uint128{0, 0},
			Uint128{0, 0},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			Uint128{0, 1},
		},
		{
			Uint128{0, 1},
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Neg()
		if got != tc.want {
			t.Errorf("%d: -%#v should %#v, but %#v", i, tc.a, tc.want, got)
		}
	}
}

func TestUint128_Lsh(t *testing.T) {
	testCases := []struct {
		a    Uint128
		n    uint
		want Uint128
	}{
		{
			Uint128{0, 0},
			0,
			Uint128{0, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			1,
			Uint128{0x01, 0xffff_ffff_ffff_fffe},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			16,
			Uint128{0xffff, 0xffff_ffff_ffff_0000},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			32,
			Uint128{0xffff_ffff, 0xffff_ffff_0000_0000},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			64,
			Uint128{0xffff_ffff_ffff_ffff, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			65,
			Uint128{0xffff_ffff_ffff_fffe, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			128,
			Uint128{0, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Lsh(tc.n)
		if got != tc.want {
			t.Errorf("%d: %#v << %d should %#v, but %#v", i, tc.a, tc.n, tc.want, got)
		}
	}
}

func TestUint128_Rsh(t *testing.T) {
	testCases := []struct {
		a    Uint128
		n    uint
		want Uint128
	}{
		{
			Uint128{0, 0},
			0,
			Uint128{0, 0},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			1,
			Uint128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			16,
			Uint128{0xffff_ffff_ffff, 0xffff_0000_0000_0000},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			32,
			Uint128{0xffff_ffff, 0xffff_ffff_0000_0000},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			64,
			Uint128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			65,
			Uint128{0, 0x7fff_ffff_ffff_ffff},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			128,
			Uint128{0, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Rsh(tc.n)
		if got != tc.want {
			t.Errorf("%d: %#v >> %d should %#v, but %#v", i, tc.a, tc.n, tc.want, got)
		}
	}
}

func TestUint128_LeadingZeros(t *testing.T) {
	testCases := []struct {
		a    Uint128
		want int
	}{
		{
			Uint128{0, 0},
			128,
		},
		{
			Uint128{0, 0xffff_ffff},
			96,
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			64,
		},
		{
			Uint128{0xffff_ffff, 0},
			32,
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			0,
		},
	}

	for i, tc := range testCases {
		got := tc.a.LeadingZeros()
		if got != tc.want {
			t.Errorf("%d: LeadingZeros of %#v should %#v, but %#v", i, tc.a, tc.want, got)
		}
	}
}

func TestUint128_TrailingZeros(t *testing.T) {
	testCases := []struct {
		a    Uint128
		want int
	}{
		{
			Uint128{0, 0},
			128,
		},
		{
			Uint128{0xffff_ffff_0000_0000, 0},
			96,
		},
		{
			Uint128{1, 0},
			64,
		},
		{
			Uint128{0, 0xffff_ffff_0000_0000},
			32,
		},
		{
			Uint128{0, 1},
			0,
		},
	}

	for i, tc := range testCases {
		got := tc.a.TrailingZeros()
		if got != tc.want {
			t.Errorf("%d: TrailingZeros %#v should %#v, but %#v", i, tc.a, tc.want, got)
		}
	}
}

func TestUint128_Len(t *testing.T) {
	testCases := []struct {
		a    Uint128
		want int
	}{
		{
			Uint128{0, 0},
			0,
		},
		{
			Uint128{0, 0xffff_ffff},
			32,
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			64,
		},
		{
			Uint128{0xffff_ffff, 0},
			96,
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			128,
		},
	}

	for i, tc := range testCases {
		got := tc.a.Len()
		if got != tc.want {
			t.Errorf("%d: Len of %#v should %#v, but %#v", i, tc.a, tc.want, got)
		}
	}
}

func TestUint128_OnesCount(t *testing.T) {
	testCases := []struct {
		a    Uint128
		want int
	}{
		{
			Uint128{0, 0},
			0,
		},
		{
			Uint128{0, 0xffff_ffff},
			32,
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			64,
		},
		{
			Uint128{0xffff_ffff, 0},
			32,
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			64,
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			128,
		},
	}

	for i, tc := range testCases {
		got := tc.a.OnesCount()
		if got != tc.want {
			t.Errorf("%d: OnesCount of %#v should %#v, but %#v", i, tc.a, tc.want, got)
		}
	}
}

func TestUint128_RotateLeft(t *testing.T) {
	testCases := []struct {
		a    Uint128
		n    int
		want Uint128
	}{
		{
			Uint128{0, 0},
			0,
			Uint128{0, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			1,
			Uint128{1, 0xffff_ffff_ffff_fffe},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			-1,
			Uint128{0x8000_0000_0000_0000, 0x7fff_ffff_ffff_ffff},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			1,
			Uint128{0xffff_ffff_ffff_fffe, 1},
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0},
			-1,
			Uint128{0x7fff_ffff_ffff_ffff, 0x8000_0000_0000_0000},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			32,
			Uint128{0xffff_ffff, 0xffff_ffff_0000_0000},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			-32,
			Uint128{0xffff_ffff_0000_0000, 0xffff_ffff},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			64,
			Uint128{0xffff_ffff_ffff_ffff, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			-64,
			Uint128{0xffff_ffff_ffff_ffff, 0},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			96,
			Uint128{0xffff_ffff_0000_0000, 0xffff_ffff},
		},
		{
			Uint128{0, 0xffff_ffff_ffff_ffff},
			-96,
			Uint128{0xffff_ffff, 0xffff_ffff_0000_0000},
		},
	}

	for i, tc := range testCases {
		got := tc.a.RotateLeft(tc.n)
		if got != tc.want {
			t.Errorf("%d: %#v.Rotate(%d) should %#v, but %#v", i, tc.a, tc.n, tc.want, got)
		}
	}
}

func TestUint128_Reverse(t *testing.T) {
	testCases := []struct {
		a, want Uint128
	}{
		{
			Uint128{0, 0},
			Uint128{0, 0},
		},
		{
			Uint128{0, 1},
			Uint128{0x8000_0000_0000_0000, 0},
		},
		{
			Uint128{0x8000_0000_0000_0000, 0},
			Uint128{0, 1},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Reverse()
		if got != tc.want {
			t.Errorf("%d: %#v.Reverse() should %#v, but %#v", i, tc.a, tc.want, got)
		}
	}
}

func TestUint128_ReverseBytes(t *testing.T) {
	testCases := []struct {
		a, want Uint128
	}{
		{
			Uint128{0, 0},
			Uint128{0, 0},
		},
		{
			Uint128{0, 0x1234_5678_9abc_def0},
			Uint128{0xf0de_bc9a_7856_3412, 0},
		},
		{
			Uint128{0xf0de_bc9a_7856_3412, 0},
			Uint128{0, 0x1234_5678_9abc_def0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.ReverseBytes()
		if got != tc.want {
			t.Errorf("%d: %#v.Reverse() should %#v, but %#v", i, tc.a, tc.want, got)
		}
	}
}

func TestUint128_String(t *testing.T) {
	testCases := []struct {
		a    Uint128
		want string
	}{
		{
			Uint128{0, 0},
			"0",
		},
		{
			// the max value of small integers
			Uint128{0, 99},
			"99",
		},
		{
			Uint128{0, 100},
			"100",
		},
		{
			// the max value of uint64
			Uint128{0, 0xffff_ffff_ffff_ffff},
			"18446744073709551615",
		},
		{
			// the max value of Uint128
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			"340282366920938463463374607431768211455",
		},
	}

	for i, tc := range testCases {
		got := tc.a.String()
		if got != tc.want {
			t.Errorf("%d: string of %#v should %q, but %q", i, tc.a, tc.want, got)
		}
	}
}

func BenchmarkUint128_String(b *testing.B) {
	b.Run("the max value of small integers", func(b *testing.B) {
		v := Uint128{0, 99}
		for i := 0; i < b.N; i++ {
			runtime.KeepAlive(v.String())
		}
	})

	b.Run("strconv.FormatUint(99, 10)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			runtime.KeepAlive(strconv.FormatUint(99, 10))
		}
	})

	b.Run("the max value of uint64", func(b *testing.B) {
		v := Uint128{0, 0xffff_ffff_ffff_ffff}
		for i := 0; i < b.N; i++ {
			runtime.KeepAlive(v.String())
		}
	})

	b.Run("strconv.FormatUint(0xffff_ffff_ffff_ffff, 10)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			runtime.KeepAlive(strconv.FormatUint(0xffff_ffff_ffff_ffff, 10))
		}
	})

	b.Run("the max value of Uint128", func(b *testing.B) {
		v := Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff}
		for i := 0; i < b.N; i++ {
			runtime.KeepAlive(v.String())
		}
	})
}

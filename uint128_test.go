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

func TestFloat64ToUint128(t *testing.T) {
	testCases := []struct {
		input float64
		want  Uint128
	}{
		{
			0,
			Uint128{0, 0},
		},
		{
			1,
			Uint128{0, 1},
		},
		{
			-1,
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
		},
		{
			// the maximum float64 value that that can correctly represent an integer
			1 << 53,
			Uint128{0, 0x20000000000000},
		},
		{
			// the maximum float64 value that can convert to uint64
			(1<<53 - 1) << 11,
			Uint128{0, 0xfffffffffffff800},
		},
		{
			// the maximum float64 value that can convert to Uint128
			(1<<53 - 1) << 75,
			Uint128{0xfffffffffffff800, 0},
		},
	}

	for i, tc := range testCases {
		got := Float64ToUint128(tc.input)
		if got != tc.want {
			t.Errorf("%d: Float64ToUint128(%f) should %#v, but %#v", i, tc.input, tc.want, got)
		}
	}
}

func TestUint128_Text(t *testing.T) {
	testCases := []struct {
		a    Uint128
		base int
		want string
	}{
		{
			Uint128{0, 0},
			2,
			"0",
		},
		{
			Uint128{0, 35},
			36,
			"z",
		},
		{
			Uint128{0, 99},
			10,
			"99",
		},
		{
			Uint128{0, 0x1234_5678_9abc_def0},
			16,
			"123456789abcdef0",
		},
		{
			Uint128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0},
			16,
			"123456789abcdef0123456789abcdef0",
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			2,
			"11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			3,
			"202201102121002021012000211012011021221022212021111001022110211020010021100121010",
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			4,
			"3333333333333333333333333333333333333333333333333333333333333333",
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			5,
			"11031110441201303134210404233413032443021130230130231310",
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			6,
			"23053353530155550541354043543542243325553444410303",
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			7,
			"3115512162124626343001006330151620356026315303",
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			8,
			"3777777777777777777777777777777777777777777",
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			9,
			"22642532235024164257285244038424203240533",
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			10,
			"340282366920938463463374607431768211455",
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			16,
			"ffffffffffffffffffffffffffffffff",
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			32,
			"7vvvvvvvvvvvvvvvvvvvvvvvvv",
		},
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			36,
			"f5lxx1zz5pnorynqglhzmsp33",
		},
	}

	buf := make([]byte, 0, 128)
	for i, tc := range testCases {
		got := tc.a.Text(tc.base)
		if got != tc.want {
			t.Errorf("%d: %#v.Text(%d) should %q, but %q", i, tc.a, tc.base, tc.want, got)
		}

		buf = tc.a.Append(buf[:0], tc.base)
		if string(buf) != tc.want {
			t.Errorf("%d: %#v.Append(buf, %d) should %q, but %q", i, tc.a, tc.base, tc.want, got)
		}
	}
}

func BenchmarkUint128_Append(b *testing.B) {
	buf := make([]byte, 0, 128)
	b.Run("the max value of small integers(base 10)", func(b *testing.B) {
		v := Uint128{0, 99}
		for i := 0; i < b.N; i++ {
			v.Append(buf, 10)
		}
	})

	b.Run("strconv.FormatUint(99, 10)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			runtime.KeepAlive(strconv.FormatUint(99, 10))
		}
	})

	b.Run("the max value of uint64 (base 2)", func(b *testing.B) {
		v := Uint128{0, 0xffff_ffff_ffff_ffff}
		for i := 0; i < b.N; i++ {
			v.Append(buf, 2)
		}
	})

	b.Run("strconv.FormatUint(0xffff_ffff_ffff_ffff, 2)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strconv.AppendUint(buf, 0xffff_ffff_ffff_ffff, 2)
		}
	})

	b.Run("the max value of uint64 (base 3)", func(b *testing.B) {
		v := Uint128{0, 0xffff_ffff_ffff_ffff}
		for i := 0; i < b.N; i++ {
			v.Append(buf, 3)
		}
	})

	b.Run("strconv.FormatUint(0xffff_ffff_ffff_ffff, 3)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strconv.AppendUint(buf, 0xffff_ffff_ffff_ffff, 3)
		}
	})

	b.Run("the max value of uint64 (base 16)", func(b *testing.B) {
		v := Uint128{0, 0xffff_ffff_ffff_ffff}
		for i := 0; i < b.N; i++ {
			v.Append(buf, 16)
		}
	})

	b.Run("strconv.FormatUint(0xffff_ffff_ffff_ffff, 16)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strconv.AppendUint(buf, 0xffff_ffff_ffff_ffff, 16)
		}
	})

	b.Run("the max value of uint64 (base 36)", func(b *testing.B) {
		v := Uint128{0, 0xffff_ffff_ffff_ffff}
		for i := 0; i < b.N; i++ {
			v.Append(buf, 36)
		}
	})

	b.Run("strconv.FormatUint(0xffff_ffff_ffff_ffff, 36)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strconv.AppendUint(buf, 0xffff_ffff_ffff_ffff, 36)
		}
	})

	b.Run("the max value of Uint128 (base 2)", func(b *testing.B) {
		v := Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff}
		for i := 0; i < b.N; i++ {
			v.Append(buf, 2)
		}
	})

	b.Run("the max value of Uint128 (base 3)", func(b *testing.B) {
		v := Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff}
		for i := 0; i < b.N; i++ {
			v.Append(buf, 3)
		}
	})

	b.Run("the max value of Uint128 (base 8)", func(b *testing.B) {
		v := Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff}
		for i := 0; i < b.N; i++ {
			v.Append(buf, 8)
		}
	})

	b.Run("the max value of Uint128 (base 16)", func(b *testing.B) {
		v := Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff}
		for i := 0; i < b.N; i++ {
			v.Append(buf, 16)
		}
	})

	b.Run("the max value of Uint128 (base 32)", func(b *testing.B) {
		v := Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff}
		for i := 0; i < b.N; i++ {
			v.Append(buf, 32)
		}
	})

	b.Run("the max value of Uint128 (base 36)", func(b *testing.B) {
		v := Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff}
		for i := 0; i < b.N; i++ {
			v.Append(buf, 36)
		}
	})

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

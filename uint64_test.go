package int128

import (
	"math/big"
	"runtime"
	"testing"
)

func TestUint64_Add(t *testing.T) {
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
			t.Errorf("%d: %v + %v should %v, but %v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func BenchmarkUint64_Add(b *testing.B) {
	x := Uint128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0}
	y := Uint128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0}
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(x.Add(y))
	}
}

func BenchmarkBigUint64_Add(b *testing.B) {
	x, _ := new(big.Int).SetString("0x1234_5678_9abc_def0_1234_5678_9abc_def0", 0)
	y, _ := new(big.Int).SetString("0x1234_5678_9abc_def0_1234_5678_9abc_def0", 0)
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(new(big.Int).Add(x, y))
	}
}

func TestUint64_Sub(t *testing.T) {
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
			t.Errorf("%d: %v - %v should %v, but %v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func BenchmarkUint64_Sub(b *testing.B) {
	x := Uint128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0}
	y := Uint128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0}
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(x.Sub(y))
	}
}

func BenchmarkBigUint64_Sub(b *testing.B) {
	x, _ := new(big.Int).SetString("0x1234_5678_9abc_def0_1234_5678_9abc_def0", 0)
	y, _ := new(big.Int).SetString("0x1234_5678_9abc_def0_1234_5678_9abc_def0", 0)
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(new(big.Int).Sub(x, y))
	}
}

func TestUint64_Mul(t *testing.T) {
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
			t.Errorf("%d: %v * %v should %v, but %v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint64_Div(t *testing.T) {
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
			t.Errorf("%d: %v / %v should %v, but %v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint64_Cmp(t *testing.T) {
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
			t.Errorf("%d: %v * %v should %d, but %d", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint64_And(t *testing.T) {
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
			t.Errorf("%d: %v & %v should %v, but %v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint64_Or(t *testing.T) {
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
			t.Errorf("%d: %v | %v should %v, but %v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint64_Xor(t *testing.T) {
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
			t.Errorf("%d: %v ^ %v should %v, but %v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint64_AndNot(t *testing.T) {
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
			t.Errorf("%d: %v &^ %v should %v, but %v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestUint64_Not(t *testing.T) {
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
			t.Errorf("%d: ^%v should %v, but %v", i, tc.a, tc.want, got)
		}
	}
}

func TestUint64_Neg(t *testing.T) {
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
			t.Errorf("%d: -%v should %v, but %v", i, tc.a, tc.want, got)
		}
	}
}

func TestUint64_Lsh(t *testing.T) {
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
			t.Errorf("%d: %v << %d should %v, but %v", i, tc.a, tc.n, tc.want, got)
		}
	}
}

func TestUint64_Rsh(t *testing.T) {
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
			t.Errorf("%d: %v >> %d should %v, but %v", i, tc.a, tc.n, tc.want, got)
		}
	}
}

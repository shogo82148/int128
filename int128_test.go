package int128

import (
	"runtime"
	"testing"
)

func TestInt128_Add(t *testing.T) {
	testCases := []struct {
		a, b, want Int128
	}{
		{
			Int128{0, 0},
			Int128{0, 0},
			Int128{0, 0},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{0, 1},
			Int128{1, 0},
		},
		{
			Int128{0, 1},
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{1, 0},
		},
		{
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{0, 1},
			Int128{0, 0},
		},
		{
			Int128{0, 1},
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{0, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Add(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v + %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func BenchmarkInt128_Add(b *testing.B) {
	x := Int128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0}
	y := Int128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0}
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(x.Add(y))
	}
}

func TestInt128_Sub(t *testing.T) {
	testCases := []struct {
		a, b, want Int128
	}{
		{
			Int128{0, 0},
			Int128{0, 0},
			Int128{0, 0},
		},
		{
			Int128{1, 0},
			Int128{0, 1},
			Int128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{0, 0},
			Int128{0, 1},
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Sub(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v - %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func BenchmarkInt128_Sub(b *testing.B) {
	x := Int128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0}
	y := Int128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0}
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(x.Sub(y))
	}
}

func TestInt128_Cmp(t *testing.T) {
	testCases := []struct {
		a, b Int128
		want int
	}{
		{
			Int128{0, 0},
			Int128{0, 0},
			0,
		},
		{
			Int128{0, 1},
			Int128{0, 0},
			1,
		},
		{
			Int128{0, 0},
			Int128{0, 1},
			-1,
		},
		{
			Int128{0, 0},
			Int128{-1, 0xffff_ffff_ffff_ffff},
			1,
		},
		{
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{0, 0},
			-1,
		},
		{
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{-1, 0xffff_ffff_ffff_fffe},
			1,
		},
		{
			Int128{-1, 0xffff_ffff_ffff_fffe},
			Int128{-1, 0xffff_ffff_ffff_ffff},
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

func TestInt128_And(t *testing.T) {
	testCases := []struct {
		a, b, want Int128
	}{
		{
			Int128{-1, 0},
			Int128{0, 0},
			Int128{0, 0},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{-1, 0},
			Int128{0, 0},
		},
		{
			Int128{-1, 0},
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{-1, 0},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{0, 0xffff_ffff_ffff_ffff},
		},
	}

	for i, tc := range testCases {
		got := tc.a.And(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v & %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestInt128_Or(t *testing.T) {
	testCases := []struct {
		a, b, want Int128
	}{
		{
			Int128{-1, 0},
			Int128{0, 0},
			Int128{-1, 0},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{-1, 0},
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{-1, 0},
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{0, 0xffff_ffff_ffff_ffff},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Or(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v | %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestInt128_Xor(t *testing.T) {
	testCases := []struct {
		a, b, want Int128
	}{
		{
			Int128{-1, 0},
			Int128{0, 0},
			Int128{-1, 0},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{-1, 0},
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{-1, 0},
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{0, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Xor(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v ^ %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestInt128_AndNot(t *testing.T) {
	testCases := []struct {
		a, b, want Int128
	}{
		{
			Int128{-1, 0},
			Int128{0, 0},
			Int128{-1, 0},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{-1, 0},
			Int128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{-1, 0},
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{0, 0},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{0, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.AndNot(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v &^ %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
	}
}

func TestInt128_Not(t *testing.T) {
	testCases := []struct {
		a, want Int128
	}{
		{
			Int128{0, 0},
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{0, 0},
		},
		{
			Int128{-1, 0},
			Int128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			Int128{-1, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Not()
		if got != tc.want {
			t.Errorf("%d: ^%#v should %#v, but %#v", i, tc.a, tc.want, got)
		}
	}
}

func TestInt128_Neg(t *testing.T) {
	testCases := []struct {
		a, want Int128
	}{
		{
			Int128{0, 0},
			Int128{0, 0},
		},
		{
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{0, 1},
		},
		{
			Int128{0, 1},
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Neg()
		if got != tc.want {
			t.Errorf("%d: -%#v should %#v, but %#v", i, tc.a, tc.want, got)
		}
	}
}

func TestInt128_Lsh(t *testing.T) {
	testCases := []struct {
		a    Int128
		n    uint
		want Int128
	}{
		{
			Int128{0, 0},
			0,
			Int128{0, 0},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			1,
			Int128{0x01, 0xffff_ffff_ffff_fffe},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			16,
			Int128{0xffff, 0xffff_ffff_ffff_0000},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			32,
			Int128{0xffff_ffff, 0xffff_ffff_0000_0000},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			64,
			Int128{-1, 0},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			65,
			Int128{-2, 0},
		},
		{
			Int128{0, 0xffff_ffff_ffff_ffff},
			128,
			Int128{0, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Lsh(tc.n)
		if got != tc.want {
			t.Errorf("%d: %#v << %d should %#v, but %#v", i, tc.a, tc.n, tc.want, got)
		}
	}
}

func TestInt128_Rsh(t *testing.T) {
	testCases := []struct {
		a    Int128
		n    uint
		want Int128
	}{
		{
			Int128{0, 0},
			0,
			Int128{0, 0},
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			1,
			Int128{0x3fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0},
			16,
			Int128{0x7fff_ffff_ffff, 0xffff_0000_0000_0000},
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0},
			32,
			Int128{0x7fff_ffff, 0xffff_ffff_0000_0000},
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0},
			64,
			Int128{0, 0x7fff_ffff_ffff_ffff},
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0},
			65,
			Int128{0, 0x3fff_ffff_ffff_ffff},
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0},
			128,
			Int128{0, 0},
		},

		// sign extension
		{
			Int128{-1, 0xffff_ffff_ffff_ffff},
			1,
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{-1, 0},
			16,
			Int128{-1, 0xffff_0000_0000_0000},
		},
		{
			Int128{-1, 0},
			32,
			Int128{-1, 0xffff_ffff_0000_0000},
		},
		{
			Int128{-1, 0},
			64,
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{-1, 0},
			65,
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{-1, 0},
			128,
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Rsh(tc.n)
		if got != tc.want {
			t.Errorf("%d: %#v >> %d should %#v, but %#v", i, tc.a, tc.n, tc.want, got)
		}
	}
}

func TestFloat64ToInt128(t *testing.T) {
	testCases := []struct {
		input float64
		want  Int128
	}{
		{
			0,
			Int128{0, 0},
		},
		{
			1,
			Int128{0, 1},
		},
		{
			-1,
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
		{
			// the maximum float64 value that that can correctly represent an integer
			1 << 53,
			Int128{0, 0x20000000000000},
		},
		{
			// the maximum float64 value that can convert to uint64
			(1<<53 - 1) << 11,
			Int128{0, 0xfffffffffffff800},
		},
		{
			// the maximum float64 value that can convert to Int128
			(1<<53 - 1) << 74,
			Int128{0x7ffffffffffffc00, 0},
		},
		{
			// the minimum float64 value that can convert to Int128
			-1 << 127,
			Int128{-0x8000_0000_0000_0000, 0},
		},
	}

	for i, tc := range testCases {
		got := Float64ToInt128(tc.input)
		if got != tc.want {
			t.Errorf("%d: Float64ToInt128(%f) should %#v, but %#v", i, tc.input, tc.want, got)
		}
	}
}

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

func TestInt128_Mul(t *testing.T) {
	testCases := []struct {
		a, b, want Int128
	}{
		{
			Int128{0, 0},
			Int128{0, 0},
			Int128{0, 0},
		},
		{
			Int128{0, 1},
			Int128{0, 1},
			Int128{0, 1},
		},
		{
			Int128{0, 1},
			Int128{1, 0},
			Int128{1, 0},
		},
		{
			Int128{1, 0},
			Int128{0, 1},
			Int128{1, 0},
		},
		{
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{0, 1},
		},
		{
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{0, 1},
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{1, 0},
			Int128{-1, 0},
		},
	}

	for i, tc := range testCases {
		got := tc.a.Mul(tc.b)
		if got != tc.want {
			t.Errorf("%d: %#v * %#v should %#v, but %#v", i, tc.a, tc.b, tc.want, got)
		}
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

func TestInt128_Text(t *testing.T) {
	testCases := []struct {
		a    Int128
		base int
		want string
	}{
		{
			Int128{0, 0},
			2,
			"0",
		},
		{
			Int128{0, 35},
			36,
			"z",
		},
		{
			Int128{0, 99},
			10,
			"99",
		},
		{
			Int128{0, 0x1234_5678_9abc_def0},
			16,
			"123456789abcdef0",
		},
		{
			Int128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0},
			16,
			"123456789abcdef0123456789abcdef0",
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			2,
			"1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			3,
			"101100201022001010121000102002120122110122221010202000122201220121120010200022001",
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			4,
			"1333333333333333333333333333333333333333333333333333333333333333",
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			5,
			"3013030220323124042102424341431241221233040112312340402",
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			6,
			"11324454543055553250455021551551121442554522203131",
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			7,
			"1406241064412313155000336513424310163013142501",
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			8,
			"1777777777777777777777777777777777777777777",
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			9,
			"11321261117012076573587122018656546120261",
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			10,
			"170141183460469231731687303715884105727",
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			16,
			"7fffffffffffffffffffffffffffffff",
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			32,
			"3vvvvvvvvvvvvvvvvvvvvvvvvv",
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			36,
			"7ksyyizzkutudzbv8aqztecjj",
		},
		{
			Int128{-0x8000_0000_0000_0000, 0},
			2,
			"-10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			Int128{-0x8000_0000_0000_0000, 0},
			3,
			"-101100201022001010121000102002120122110122221010202000122201220121120010200022002",
		},
		{
			Int128{-0x8000_0000_0000_0000, 0},
			4,
			"-2000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			Int128{-0x8000_0000_0000_0000, 0},
			5,
			"-3013030220323124042102424341431241221233040112312340403",
		},
		{
			Int128{-0x8000_0000_0000_0000, 0},
			6,
			"-11324454543055553250455021551551121442554522203132",
		},
		{
			Int128{-0x8000_0000_0000_0000, 0},
			7,
			"-1406241064412313155000336513424310163013142502",
		},
		{
			Int128{-0x8000_0000_0000_0000, 0},
			8,
			"-2000000000000000000000000000000000000000000",
		},
		{
			Int128{-0x8000_0000_0000_0000, 0},
			9,
			"-11321261117012076573587122018656546120262",
		},
		{
			Int128{-0x8000_0000_0000_0000, 0},
			10,
			"-170141183460469231731687303715884105728",
		},
		{
			Int128{-0x8000_0000_0000_0000, 0},
			16,
			"-80000000000000000000000000000000",
		},
		{
			Int128{-0x8000_0000_0000_0000, 0},
			32,
			"-40000000000000000000000000",
		},
		{
			Int128{-0x8000_0000_0000_0000, 0},
			36,
			"-7ksyyizzkutudzbv8aqztecjk",
		},
	}

	buf := make([]byte, 0, 128+1)
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

func TestInt128_String(t *testing.T) {
	testCases := []struct {
		a    Int128
		want string
	}{
		{
			Int128{0, 0},
			"0",
		},
		{
			// the max value of small integers
			Int128{0, 99},
			"99",
		},
		{
			Int128{0, 100},
			"100",
		},
		{
			// the max value of uint64
			Int128{0, 0xffff_ffff_ffff_ffff},
			"18446744073709551615",
		},
		{
			// the max value of Int128
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			"170141183460469231731687303715884105727",
		},
		{
			// the minium value of Int128
			Int128{-0x8000_0000_0000_0000, 0},
			"-170141183460469231731687303715884105728",
		},
	}

	for i, tc := range testCases {
		got := tc.a.String()
		if got != tc.want {
			t.Errorf("%d: string of %#v should %q, but %q", i, tc.a, tc.want, got)
		}
	}
}

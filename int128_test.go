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

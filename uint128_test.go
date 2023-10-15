package int128

import (
	"math"
	"math/big"
	"runtime"
	"strconv"
	"testing"
	"testing/quick"
)

// uint128Input is used for benchmarks to prevent compiler optimizations.
var uint128Input = Uint128{0, 42}

var bigModUint128, _ = new(big.Int).SetString("100000000000000000000000000000000", 16)

// uint128ToBig converts x to a big.Int.
func uint128ToBig(b *big.Int, x Uint128) *big.Int {
	if b == nil {
		b = new(big.Int)
	}
	b.SetBytes([]byte{
		byte(x.H >> 56),
		byte(x.H >> 48),
		byte(x.H >> 40),
		byte(x.H >> 32),
		byte(x.H >> 24),
		byte(x.H >> 16),
		byte(x.H >> 8),
		byte(x.H),
		byte(x.L >> 56),
		byte(x.L >> 48),
		byte(x.L >> 40),
		byte(x.L >> 32),
		byte(x.L >> 24),
		byte(x.L >> 16),
		byte(x.L >> 8),
		byte(x.L),
	})
	return b
}

// bigToUint128 converts x to a Uint128.
func bigToUint128(x *big.Int) Uint128 {
	var buf [16]byte
	z := new(big.Int).Mod(x, bigModUint128)
	z.FillBytes(buf[:])
	return Uint128{
		H: (uint64(buf[0]) << 56) | (uint64(buf[1]) << 48) | (uint64(buf[2]) << 40) | (uint64(buf[3]) << 32) |
			(uint64(buf[4]) << 24) | (uint64(buf[5]) << 16) | (uint64(buf[6]) << 8) | uint64(buf[7]),
		L: (uint64(buf[8]) << 56) | (uint64(buf[9]) << 48) | (uint64(buf[10]) << 40) | (uint64(buf[11]) << 32) |
			(uint64(buf[12]) << 24) | (uint64(buf[13]) << 16) | (uint64(buf[14]) << 8) | uint64(buf[15]),
	}
}

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
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Add(uint128Input))
	}
}

func BenchmarkBigUint128_Add(b *testing.B) {
	x, _ := new(big.Int).SetString("0x1234_5678_9abc_def0_1234_5678_9abc_def0", 0)
	y, _ := new(big.Int).SetString("0x1234_5678_9abc_def0_1234_5678_9abc_def0", 0)
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(new(big.Int).Add(x, y))
	}
}

func TestUint128_AddQuick(t *testing.T) {
	f := func(a, b Uint128) Uint128 {
		return a.Add(b)
	}
	g := func(a, b Uint128) Uint128 {
		bigA := uint128ToBig(new(big.Int), a)
		bigB := uint128ToBig(new(big.Int), b)
		bigA.Add(bigA, bigB)
		return bigToUint128(bigA)
	}
	if err := quick.CheckEqual(f, g, &quick.Config{
		MaxCountScale: 1000,
	}); err != nil {
		t.Error(err)
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
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Sub(uint128Input))
	}
}

func BenchmarkBigUint128_Sub(b *testing.B) {
	x, _ := new(big.Int).SetString("0x1234_5678_9abc_def0_1234_5678_9abc_def0", 0)
	y, _ := new(big.Int).SetString("0x1234_5678_9abc_def0_1234_5678_9abc_def0", 0)
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(new(big.Int).Sub(x, y))
	}
}

func TestUint128_SubQuick(t *testing.T) {
	f := func(a, b Uint128) Uint128 {
		return a.Sub(b)
	}
	g := func(a, b Uint128) Uint128 {
		bigA := uint128ToBig(new(big.Int), a)
		bigB := uint128ToBig(new(big.Int), b)
		bigA.Sub(bigA, bigB)
		return bigToUint128(bigA)
	}
	if err := quick.CheckEqual(f, g, &quick.Config{
		MaxCountScale: 1000,
	}); err != nil {
		t.Error(err)
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

func TestUint128_MulQuick(t *testing.T) {
	f := func(a, b Uint128) Uint128 {
		return a.Mul(b)
	}
	g := func(a, b Uint128) Uint128 {
		bigA := uint128ToBig(new(big.Int), a)
		bigB := uint128ToBig(new(big.Int), b)
		bigA.Mul(bigA, bigB)
		return bigToUint128(bigA)
	}
	if err := quick.CheckEqual(f, g, &quick.Config{
		MaxCountScale: 1000,
	}); err != nil {
		t.Error(err)
	}
}

func BenchmarkUint128_Mul(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Mul(uint128Input))
	}
}

func TestUint128_DivMod(t *testing.T) {
	testCases := []struct {
		a, b, div, mod Uint128
	}{
		{
			Uint128{0x101, 0x123},
			Uint128{0, 0x100},
			Uint128{0x1, 0x100_0000_0000_0001},
			Uint128{0, 0x23},
		},
		{
			Uint128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			Uint128{1, 0},
			Uint128{0, 0x7fff_ffff_ffff_ffff},
			Uint128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Uint128{0x7fff_ffff_ffff_ffff, 0},
			Uint128{1, 0xffff_ffff_ffff_ffff},
			Uint128{0, 0x3fff_ffff_ffff_ffff},
			Uint128{1, 0x3fff_ffff_ffff_ffff},
		},
	}

	for i, tc := range testCases {
		div, mod := tc.a.DivMod(tc.b)
		if div != tc.div {
			t.Errorf("%d: %#v / %#v should %#v, but %#v", i, tc.a, tc.b, tc.div, div)
		}
		if mod != tc.mod {
			t.Errorf("%d: %#v %% %#v should %#v, but %#v", i, tc.a, tc.b, tc.mod, mod)
		}

		div = tc.a.Div(tc.b)
		if div != tc.div {
			t.Errorf("%d: %#v / %#v should %#v, but %#v", i, tc.a, tc.b, tc.div, div)
		}
		mod = tc.a.Mod(tc.b)
		if mod != tc.mod {
			t.Errorf("%d: %#v %% %#v should %#v, but %#v", i, tc.a, tc.b, tc.mod, mod)
		}
	}
}

func TestUint128_DivModQuick(t *testing.T) {
	f := func(l uint, a, b Uint128) (Uint128, Uint128) {
		a.H |= 1 << 63
		b.H |= 1 << 63
		a = a.Rsh(l % 128)
		b = b.Rsh(l % 128)
		if b == (Uint128{0, 0}) {
			return Uint128{0, 0}, Uint128{0, 0}
		}
		return a.DivMod(b)
	}
	g := func(l uint, a, b Uint128) (Uint128, Uint128) {
		a.H |= 1 << 63
		b.H |= 1 << 63
		a = a.Rsh(l % 128)
		b = b.Rsh(l % 128)
		if b == (Uint128{0, 0}) {
			return Uint128{0, 0}, Uint128{0, 0}
		}
		bigA := uint128ToBig(new(big.Int), a)
		bigB := uint128ToBig(new(big.Int), b)
		div, mod := new(big.Int).DivMod(bigA, bigB, new(big.Int))
		return bigToUint128(div), bigToUint128(mod)
	}
	if err := quick.CheckEqual(f, g, &quick.Config{
		MaxCountScale: 10000,
	}); err != nil {
		t.Error(err)
	}
}

func BenchmarkUint128_DivMod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		div, mod := uint128Input.DivMod(uint128Input)
		runtime.KeepAlive(div)
		runtime.KeepAlive(mod)
	}
}

func BenchmarkUint128_Div(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Div(uint128Input))
	}
}

func BenchmarkUint128_Mod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Mod(uint128Input))
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

func BenchmarkUint128_And(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.And(uint128Input))
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

func BenchmarkUint128_Or(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Or(uint128Input))
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

func BenchmarkUint128_Xor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Xor(uint128Input))
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

func BenchmarkUint128_AndNot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.AndNot(uint128Input))
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

func BenchmarkUint128_Not(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Not())
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

func BenchmarkUint128_Neg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Neg())
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
			Uint128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0},
			0,
			Uint128{0x1234_5678_9abc_def0, 0x1234_5678_9abc_def0},
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
		{
			Uint128{0xffff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			math.MaxUint,
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

func BenchmarkUint128_Lsh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Lsh(uint(i) % 128))
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

func BenchmarkUint128_Rsh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Rsh(uint(i) % 128))
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

func BenchmarkUint128_LeadingZeros(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.LeadingZeros())
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

func BenchmarkUint128_TrailingZeros(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.TrailingZeros())
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

func BenchmarkUint128_Len(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Len())
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

func BenchmarkUint128_OnesCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.OnesCount())
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

func BenchmarkUint128_RotateLeft(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.RotateLeft(i % 128))
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

func BenchmarkUint128_Reverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.Reverse())
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

func BenchmarkUint128_ReverseBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(uint128Input.ReverseBytes())
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
			0.5,
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
			// the maximum float64 value that can convert to 65-bits integer
			(1<<53 - 1) << 12,
			Uint128{1, 0xfffffffffffff000},
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

func BenchmarkFloat64ToUint128(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(Float64ToUint128(float64(i)))
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
			Uint128{0x4b3b4ca85a86c47a, 0x098a_2240_0000_0000},
			"100000000000000000000000000000000000000",
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

package int128

import (
	"math/big"
	"runtime"
	"testing"
	"testing/quick"
)

// int128Input is used for benchmarks to prevent compiler optimizations.
var int128Input = Int128{0, 42}

// int128ToBig converts x to a big.Int.
func int128ToBig(b *big.Int, x Int128) *big.Int {
	if b == nil {
		b = new(big.Int)
	}
	var neg bool
	if x.H < 0 {
		neg = true
		x = x.Neg()
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
	if neg {
		b = b.Neg(b)
	}
	return b
}

// bigToInt128 converts x to a Uint128.
func bigToInt128(x *big.Int) Int128 {
	var buf [16]byte
	z := new(big.Int).Mod(x, bigModUint128)
	z.FillBytes(buf[:])
	ret := Int128{
		H: (int64(buf[0]) << 56) | (int64(buf[1]) << 48) | (int64(buf[2]) << 40) | (int64(buf[3]) << 32) |
			(int64(buf[4]) << 24) | (int64(buf[5]) << 16) | (int64(buf[6]) << 8) | int64(buf[7]),
		L: (uint64(buf[8]) << 56) | (uint64(buf[9]) << 48) | (uint64(buf[10]) << 40) | (uint64(buf[11]) << 32) |
			(uint64(buf[12]) << 24) | (uint64(buf[13]) << 16) | (uint64(buf[14]) << 8) | uint64(buf[15]),
	}
	if z.Sign() < 0 {
		ret = ret.Neg()
	}
	return ret
}

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

func TestInt128_AddQuick(t *testing.T) {
	f := func(a, b Int128) Int128 {
		return a.Add(b)
	}
	g := func(a, b Int128) Int128 {
		bigA := int128ToBig(new(big.Int), a)
		bigB := int128ToBig(new(big.Int), b)
		bigA.Add(bigA, bigB)
		return bigToInt128(bigA)
	}
	if err := quick.CheckEqual(f, g, &quick.Config{
		MaxCountScale: 1000,
	}); err != nil {
		t.Error(err)
	}
}

func BenchmarkInt128_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Add(int128Input))
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

func TestInt128_SubQuick(t *testing.T) {
	f := func(a, b Int128) Int128 {
		return a.Sub(b)
	}
	g := func(a, b Int128) Int128 {
		bigA := int128ToBig(new(big.Int), a)
		bigB := int128ToBig(new(big.Int), b)
		bigA.Sub(bigA, bigB)
		return bigToInt128(bigA)
	}
	if err := quick.CheckEqual(f, g, &quick.Config{
		MaxCountScale: 1000,
	}); err != nil {
		t.Error(err)
	}
}

func BenchmarkInt128_Sub(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Sub(int128Input))
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
			// -1 * -1 = 1
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{0, 1},
		},
		{
			// -1 * 1 = -1
			Int128{-1, 0xffff_ffff_ffff_ffff},
			Int128{0, 1},
			Int128{-1, 0xffff_ffff_ffff_ffff},
		},
		{
			// -1 * 1<<64 = -1<<64
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

func TestInt128_MulQuick(t *testing.T) {
	f := func(a, b Int128) Int128 {
		return a.Mul(b)
	}
	g := func(a, b Int128) Int128 {
		bigA := int128ToBig(new(big.Int), a)
		bigB := int128ToBig(new(big.Int), b)
		bigA.Mul(bigA, bigB)
		return bigToInt128(bigA)
	}
	if err := quick.CheckEqual(f, g, &quick.Config{
		MaxCountScale: 1000,
	}); err != nil {
		t.Error(err)
	}
}

func BenchmarkInt128_Mul(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Mul(int128Input))
	}
}

func TestInt128_DivMod(t *testing.T) {
	testCases := []struct {
		a, b, div, mod Int128
	}{
		{
			Int128{0x101, 0x123},
			Int128{0, 0x100},
			Int128{0x1, 0x100_0000_0000_0001},
			Int128{0, 0x23},
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			Int128{1, 0},
			Int128{0, 0x7fff_ffff_ffff_ffff},
			Int128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0},
			Int128{1, 0xffff_ffff_ffff_ffff},
			Int128{0, 0x3fff_ffff_ffff_ffff},
			Int128{1, 0x3fff_ffff_ffff_ffff},
		},
		{
			Int128{0, 5},
			Int128{0, 2},
			Int128{0, 2},
			Int128{0, 1},
		},
		{
			Int128{0, 5},
			Int128{0, 2}.Neg(),
			Int128{0, 2}.Neg(),
			Int128{0, 1},
		},
		{
			Int128{0, 5}.Neg(),
			Int128{0, 2},
			Int128{0, 3}.Neg(),
			Int128{0, 1},
		},
		{
			Int128{0, 5}.Neg(),
			Int128{0, 2}.Neg(),
			Int128{0, 3},
			Int128{0, 1},
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

func TestInt128_DivModQuick(t *testing.T) {
	f := func(a, b Int128) (Int128, Int128) {
		if b == (Int128{0, 0}) {
			return Int128{0, 0}, Int128{0, 0}
		}
		return a.DivMod(b)
	}
	g := func(a, b Int128) (Int128, Int128) {
		if b == (Int128{0, 0}) {
			return Int128{0, 0}, Int128{0, 0}
		}
		bigA := int128ToBig(new(big.Int), a)
		bigB := int128ToBig(new(big.Int), b)
		div, mod := new(big.Int).DivMod(bigA, bigB, new(big.Int))
		return bigToInt128(div), bigToInt128(mod)
	}
	if err := quick.CheckEqual(f, g, &quick.Config{
		MaxCountScale: 1000,
	}); err != nil {
		t.Error(err)
	}
}

func BenchmarkInt128_DivMod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		div, mod := int128Input.DivMod(int128Input)
		runtime.KeepAlive(div)
		runtime.KeepAlive(mod)
	}
}

func BenchmarkInt128_Div(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Div(int128Input))
	}
}

func BenchmarkInt128_Mod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Mod(int128Input))
	}
}

func TestInt128_QuoRem(t *testing.T) {
	testCases := []struct {
		a, b, div, mod Int128
	}{
		{
			Int128{0x101, 0x123},
			Int128{0, 0x100},
			Int128{0x1, 0x100_0000_0000_0001},
			Int128{0, 0x23},
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0xffff_ffff_ffff_ffff},
			Int128{1, 0},
			Int128{0, 0x7fff_ffff_ffff_ffff},
			Int128{0, 0xffff_ffff_ffff_ffff},
		},
		{
			Int128{0x7fff_ffff_ffff_ffff, 0},
			Int128{1, 0xffff_ffff_ffff_ffff},
			Int128{0, 0x3fff_ffff_ffff_ffff},
			Int128{1, 0x3fff_ffff_ffff_ffff},
		},
		{
			Int128{0, 5},
			Int128{0, 2},
			Int128{0, 2},
			Int128{0, 1},
		},
		{
			Int128{0, 5},
			Int128{0, 2}.Neg(),
			Int128{0, 2}.Neg(),
			Int128{0, 1},
		},
		{
			Int128{0, 5}.Neg(),
			Int128{0, 2},
			Int128{0, 2}.Neg(),
			Int128{0, 1}.Neg(),
		},
		{
			Int128{0, 5}.Neg(),
			Int128{0, 2}.Neg(),
			Int128{0, 2},
			Int128{0, 1}.Neg(),
		},
	}

	for i, tc := range testCases {
		div, mod := tc.a.QuoRem(tc.b)
		if div != tc.div {
			t.Errorf("%d: %#v / %#v should %#v, but %#v", i, tc.a, tc.b, tc.div, div)
		}
		if mod != tc.mod {
			t.Errorf("%d: %#v %% %#v should %#v, but %#v", i, tc.a, tc.b, tc.mod, mod)
		}

		div = tc.a.Quo(tc.b)
		if div != tc.div {
			t.Errorf("%d: %#v / %#v should %#v, but %#v", i, tc.a, tc.b, tc.div, div)
		}
		mod = tc.a.Rem(tc.b)
		if mod != tc.mod {
			t.Errorf("%d: %#v %% %#v should %#v, but %#v", i, tc.a, tc.b, tc.mod, mod)
		}
	}
}

func TestInt128_QuoRemQuick(t *testing.T) {
	f := func(a, b Int128) (Int128, Int128) {
		if b == (Int128{0, 0}) {
			return Int128{0, 0}, Int128{0, 0}
		}
		return a.QuoRem(b)
	}
	g := func(a, b Int128) (Int128, Int128) {
		if b == (Int128{0, 0}) {
			return Int128{0, 0}, Int128{0, 0}
		}
		bigA := int128ToBig(new(big.Int), a)
		bigB := int128ToBig(new(big.Int), b)
		quo, rem := new(big.Int).QuoRem(bigA, bigB, new(big.Int))
		return bigToInt128(quo), bigToInt128(rem)
	}
	if err := quick.CheckEqual(f, g, &quick.Config{
		MaxCountScale: 1000,
	}); err != nil {
		t.Error(err)
	}
}

func BenchmarkInt128_QuoRem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		div, mod := int128Input.QuoRem(int128Input)
		runtime.KeepAlive(div)
		runtime.KeepAlive(mod)
	}
}

func BenchmarkInt128_Quo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Quo(int128Input))
	}
}

func BenchmarkInt128_Rem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Rem(int128Input))
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

func BenchmarkInt128_And(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.And(int128Input))
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

func BenchmarkInt128_Or(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Or(int128Input))
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

func BenchmarkInt128_Xor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Xor(int128Input))
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

func BenchmarkInt128_AndNot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.AndNot(int128Input))
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

func BenchmarkInt128_Not(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Not())
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

func BenchmarkInt128_Neg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Neg())
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

func BenchmarkInt128_Lsh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Lsh(uint(i) % 128))
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

func BenchmarkInt128_Rsh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(int128Input.Rsh(uint(i) % 128))
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
			0.5,
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
			// the maximum float64 value that can convert to 65-bits integer
			(1<<53 - 1) << 12,
			Int128{1, 0xfffffffffffff000},
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

func BenchmarkFloat64ToInt128(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(Float64ToInt128(float64(i)))
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

package int128

import (
	"math"
	"math/bits"
)

// Uint128 is a 128-bit unsigned integer.
type Uint128 struct {
	H uint64
	L uint64
}

// Add returns the sum a+b.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) Add(b Uint128) Uint128 {
	l, carry := bits.Add64(a.L, b.L, 0)
	h, _ := bits.Add64(a.H, b.H, carry)
	return Uint128{h, l}
}

// Sub returns the difference x-y.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) Sub(b Uint128) Uint128 {
	l, borrow := bits.Sub64(a.L, b.L, 0)
	h, _ := bits.Sub64(a.H, b.H, borrow)
	return Uint128{h, l}
}

// Mul returns the product x*y.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) Mul(b Uint128) Uint128 {
	h, l := bits.Mul64(a.L, b.L)
	// they are same as a.H * b.L, but use bits.Mul64 for avoiding compiler optimization
	_, h1 := bits.Mul64(a.H, b.L)
	_, h2 := bits.Mul64(a.L, b.H)
	return Uint128{h + h1 + h2, l}
}

// Div returns the quotient a/b for b != 0.
// If b == 0, a division-by-zero run-time panic occurs.
func (a Uint128) Div(b Uint128) Uint128 {
	if b.H == 0 {
		// optimize for uint128 / uint64
		h := a.H / b.L
		l, _ := bits.Div64(a.H%b.L, a.L, b.L)
		return Uint128{h, l}
	}

	n := uint(bits.LeadingZeros64(b.H))
	x := a.Rsh(1)
	y := b.Lsh(n)
	q, _ := bits.Div64(x.H, x.L, y.H)
	q >>= 63 - n
	if q > 0 {
		q--
	}

	h, l := bits.Mul64(b.L, q)
	h += b.H * q
	r := a.Sub(Uint128{h, l})
	if r.Cmp(b) >= 0 {
		q++
	}
	return Uint128{0, q}
}

// Mod returns the modulus x%y for y != 0.
// If y == 0, a division-by-zero run-time panic occurs.
func (a Uint128) Mod(b Uint128) Uint128 {
	if b.H == 0 {
		// optimize for uint128 / uint64
		_, rem := bits.Div64(a.H%b.L, a.L, b.L)
		return Uint128{0, rem}
	}

	n := uint(bits.LeadingZeros64(b.H))
	x := a.Rsh(1)
	y := b.Lsh(n)
	q, _ := bits.Div64(x.H, x.L, y.H)
	q >>= 63 - n
	if q > 0 {
		q--
	}

	h, l := bits.Mul64(b.L, q)
	h += b.H * q
	r := a.Sub(Uint128{h, l})
	if r.Cmp(b) >= 0 {
		r = r.Sub(b)
	}
	return r
}

// DivMod returns the quotient and remainder of a/b for b != 0.
// If b == 0, a division-by-zero run-time panic occurs.
//
//	q = a div b  such that
//	m = a - b*q  with 0 <= m < |y|
func (a Uint128) DivMod(b Uint128) (Uint128, Uint128) {
	if b.H == 0 {
		// optimize for uint128 / uint64
		h := a.H / b.L
		l, rem := bits.Div64(a.H%b.L, a.L, b.L)
		return Uint128{h, l}, Uint128{0, rem}
	}

	n := uint(bits.LeadingZeros64(b.H))
	x := a.Rsh(1)
	y := b.Lsh(n)
	q, _ := bits.Div64(x.H, x.L, y.H)
	q >>= 63 - n
	if q > 0 {
		q--
	}

	h, l := bits.Mul64(b.L, q)
	h += b.H * q
	r := a.Sub(Uint128{h, l})
	if r.Cmp(b) >= 0 {
		q++
		r = r.Sub(b)
	}
	return Uint128{0, q}, r
}

// Quo returns the quotient a/b for b != 0.
// If b == 0, a division-by-zero run-time panic occurs.
// Quo is the same as Div in Uint128.
func (a Uint128) Quo(b Uint128) Uint128 {
	return a.Div(b)
}

// Rem returns he remainder a%b for b != 0.
// If b == 0, a division-by-zero run-time panic occurs.
// Rem is the same as Mod in Uint128.
func (a Uint128) Rem(b Uint128) Uint128 {
	return a.Mod(b)
}

// QuoRem returns the quotient a/b and the remainder a%b for b != 0.
// a division-by-zero run-time panic occurs.
//
//	q = a/b      with the result truncated to zero
//	r = a - b*q
//
// QuoRem is the same as DivMod in Uint128.
func (a Uint128) QuoRem(b Uint128) (Uint128, Uint128) {
	return a.DivMod(b)
}

// Cmp compares a and b and returns:
//
//	-1 if a <  b
//	 0 if a == b
//	+1 if a >  b
func (a Uint128) Cmp(b Uint128) int {
	if a.H == b.H {
		if a.L == b.L {
			return 0
		} else if a.L > b.L {
			return 1
		} else {
			return -1
		}
	} else if a.H > b.H {
		return 1
	} else {
		return -1
	}
}

// And returns the bitwise AND a&b.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) And(b Uint128) Uint128 {
	return Uint128{a.H & b.H, a.L & b.L}
}

// Or returns the bitwise OR a|b.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) Or(b Uint128) Uint128 {
	return Uint128{a.H | b.H, a.L | b.L}
}

// Xor returns the bitwise XOR a^b.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) Xor(b Uint128) Uint128 {
	return Uint128{a.H ^ b.H, a.L ^ b.L}
}

// AndNot returns the bitwise AND NOT a&^b.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) AndNot(b Uint128) Uint128 {
	return Uint128{a.H &^ b.H, a.L &^ b.L}
}

// Not returns the bitwise NOT ^a.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) Not() Uint128 {
	return Uint128{^a.H, ^a.L}
}

// Neg returns the negation -a.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) Neg() Uint128 {
	l, borrow := bits.Sub64(0, a.L, 0)
	h, _ := bits.Sub64(0, a.H, borrow)
	return Uint128{h, l}
}

// Lsh returns the logical left shift a<<i.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) Lsh(i uint) Uint128 {
	// This operation may overflow, but it's okay because when it overflows,
	// the result is always greater than or equal to 64.
	// And shifts of 64 bits or more always result in 0, so they don't affect the final result.
	n := uint(i - 64)
	m := uint(64 - i)

	return Uint128{a.H<<i | a.L<<n | a.L>>m, a.L << i}
}

// Rsh returns the logical right shift a>>i.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) Rsh(i uint) Uint128 {
	// This operation may overflow, but it's okay because when it overflows,
	// the result is always greater than or equal to 64.
	// And shifts of 64 bits or more always result in 0, so they don't affect the final result.
	n := uint(i - 64)
	m := uint(64 - i)

	return Uint128{a.H >> i, a.H>>n | a.H<<m | a.L>>i}
}

// LeadingZeros returns the number of leading zero bits in a; the result is 128 for a == 0.
func (a Uint128) LeadingZeros() int {
	if a.H == 0 {
		return 64 + bits.LeadingZeros64(a.L)
	}
	return bits.LeadingZeros64(a.H)
}

// TrailingZeros returns the number of trailing zero bits in a; the result is 128 for a == 0.
func (a Uint128) TrailingZeros() int {
	if a.L == 0 {
		return 64 + bits.TrailingZeros64(a.H)
	}
	return bits.TrailingZeros64(a.L)
}

// Len returns the minimum number of bits required to represent a; the result is 0 for a == 0.
func (a Uint128) Len() int {
	if a.H == 0 {
		return bits.Len64(a.L)
	}
	return 64 + bits.Len64(a.H)
}

// OnesCount returns the number of one bits ("population count") in a.
func (a Uint128) OnesCount() int {
	return bits.OnesCount64(a.H) + bits.OnesCount64(a.L)
}

// RotateLeft returns the value of a rotated left by (k mod 128) bits.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) RotateLeft(k int) Uint128 {
	const n = 128
	s := uint(k) & (n - 1)
	t := n - s

	return Uint128{a.H<<s | a.L<<(s-64) | a.L>>(64-s) | a.H>>t, a.L<<s | a.H>>(t-64) | a.H<<(64-t) | a.L>>t}
}

// Reverse returns the value of a with its bits in reversed order.
func (a Uint128) Reverse() Uint128 {
	return Uint128{bits.Reverse64(a.L), bits.Reverse64(a.H)}
}

// ReverseBytes returns the value of a with its bytes in reversed order.
//
// This function's execution time does not depend on the inputs.
func (a Uint128) ReverseBytes() Uint128 {
	return Uint128{bits.ReverseBytes64(a.L), bits.ReverseBytes64(a.H)}
}

// Int128 returns a as a signed 128-bit integer.
func (a Uint128) Int128() Int128 {
	return Int128{int64(a.H), a.L}
}

// Float64ToUint128 returns the nearest Uint128 representation of v.
func Float64ToUint128(v float64) Uint128 {
	b := math.Float64bits(v)
	exp := int((b>>52)&0x7FF) - 1023
	frac := b&0xFFFFFFFFFFFFF | 0x10000000000000
	if exp < 0 {
		return Uint128{}
	}
	ret := Uint128{0, frac}
	if exp < 52 {
		ret = ret.Rsh(uint(52 - exp))
	} else {
		// Uint128 cannot represent values greater or equal 1 << 128,
		// however the spec says: https://golang.org/ref/spec#Conversions
		// > if the result type cannot represent the value the conversion succeeds
		// > but the result value is implementation-dependent.
		// so we don't care these case.
		ret = ret.Lsh(uint(exp - 52))
	}
	if b&0x8000000000000000 != 0 {
		ret = ret.Neg()
	}
	return ret
}

// Text returns the string representation of a in the given base.
// Base must be between 2 and 62, inclusive.
// The result uses the lower-case letters 'a' to 'z' for digit values 10 to 35,
// and the upper-case letters 'A' to 'Z' for digit values 36 to 61. No prefix (such as "0x") is added to the string.
func (a Uint128) Text(base int) string {
	if base == 10 && a.H == 0 && a.L < nSmalls {
		return small(int(a.L))
	}
	_, s := formatUint128(nil, a.H, a.L, base, false, false)
	return s
}

// Append appends the string representation of a, as generated by a.Text(base), to buf and returns the extended buffer.
func (a Uint128) Append(dst []byte, base int) []byte {
	if base == 10 && a.H == 0 && a.L < nSmalls {
		return append(dst, small(int(a.L))...)
	}
	d, _ := formatUint128(dst, a.H, a.L, base, false, true)
	return d
}

// String returns the decimal representation of a as generated by a.Text(10).
func (a Uint128) String() string {
	if a.H == 0 && a.L < nSmalls {
		return small(int(a.L))
	}
	_, s := formatUint128(nil, a.H, a.L, 10, false, false)
	return s
}

const nSmalls = 100

const smallsString = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

func small(n int) string {
	if n < 10 {
		return digits[n : n+1]
	}
	return smallsString[n*2 : n*2+2]
}

func formatUint128(dst []byte, h, l uint64, base int, neg bool, append_ bool) ([]byte, string) {
	if base < 2 || base > len(digits) {
		panic("int128: illegal Append/Format base")
	}

	var s [128 + 1]byte // +1 is for the sign
	i := len(s)

	if neg {
		var carry uint64
		l, carry = bits.Add64(^l, 1, 0)
		h, _ = bits.Add64(^h, 0, carry)
	}

	if base == 10 {
		// common case: use constants for / because
		// the compiler can optimize it into a multiply+shift
		for h != 0 {
			var r uint64
			l, r = bits.Div64(h%1e18, l, 1e18)
			h /= 1e18
			for j := 0; j < 9; j++ {
				is := (r % 100) * 2
				r /= 100
				i -= 2
				s[i+1] = smallsString[is+1]
				s[i+0] = smallsString[is+0]
			}
		}
		for l >= 100 {
			is := (l % 100) * 2
			l /= 100
			i -= 2
			s[i+1] = smallsString[is+1]
			s[i+0] = smallsString[is+0]
		}

		if l >= 10 {
			i--
			s[i] = digits[l%10]
			l /= 10
		}
		i--
		s[i] = digits[l]
	} else if isPowerOfTwo(base) {
		// Use shifts and masks instead of / and %.
		shift := uint(bits.TrailingZeros(uint(base))) & 7
		b := uint64(base)
		mask := uint(b) - 1 // == 1<<shift - 1
		for h != 0 {
			i--
			s[i] = digits[uint(l)&mask]
			l = h<<(64-shift) | l>>shift
			h >>= shift
		}
		for l >= b {
			i--
			s[i] = digits[uint(l)&mask]
			l >>= shift
		}
		// l < base
		i--
		s[i] = digits[uint(l)]
	} else {
		// general case
		b := uint64(base)
		for h != 0 {
			var r uint64
			l, r = bits.Div64(h%b, l, b)
			h /= b
			i--
			s[i] = digits[uint(r)]
		}
		for l >= b {
			i--
			q := l / b
			s[i] = digits[uint(l-q*b)]
			l = q
		}
		// l < base
		i--
		s[i] = digits[uint(l)]
	}

	// add the sign
	if neg {
		i--
		s[i] = '-'
	}

	if append_ {
		return append(dst, s[i:]...), ""
	}
	return nil, string(s[i:])
}

func isPowerOfTwo(x int) bool {
	return x&(x-1) == 0
}

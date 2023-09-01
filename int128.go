package int128

import (
	"math"
	"math/bits"
)

// Int128 is a signed 128-bit integer.
type Int128 struct {
	H int64
	L uint64
}

// Add returns the sum a+b.
//
// This function's execution time does not depend on the inputs.
func (a Int128) Add(b Int128) Int128 {
	l, carry := bits.Add64(a.L, b.L, 0)
	h, _ := bits.Add64(uint64(a.H), uint64(b.H), carry)
	return Int128{int64(h), l}
}

// Sub returns the difference x-y.
//
// This function's execution time does not depend on the inputs.
func (a Int128) Sub(b Int128) Int128 {
	l, borrow := bits.Sub64(a.L, b.L, 0)
	h, _ := bits.Sub64(uint64(a.H), uint64(b.H), borrow)
	return Int128{int64(h), l}
}

// Mul returns the product x*y.
func (a Int128) Mul(b Int128) Int128 {
	neg := false
	if a.H < 0 {
		neg = !neg
		a = a.Neg()
	}
	if b.H < 0 {
		neg = !neg
		b = b.Neg()
	}

	h, l := bits.Mul64(a.L, b.L)
	h1 := uint64(a.H) * b.L
	h2 := a.L * uint64(b.H)
	ret := Int128{int64(h + h1 + h2), l}

	if neg {
		ret = ret.Neg()
	}
	return ret
}

// Div returns the quotient a/b for b != 0.
// If b == 0, a division-by-zero run-time panic occurs.
func (a Int128) Div(b Int128) Int128 {
	var negA, negB bool
	if a.H < 0 {
		negA = true
		a = a.Neg()
	}
	if b.H < 0 {
		negB = true
		b = b.Neg()
	}

	div := a.Uint128().Div(b.Uint128())
	if negA {
		div = div.Add(Uint128{0, 1})
	}
	if negA != negB {
		div = div.Neg()
	}

	return div.Int128()

}

// Mod returns the modulus x%y for y != 0.
// If y == 0, a division-by-zero run-time panic occurs.
func (a Int128) Mod(b Int128) Int128 {
	var negA bool
	if a.H < 0 {
		negA = true
		a = a.Neg()
	}
	if b.H < 0 {
		b = b.Neg()
	}

	mod := a.Uint128().Mod(b.Uint128())
	if negA {
		mod = mod.Neg().Add(b.Uint128())
	}

	return mod.Int128()
}

// DivMod returns the quotient and remainder of a/b for b != 0.
// If b == 0, a division-by-zero run-time panic occurs.
//
// DivMod implements Euclidean division and modulus (unlike Go):
//
//	q = a div b  such that
//	m = a - b*q  with 0 <= m < |y|
func (a Int128) DivMod(b Int128) (Int128, Int128) {
	var negA, negB bool
	if a.H < 0 {
		negA = true
		a = a.Neg()
	}
	if b.H < 0 {
		negB = true
		b = b.Neg()
	}

	div, mod := a.Uint128().DivMod(b.Uint128())
	if negA != negB {
		div = div.Neg()
	}
	if negA {
		mod = mod.Neg().Add(b.Uint128())
		if negB {
			div = div.Add(Uint128{0, 1})
		} else {
			div = div.Sub(Uint128{0, 1})
		}
	}

	return div.Int128(), mod.Int128()
}

// Quo returns the quotient a/b for b != 0.
// If b == 0, a division-by-zero run-time panic occurs.
// Quo implements truncated division (like Go); see QuoRem for more details.
func (a Int128) Quo(b Int128) Int128 {
	var negA, negB bool
	if a.H < 0 {
		negA = true
		a = a.Neg()
	}
	if b.H < 0 {
		negB = true
		b = b.Neg()
	}

	div := a.Uint128().Div(b.Uint128())
	if negA != negB {
		div = div.Neg()
	}
	return div.Int128()
}

// Rem returns he remainder a%b for b != 0.
// If b == 0, a division-by-zero run-time panic occurs.
// Rem implements truncated modulus (like Go); see QuoRem for more details.
func (a Int128) Rem(b Int128) Int128 {
	var negA bool
	if a.H < 0 {
		negA = true
		a = a.Neg()
	}
	if b.H < 0 {
		b = b.Neg()
	}

	mod := a.Uint128().Mod(b.Uint128())
	if negA {
		mod = mod.Neg()
	}
	return mod.Int128()
}

// QuoRem returns the quotient a/b and the remainder a%b for b != 0.
// a division-by-zero run-time panic occurs.
//
// QuoRem implements T-division and modulus (like Go):
//
//	q = a/b      with the result truncated to zero
//	r = a - b*q
func (a Int128) QuoRem(b Int128) (Int128, Int128) {
	var negA, negB bool
	if a.H < 0 {
		negA = true
		a = a.Neg()
	}
	if b.H < 0 {
		negB = true
		b = b.Neg()
	}

	div, mod := a.Uint128().DivMod(b.Uint128())
	if negA != negB {
		div = div.Neg()
	}
	if negA {
		mod = mod.Neg()
	}

	return div.Int128(), mod.Int128()
}

// Cmp compares a and b and returns:
//
//	-1 if a <  b
//	 0 if a == b
//	+1 if a >  b
func (a Int128) Cmp(b Int128) int {
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
func (a Int128) And(b Int128) Int128 {
	return Int128{a.H & b.H, a.L & b.L}
}

// Or returns the bitwise OR a|b.
//
// This function's execution time does not depend on the inputs.
func (a Int128) Or(b Int128) Int128 {
	return Int128{a.H | b.H, a.L | b.L}
}

// Xor returns the bitwise XOR a^b.
//
// This function's execution time does not depend on the inputs.
func (a Int128) Xor(b Int128) Int128 {
	return Int128{a.H ^ b.H, a.L ^ b.L}
}

// AndNot returns the bitwise AND NOT a&^b.
//
// This function's execution time does not depend on the inputs.
func (a Int128) AndNot(b Int128) Int128 {
	return Int128{a.H &^ b.H, a.L &^ b.L}
}

// Not returns the bitwise NOT ^a.
//
// This function's execution time does not depend on the inputs.
func (a Int128) Not() Int128 {
	return Int128{^a.H, ^a.L}
}

// Neg returns the negation -a.
//
// This function's execution time does not depend on the inputs.
func (a Int128) Neg() Int128 {
	l, borrow := bits.Sub64(0, a.L, 0)
	h, _ := bits.Sub64(0, uint64(a.H), borrow)
	return Int128{int64(h), l}
}

// Lsh returns the logical left shift a<<i.
//
// This function's execution time does not depend on the inputs.
func (a Int128) Lsh(i uint) Int128 {
	n, v := bits.Sub(i, 64, 0)
	m := ^n + 1

	mask := uint64(int(v) - 1)
	return Int128{a.H<<i | int64(mask&(a.L<<n)) | int64(^mask&(a.L>>m)), a.L << i}
}

// Rsh returns the logical right shift a>>i.
//
// This function's execution time does not depend on the inputs.
func (a Int128) Rsh(i uint) Int128 {
	n, v := bits.Sub(i, 64, 0)
	m := ^n + 1

	mask := uint64(int(v) - 1)
	return Int128{a.H >> i, mask&uint64(a.H>>n) | ^mask&uint64(a.H<<m) | a.L>>i}
}

// Uint128 returns a as a unsigned 128-bit integer.
func (a Int128) Uint128() Uint128 {
	return Uint128{uint64(a.H), a.L}
}

// Float64ToUint128 returns the nearest Uint128 representation of v.
func Float64ToInt128(v float64) Int128 {
	b := math.Float64bits(v)
	exp := int((b>>52)&0x7FF) - 1023
	frac := b&0xFFFFFFFFFFFFF | 0x10000000000000
	if exp < 0 {
		return Int128{}
	}
	ret := Int128{0, frac}
	if exp < 52 {
		ret = ret.Rsh(uint(52 - exp))
	} else {
		// Int128 cannot represent values greater or equal 1 << 128,
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
func (a Int128) Text(base int) string {
	if base == 10 && a.H == 0 && a.L < nSmalls {
		return small(int(a.L))
	}
	_, s := formatUint128(nil, uint64(a.H), a.L, base, a.H < 0, false)
	return s
}

// Append appends the string representation of a, as generated by a.Text(base), to buf and returns the extended buffer.
func (a Int128) Append(dst []byte, base int) []byte {
	if base == 10 && a.H == 0 && a.L < nSmalls {
		return append(dst, small(int(a.L))...)
	}
	d, _ := formatUint128(dst, uint64(a.H), a.L, base, a.H < 0, true)
	return d
}

// String returns the decimal representation of a as generated by a.Text(10).
func (a Int128) String() string {
	if a.H == 0 && a.L < nSmalls {
		return small(int(a.L))
	}
	_, s := formatUint128(nil, uint64(a.H), a.L, 10, a.H < 0, false)
	return s
}

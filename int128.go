package int128

import "math/bits"

type Int128 struct {
	H int64
	L uint64
}

func (a Int128) Add(b Int128) Int128 {
	l, carry := bits.Add64(a.L, b.L, 0)
	h, _ := bits.Add64(uint64(a.H), uint64(b.H), carry)
	return Int128{int64(h), l}
}

func (a Int128) Sub(b Int128) Int128 {
	l, borrow := bits.Sub64(a.L, b.L, 0)
	h, _ := bits.Sub64(uint64(a.H), uint64(b.H), borrow)
	return Int128{int64(h), l}
}

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

func (a Int128) And(b Int128) Int128 {
	return Int128{a.H & b.H, a.L & b.L}
}

func (a Int128) Or(b Int128) Int128 {
	return Int128{a.H | b.H, a.L | b.L}
}

func (a Int128) Xor(b Int128) Int128 {
	return Int128{a.H ^ b.H, a.L ^ b.L}
}

func (a Int128) AndNot(b Int128) Int128 {
	return Int128{a.H &^ b.H, a.L &^ b.L}
}

func (a Int128) Not() Int128 {
	return Int128{^a.H, ^a.L}
}

func (a Int128) Neg() Int128 {
	l, carry := bits.Add64(^a.L, 1, 0)
	h, _ := bits.Add64(uint64(^a.H), 0, carry)
	return Int128{int64(h), l}
}

func (a Int128) Lsh(i uint) Int128 {
	if i < 64 {
		return Int128{a.H<<i | int64(a.L>>(64-i)), a.L << i}
	} else {
		return Int128{int64(a.L << (i - 64)), 0}
	}
}

func (a Int128) Rsh(i uint) Int128 {
	if i < 64 {
		return Int128{a.H >> i, uint64(a.H<<(64-i)) | a.L>>i}
	} else {
		return Int128{a.H >> 63, uint64(a.H >> (i - 64))}
	}
}

func (a Int128) Uint128() Uint128 {
	return Uint128{uint64(a.H), a.L}
}

func Float64ToInt128(v float64) Int128 {
	neg := false
	if v < 0 {
		neg = true
		v = -v
	}
	var ret Int128
	if v < 1<<64 {
		ret = Int128{0, uint64(v)}
	} else {
		// Int128 cannot represent values greater or equal 1 << 128,
		// however the spec says: https://golang.org/ref/spec#Conversions
		// > if the result type cannot represent the value the conversion succeeds
		// > but the result value is implementation-dependent.
		// so we don't care these case.
		ret = Int128{int64(v / (1 << 64)), 0}
	}
	if neg {
		ret = ret.Neg()
	}
	return ret
}

func (a Int128) Text(base int) string {
	if base == 10 && a.H == 0 && a.L < nSmalls {
		return small(int(a.L))
	}
	_, s := formatUint128(nil, uint64(a.H), a.L, base, a.H < 0, false)
	return s
}

func (a Int128) Append(dst []byte, base int) []byte {
	if base == 10 && a.H == 0 && a.L < nSmalls {
		return append(dst, small(int(a.L))...)
	}
	d, _ := formatUint128(dst, uint64(a.H), a.L, base, a.H < 0, true)
	return d
}

func (a Int128) String() string {
	if a.H == 0 && a.L < nSmalls {
		return small(int(a.L))
	}
	_, s := formatUint128(nil, uint64(a.H), a.L, 10, a.H < 0, false)
	return s
}

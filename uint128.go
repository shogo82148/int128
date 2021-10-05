package int128

import "math/bits"

type Uint128 struct {
	H uint64
	L uint64
}

func (a Uint128) Add(b Uint128) Uint128 {
	l, carry := bits.Add64(a.L, b.L, 0)
	h, _ := bits.Add64(a.H, b.H, carry)
	return Uint128{h, l}
}

func (a Uint128) Sub(b Uint128) Uint128 {
	l, borrow := bits.Sub64(a.L, b.L, 0)
	h, _ := bits.Sub64(a.H, b.H, borrow)
	return Uint128{h, l}
}

func (a Uint128) Mul(b Uint128) Uint128 {
	h, l := bits.Mul64(a.L, b.L)
	h1 := a.H * b.L
	h2 := a.L * b.H
	return Uint128{h + h1 + h2, l}
}

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

func (a Uint128) And(b Uint128) Uint128 {
	return Uint128{a.H & b.H, a.L & b.L}
}

func (a Uint128) Or(b Uint128) Uint128 {
	return Uint128{a.H | b.H, a.L | b.L}
}

func (a Uint128) Xor(b Uint128) Uint128 {
	return Uint128{a.H ^ b.H, a.L ^ b.L}
}

func (a Uint128) AndNot(b Uint128) Uint128 {
	return Uint128{a.H &^ b.H, a.L &^ b.L}
}

func (a Uint128) Not() Uint128 {
	return Uint128{^a.H, ^a.L}
}

func (a Uint128) Neg() Uint128 {
	l, carry := bits.Add64(^a.L, 1, 0)
	h, _ := bits.Add64(^a.H, 0, carry)
	return Uint128{h, l}
}

func (a Uint128) Lsh(i uint) Uint128 {
	if i < 64 {
		return Uint128{a.H<<i | a.L>>(64-i), a.L << i}
	} else {
		return Uint128{a.L << (i - 64), 0}
	}
}

func (a Uint128) Rsh(i uint) Uint128 {
	if i < 64 {
		return Uint128{a.H >> i, a.H<<(64-i) | a.L>>i}
	} else {
		return Uint128{0, a.H >> (i - 64)}
	}
}

func (a Uint128) LeadingZeros() int {
	if a.H == 0 {
		return 64 + bits.LeadingZeros64(a.L)
	}
	return bits.LeadingZeros64(a.H)
}

func (a Uint128) TrailingZeros() int {
	if a.L == 0 {
		return 64 + bits.TrailingZeros64(a.H)
	}
	return bits.TrailingZeros64(a.L)
}

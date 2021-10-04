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

	// the lower limit of the result
	ret := a.H / (b.H + 1)

	// TODO: fix me
	h, l := bits.Mul64(b.L, ret)
	h += b.H * ret
	if h < a.H || (h == a.H && l <= b.L) {
		return Uint128{0, ret}
	}
	return Uint128{0, ret}
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

func (a Uint128) Lsh(i int) Uint128 {
	h := a.H << i
	if i <= 64 {
		h |= a.L >> (64 - i)
	} else {
		h |= a.L << (i - 64)
	}
	return Uint128{h, a.L << i}
}

func (a Uint128) Rsh(i int) Uint128 {
	l := a.L >> i
	if i <= 64 {
		l |= a.H << (64 - i)
	} else {
		l |= a.H >> (i - 64)
	}
	return Uint128{a.H >> i, l}
}

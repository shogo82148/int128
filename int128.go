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

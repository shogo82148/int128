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

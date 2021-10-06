package int128

import (
	"math/bits"
)

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

func (a Uint128) Len() int {
	if a.H == 0 {
		return bits.Len64(a.L)
	}
	return 64 + bits.Len64(a.H)
}

func (a Uint128) OnesCount() int {
	return bits.OnesCount64(a.H) + bits.OnesCount64(a.L)
}

func (a Uint128) RotateLeft(k int) Uint128 {
	const n = 128
	s := uint(k) & (n - 1)
	t := n - s
	if s < 64 {
		return Uint128{a.H<<s | a.L>>(64-s), a.L<<s | a.H>>(t-64)}
	} else {
		return Uint128{a.L<<(s-64) | a.H>>t, a.H<<(64-t) | a.L>>t}
	}
}

func (a Uint128) Reverse() Uint128 {
	return Uint128{bits.Reverse64(a.L), bits.Reverse64(a.H)}
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

func (a Uint128) String() string {
	if a.H == 0 && a.L < nSmalls {
		return small(int(a.L))
	}

	// 1e19 is available on uint64
	// however the exponent should be even, because we handle twe digits at time.
	const power10 = 1e18

	var s [40]byte
	i := len(s)

	h, l := a.H, a.L
	for h != 0 {
		var r uint64
		l, r = bits.Div64(h%power10, l, power10)
		h /= power10

		for r > 0 {
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

	return string(s[i:])
}

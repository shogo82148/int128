package int128

import (
	"bytes"
	"fmt"
)

var _ fmt.Formatter = Int128{}

// Format implements [fmt.Formatter].
func (a Int128) Format(s fmt.State, verb rune) {
	var out []byte
	var prefix []byte

	if verb == 'v' {
		if s.Flag('#') {
			fmt.Fprintf(s, "Int128{H: %#016x, L: %#016x}", a.H, a.L)
		} else {
			out = a.Append(out, 10)
			s.Write(out)
		}
		return
	}

	if s.Flag('+') {
		if a.H >= 0 {
			prefix = append(prefix, '+')
		} else {
			prefix = append(prefix, '-')
			a = a.Neg()
		}
	} else if s.Flag(' ') {
		if a.H >= 0 {
			prefix = append(prefix, ' ')
		} else {
			prefix = append(prefix, '-')
			a = a.Neg()
		}
	} else {
		if a.H < 0 {
			prefix = append(prefix, '-')
			a = a.Neg()
		}
	}

	switch verb {
	case 'b':
		out = a.Append(out, 2)
		if s.Flag('#') {
			prefix = append(prefix, '0', 'b')
		}
	case 'o':
		out = a.Append(out, 8)
		if s.Flag('#') && !(len(out) > 0 && out[0] == '0') {
			prefix = append(prefix, '0')
		}
	case 'O':
		out = a.Append(out, 8)
		prefix = append(prefix, '0', 'o')
		if s.Flag('#') && !(len(out) > 0 && out[0] == '0') {
			prefix = append(prefix, '0')
		}
	case 'd':
		out = a.Append(out, 10)
	case 'x':
		out = a.Append(out, 16)
		if s.Flag('#') {
			prefix = append(prefix, '0', 'x')
		}
	case 'X':
		out = a.Append(out, 16)
		out = bytes.ToUpper(out)
		if s.Flag('#') {
			prefix = append(prefix, '0', 'X')
		}
	}

	if w, ok := s.Width(); ok {
		var buf [8]byte
		if s.Flag('0') {
			if len(prefix) > 0 {
				s.Write(prefix)
			}

			// pad with zeros
			buf[0] = '0'
			for i := len(prefix) + len(out); i < w; i++ {
				s.Write(buf[:1])
			}
			s.Write(out)
		} else if s.Flag('-') {
			if len(prefix) > 0 {
				s.Write(prefix)
			}
			s.Write(out)

			// pad with spaces
			buf[0] = ' '
			for i := len(prefix) + len(out); i < w; i++ {
				s.Write(buf[:1])
			}
		} else {
			// pad with spaces
			buf[0] = ' '
			for i := len(prefix) + len(out); i < w; i++ {
				s.Write(buf[:1])
			}
			if len(prefix) > 0 {
				s.Write(prefix)
			}
			s.Write(out)
		}
		return
	}

	if len(prefix) > 0 {
		s.Write(prefix)
	}
	s.Write(out)
}

var _ fmt.Formatter = Uint128{}

// Format implements [fmt.Formatter].
func (a Uint128) Format(s fmt.State, ch rune) {
	var out []byte
	var prefix []byte

	if s.Flag('+') {
		prefix = append(prefix, '+')
	} else if s.Flag(' ') {
		prefix = append(prefix, ' ')
	}

	switch ch {
	case 'b':
		out = a.Append(out, 2)
		if s.Flag('#') {
			prefix = append(prefix, '0', 'b')
		}
	case 'o':
		out = a.Append(out, 8)
		if s.Flag('#') && !(len(out) > 0 && out[0] == '0') {
			prefix = append(prefix, '0')
		}
	case 'O':
		out = a.Append(out, 8)
		prefix = append(prefix, '0', 'o')
		if s.Flag('#') && !(len(out) > 0 && out[0] == '0') {
			prefix = append(prefix, '0')
		}
	case 'd':
		out = a.Append(out, 10)
	case 'x':
		out = a.Append(out, 16)
		if s.Flag('#') {
			prefix = append(prefix, '0', 'x')
		}
	case 'X':
		out = a.Append(out, 16)
		out = bytes.ToUpper(out)
		if s.Flag('#') {
			prefix = append(prefix, '0', 'X')
		}
	case 'v':
		if s.Flag('#') {
			fmt.Fprintf(s, "Uint128{H: %#016x, L: %#016x}", a.H, a.L)
		} else {
			out = a.Append(out, 10)
			s.Write(out)
		}
		return
	}

	if w, ok := s.Width(); ok {
		var buf [8]byte
		if s.Flag('0') {
			if len(prefix) > 0 {
				s.Write(prefix)
			}

			// pad with zeros
			buf[0] = '0'
			for i := len(prefix) + len(out); i < w; i++ {
				s.Write(buf[:1])
			}
			s.Write(out)
		} else if s.Flag('-') {
			if len(prefix) > 0 {
				s.Write(prefix)
			}
			s.Write(out)

			// pad with spaces
			buf[0] = ' '
			for i := len(prefix) + len(out); i < w; i++ {
				s.Write(buf[:1])
			}
		} else {
			// pad with spaces
			buf[0] = ' '
			for i := len(prefix) + len(out); i < w; i++ {
				s.Write(buf[:1])
			}
			s.Write(prefix)
			s.Write(out)
		}
		return
	}

	if len(prefix) > 0 {
		s.Write(prefix)
	}
	s.Write(out)
}

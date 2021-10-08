package int128

func (a Int128) MarshalText() ([]byte, error) {
	text := a.Append(nil, 10)
	return text, nil
}

func (a Int128) MarshalJSON() ([]byte, error) {
	text := a.Append(nil, 10)
	return text, nil
}

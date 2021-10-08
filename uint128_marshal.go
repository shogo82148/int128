package int128

func (a Uint128) MarshalText() ([]byte, error) {
	text := a.Append(nil, 10)
	return text, nil
}

func (a Uint128) MarshalJSON() ([]byte, error) {
	text := a.Append(nil, 10)
	return text, nil
}

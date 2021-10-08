package int128

import (
	"encoding"
	"encoding/json"
	"testing"
)

var _ = json.Marshaler(Uint128{})
var _ = encoding.TextMarshaler(Uint128{})

func TestUint128_MarshalJSON(t *testing.T) {
	a := Uint128{0, 12345}
	data, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "12345" {
		t.Errorf("want %q, got %q", "12345", string(data))
	}
}

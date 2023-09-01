package int128_test

import (
	"fmt"

	"github.com/shogo82148/int128"
)

func ExampleInt128_Add() {
	a := int128.Int128{0, 1} // = 1
	b := int128.Int128{0, 2} // = 2
	c := a.Add(b)
	fmt.Println(c)
	// Output: 3
}

func ExampleInt128_Sub() {
	a := int128.Int128{0, 1} // = 1
	b := int128.Int128{0, 2} // = 2
	c := a.Sub(b)
	fmt.Println(c)
	// Output: -1
}

func ExampleInt128_Mul() {
	a := int128.Int128{0, 1} // = 1
	b := int128.Int128{0, 2} // = 2
	c := a.Mul(b)
	fmt.Println(c)
	// Output: 2
}

func ExampleUint128_Add() {
	a := int128.Uint128{0, 1} // = 1
	b := int128.Uint128{0, 2} // = 2
	c := a.Add(b)
	fmt.Println(c)
	// Output: 3
}

func ExampleUint128_Sub() {
	a := int128.Uint128{0, 3} // = 3
	b := int128.Uint128{0, 2} // = 2
	c := a.Sub(b)
	fmt.Println(c)
	// Output: 1
}

# int128

[![test](https://github.com/shogo82148/int128/actions/workflows/test.yml/badge.svg)](https://github.com/shogo82148/int128/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/shogo82148/int128.svg)](https://pkg.go.dev/github.com/shogo82148/int128)

The package `int128` provides 128-bit integer types.

## SYNOPSIS

```go
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
```

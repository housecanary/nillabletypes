# nillabletypes - A library of nil- and initialization-aware types in Go

nillabletypes is a library of Go types that provide nil and initialized states.

The nil state is useful when you'd like to communicate that a scalar value does
not exist.

The initialized state is useful when you'd like to ensure that a variable's
value has been set, typically after you pass an existing variable by reference
to some other function.

## Available Types

* `Bool`: represents a nil-able `bool` type.
* `Float`: represents a nil-able `float64` type.
* `Int32`: represents a nil-able `int32` type.
* `Int64`: represents a nil-able `int64` type.
* `Int`: a typealias for either Int32 or Int64, depending on whether the target
  architecture is 32- or 64-bit.
* `String`: represents a nil-able `string` type.
* `Date`: represents a nil-able date encoded as an ISO string.

## Supported Interfaces

Each type satisfies the following interfaces:

* [fmt/Stringer](https://golang.org/pkg/fmt/#Stringer)
* [encoding/json/Unmarshaler](https://golang.org/pkg/encoding/json/#Unmarshaler)
* [encoding/json/Marshaler](https://golang.org/pkg/encoding/json/#Marshaler)
* [database/sql/driver/Valuer](https://golang.org/pkg/database/sql/driver/#Valuer)
* [database/sql/Scanner](https://golang.org/pkg/database/sql/#Scanner)

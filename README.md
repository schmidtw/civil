# civil

[![Go Report Card](https://goreportcard.com/badge/github.com/schmidtw/civil)](https://goreportcard.com/report/github.com/schmidtw/civil)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/schmidtw/civil/blob/main/LICENSE)
[![GitHub Release](https://img.shields.io/github/release/schmidtw/civil.svg)](https://github.com/schmidtw/civil/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/schmidtw/civil)](https://pkg.go.dev/github.com/schmidtw/civil)

A zero-dependency Go package for representing civil dates.

This is a standalone fork of Google's [`cloud.google.com/go/civil`](https://pkg.go.dev/cloud.google.com/go/civil)
package, extracted so you can use civil date types without pulling in the entire
Google Cloud SDK dependency tree.

## Installation

```bash
go get github.com/schmidtw/civil
```

## Usage

See [example_test.go](example_test.go) for a complete working example.

The `Date` type implements `encoding.TextMarshaler` and `encoding.TextUnmarshaler`,
so it works with JSON and other text-based encodings out of the box.

## Why?

The original `cloud.google.com/go/civil` package lives inside the Google Cloud Go
SDK, which brings in a large dependency graph. If all you need is a civil date
type, this package gives you the same API with **zero external dependencies**.

## Differences from the original

- `ParseUint32Date` -- parse dates from `uint32` values in `YYYYMMDD` format
- `MarshalText` / `UnmarshalText` -- zero-value `Date{}` marshals to `""` and
  both `""` and `"0000-00-00"` unmarshal to the zero value instead of returning
  an error
- The `database/sql` `Scanner` and `driver.Valuer` interfaces are intentionally
  omitted to keep the package free of dependencies beyond the standard library

## License

Apache-2.0. See [LICENSE](LICENSE) for details.

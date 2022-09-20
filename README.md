# dwolla-v2-go

[![Build Status][1]][2] [![Code Coverage][3]][4] [![GoDoc][5]][6] [![MIT][7]][8] [![Go Report Card][9]][10]

[1]: https://circleci.com/gh/bluexlab/dwolla-v2-go.svg?style=svg
[2]: https://circleci.com/gh/bluexlab/dwolla-v2-go
[3]: https://codecov.io/gh/bluexlab/dwolla-v2-go/branch/master/graph/badge.svg
[4]: https://codecov.io/gh/bluexlab/dwolla-v2-go
[5]: https://godoc.org/github.com/bluexlab/dwolla-v2-go?status.svg
[6]: https://godoc.org/github.com/bluexlab/dwolla-v2-go
[7]: https://img.shields.io/badge/License-MIT-yellow.svg
[8]: LICENSE
[9]: https://goreportcard.com/badge/github.com/bluexlab/dwolla-v2-go
[10]: https://goreportcard.com/report/github.com/bluexlab/dwolla-v2-go

A Go wrapper for the Dwolla API V2

## Requirements

* Go v1.11+ (uses modules)

## Install

```bash
go get -u github.com/bluexlab/dwolla-v2-go
```

## Usage

To instantiate the client:

```go
package main

import (
	"context"
	"fmt"

	"github.com/bluexlab/dwolla-v2-go"
)

var ctx = context.Background()

func main() {
	client := dwolla.New("<your dwolla key here>", "<your dwolla secret here>", dwolla.Production)

	# Or if using the Dwolla sandbox
	#client := dwolla.New("<your dwolla key here>", "<your dwolla secret here>", dwolla.Sandbox)
}
```

To retrieve dwolla account information:

```go
res, err := client.Account.Retrieve(ctx)

if err != nil {
	fmt.Println("Error:", err)
	return err
}

fmt.Println("Account ID:", res.ID)
fmt.Println("Account Name:", res.Name)
```

See the [GoDoc](https://godoc.org/github.com/bluexlab/dwolla-v2-go) for the full API.

## License

MIT License

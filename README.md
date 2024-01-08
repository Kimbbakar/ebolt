# Ebolt
[![Go Report Card](https://goreportcard.com/badge/github.com/Kimbbakar/ebolt)](https://goreportcard.com/report/github.com/Kimbbakar/ebolt)
[![GoDoc](https://pkg.go.dev/badge/github.com/kimbbakar/ebolt?status.svg)](https://pkg.go.dev/github.com/kimbbakar/ebolt?tab=doc)

`Ebolt` is a lightweight Go package designed as a versatile wrapper for the bbolt key-value store. It offers a simplified interface with generic PUT, GET, and DELETE methods, along with the added convenience of expiration for key-value pairs.

## Installation

To use `ebolt` in your Go project, you can install it using the following `go get` command:

```bash
go get github.com/kimbbakar/ebolt
```

## Usage
Here's a quick example of how to use ebolt in your Go code:

```
package main

import (
	"time"

	ebolt "github.com/kimbbakar/ebolt"
)

func main() {

	// Init ebolt befor every restart to sweep expired key.
	// Passing bucket name will create bucket
	ebolt.InitEbolt(nil)

	// key-value without any expiry
	ebolt.DB.Put("my-key", "my-value", nil)

	// key-value with any expiry
	exp := time.Minute * 10
	ebolt.DB.Put("my-key", "my-value", &exp)

	// Delete key
	ebolt.DB.Delete("my-key")

	// Get Key
	ebolt.DB.Get("my-key")
}
```

## Features
- Generic PUT, GET, DELETE operations
- Expiration for key-value pairs

## Contributing
Contributions are welcome! If you have ideas for improvements or find any issues, feel free to open an issue or submit a pull request.

# License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

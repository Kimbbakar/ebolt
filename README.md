# ebolt

`ebolt` is a simple Go package that provides a generic wrapper on top of the `bbolt` key-value store with additional features like PUT, GET, DELETE, and expiration.

## Installation

To use `ebolt` in your Go project, you can install it using the following `go get` command:

```bash
go get github.com/kimbbakar/ebolt
```

Usage
Here's a quick example of how to use ebolt in your Go code:

```
package main

import (
	"time"

	eblot "github.com/kimbbakar/ebolt"
)

func main() {

	// Init ebolt befor every restart to sweep expired key.
	// Passing bucket name will create bucket
	eblot.InitEbolt(nil)

	// key-value without any expiry
	eblot.GetEbolt().Put("my-key", "my-value", nil)

	// key-value with any expiry
	exp := time.Minute * 10
	eblot.GetEbolt().Put("my-key", "my-value", &exp)

	// Delete key
	eblot.GetEbolt().Delete("my-key")

	// Get Key
	eblot.GetEbolt().Get("my-key")
}
```

## Features
- Generic PUT, GET, DELETE operations
- Key-Value store on top of bbolt
- Expiration for key-value pairs

## Contributing
Contributions are welcome! If you have ideas for improvements or find any issues, feel free to open an issue or submit a pull request.

# License
This project is licensed under the MIT License - see the LICENSE file for details.

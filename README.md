# Temporalink

Inspired by `temporalite` and copy/pasted heavily from `temporalio/cli`, this library allows you to run a Temporal instance embedded with your application.  This is experimental and not intended for production workloads.

## Example

```go
package main

import "fmt"
import "github.com/bmorton/temporalink"

func main() {
	ip := "0.0.0.0"
	serverPort := 7233
	uiPort := 8080
	
	s, err := temporalink.NewEmbeddedTemporal(ip, serverPort, uiPort)
	if err != nil {
		panic(err)
	}

	if err := s.Start(); err != nil {
		fmt.Printf("error while starting server: %w\n", err)
	}
}
```

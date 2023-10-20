# Temporalink

Inspired by `temporalite` and copy/pasted heavily from `temporalio/cli`, this library allows you to run a Temporal instance embedded with your application.  This is experimental and not intended for production workloads.

## Example

Install the `metrics` dependency and the `temporalink` package:
```
$ go get go.opentelemetry.io/otel/sdk/metric@v0.39.0
$ go get github.com/bmorton/temporalink
```

Then run a server:
```go
package main

import (
	"context"
	"fmt"

	"github.com/bmorton/temporalink"
)

func main() {
	ip := "0.0.0.0"
	serverPort := 7233
	uiPort := 8080

	s, err := temporalink.NewEmbeddedTemporal(ip, serverPort, uiPort)
	if err != nil {
		panic(err)
	}

	if err := s.Start(context.Background()); err != nil {
		fmt.Printf("error while starting server: %s\n", err)
	}
}
```

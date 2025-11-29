# Route122

Route122 is a Go routing library that vendors Go 1.22's new mux functionality, designed to be easily integrated with non-net/http handler frameworks. It provides the same pattern matching and routing capabilities as Go's standard library while being framework-agnostic.

## Installation

```bash
go get github.com/aisk/route122@latest
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/aisk/route122"
)

func main() {
    router := route122.New()

    // Register a route
    router.Handle("GET /users/{id}", "userHandler")

    // Match a request
    match, found := router.Match("GET", "", "/users/123")
    if found {
        fmt.Printf("Found route: %s\n", match.Pattern)
        fmt.Printf("Parameters: %v\n", match.Params)
        fmt.Printf("Handler: %v\n", match.Handler)
    }
}
```

For detailed pattern syntax and advanced usage, please refer to Go's official documentation on pattern-based routing.

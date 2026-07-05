# Doiit Middleware

`doiitware` is the shareable Go package for the Doiit platform. It contains the reusable Gin middleware and helper primitives that should behave the same way in every service.

## What It Provides

- `NewIPRateLimiter` for per-IP token bucket limiting.
- `RateLimitMiddleware` with configurable CIDR allowlists.
- `CustomRecoveryMiddleware` for safe panic recovery responses.
- `SecurityHeadersMiddleware` for standard browser security headers.
- `RequestLoggerMiddleware` for request IDs and latency logging.
- `RequestSizeMiddleware` for body size enforcement.
- `SanitizeInputMiddleware` for basic input hardening.
- CORS helpers and JWT claim utilities used by the sample services.

## Package Layout

```text
doiitware/
	go.mod
	middleware.go
	src/
		config/
		cors/
		helpers/
		middlewares/
```

## Installation

If the package is published from GitHub, downstream services should import it using the repository path and a version tag, for example:

```go
import doiitware "github.com/fasunle/doiitware"
```

If you are working in the current workspace, the local module path is `doiitware`.

## Quick Start

```go
package main

import (
	doiitware "doiitware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	limiter := doiitware.NewIPRateLimiter(100, 200)

	router.Use(doiitware.CustomRecoveryMiddleware())
	router.Use(doiitware.RateLimitMiddleware(limiter, nil))
	router.Use(doiitware.SecurityHeadersMiddleware("https://api.example.com"))

	_ = router.Run()
}
```

## Example Use Cases

### Shared edge protection

Use the rate limiter and recovery middleware in `api-gateway` so the gateway and all services return consistent responses under load or failure.

### New microservice bootstrap

Start a new service from `service-starter` and import the shared middleware package instead of reimplementing request logging, recovery, and security headers.

### User-facing browser APIs

Use the security headers and CORS helpers in browser-facing endpoints to keep a consistent policy across the platform.

## Testing

The package includes tests for:

- limiter reuse for the same IP,
- distinct limiters for different IPs,
- whitelist bypass behavior,
- 429 responses when the limiter rejects a request,
- panic recovery returning a safe JSON response.

Run the suite with:

```bash
go test ./...
```

## Release Process

The repository includes a GitHub Actions workflow that:

- validates the package with `gofmt`, `go test`, and `go vet`,
- uses Release Please to create versioned releases.

For Go libraries, a Git tag is the publication boundary. Downstream services should depend on a tagged release instead of a branch.

## Conventions

- Keep reusable code here only if it is genuinely shared across services.
- Prefer small constructors and middleware functions that are easy to compose.
- Add tests for any shared behavior before adding more consumers.

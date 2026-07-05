// Package helpers contains reusable primitives used by the shared doiitware middleware stack.
//
// These helpers are intentionally small and composable. They are useful when a service needs
// the same request ID generation, self-signed certificate creation, or per-IP rate limiter
// behavior as the other Doiit services.
//
// Common use cases include:
//
//   - generating request IDs for tracing and logs
//   - creating a local TLS certificate for development and test environments
//   - maintaining one token bucket limiter per client IP
package helpers

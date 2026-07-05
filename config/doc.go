// Package config defines shared constants used by the doiitware middleware and helper packages.
//
// The values in this package are intentionally centralized so that services can share
// the same security policy, timeout behavior, rate limiting defaults, and JWT settings.
//
// Common use cases include:
//
//   - customizing middleware behavior from a single configuration source
//   - keeping security header names consistent across services
//   - reusing JWT issuer, audience, and expiry defaults for token generation
package config

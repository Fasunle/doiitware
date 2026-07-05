package config

import "time"

const (
	// Security Headers
	HeaderStrictTransport   = "Strict-Transport-Security"
	HeaderContentSecurity   = "Content-Security-Policy"
	HeaderXFrameOptions     = "X-Frame-Options"
	HeaderXContentType      = "X-Content-Type-Options"
	HeaderReferrerPolicy    = "Referrer-Policy"
	HeaderPermissionsPolicy = "Permissions-Policy"
	HeaderXSSProtection     = "X-XSS-Protection"

	// Rate Limiting
	RateLimitRequestsPerSecond = 100
	RateLimitBurstSize         = 200
	RateLimitWhitelist         = "127.0.0.1,::1,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"

	// Security Tuning
	MaxRequestSize      = 10 << 20 // 10 MB
	ReadTimeout         = 15 * time.Second
	WriteTimeout        = 15 * time.Second
	IdleTimeout         = 60 * time.Second
	ReadHeaderTimeout   = 5 * time.Second
	ShutdownGracePeriod = 30 * time.Second

	// JWT
	JwtIssuer        = "doiit-platform"
	JwtAudience      = "doiit-users"
	JwtExpiry        = 24 * time.Hour
	JwtRefreshExpiry = 7 * 24 * time.Hour
)

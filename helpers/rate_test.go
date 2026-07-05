package helpers

import (
	"testing"

	"golang.org/x/time/rate"
)

func TestNewIPRateLimiterReusesLimiterForSameIP(t *testing.T) {
	t.Parallel()

	limiter := NewIPRateLimiter(rate.Limit(1), 2)
	first := limiter.GetLimiter("127.0.0.1")
	second := limiter.GetLimiter("127.0.0.1")

	if first != second {
		t.Fatalf("expected the same limiter instance for the same IP")
	}
}

func TestNewIPRateLimiterCreatesDistinctLimitersForDifferentIPs(t *testing.T) {
	t.Parallel()

	limiter := NewIPRateLimiter(rate.Limit(1), 2)
	first := limiter.GetLimiter("127.0.0.1")
	second := limiter.GetLimiter("192.168.1.10")

	if first == second {
		t.Fatalf("expected different limiter instances for different IPs")
	}
}

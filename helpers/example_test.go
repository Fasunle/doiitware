package helpers

func ExampleGenerateRequestID() {
	_ = GenerateRequestID()
}

func ExampleGenerateCorrelationID() {
	_ = GenerateCorrelationID()
}

func ExampleNewIPRateLimiter() {
	limiter := NewIPRateLimiter(10, 20)
	_ = limiter.GetLimiter("127.0.0.1")
}

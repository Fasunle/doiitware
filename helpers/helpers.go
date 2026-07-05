package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateCorrelationID returns a trace-friendly request correlation ID.
func GenerateCorrelationID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("req_%d_%s", time.Now().UnixNano(), hex.EncodeToString(b))
}

// GenerateRequestID returns a random request ID suitable for logging and response headers.
func GenerateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

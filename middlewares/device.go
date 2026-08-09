package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/fasunle/doiitware/config"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

// DeviceMiddleware extracts device information from request headers and client info
func DeviceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceInfo := extractDeviceInfo(c)

		// Store in context
		c.Set(config.DeviceContextKey, deviceInfo)
		c.Set(config.DeviceIDContextKey, deviceInfo.DeviceID)

		c.Next()
	}
}

// extractDeviceInfo extracts device information from the request
func extractDeviceInfo(c *gin.Context) config.DeviceInfo {
	// Try to get device ID from header first
	deviceID := c.GetHeader("X-Device-ID")

	// If not in header, generate one from client info
	if deviceID == "" {
		deviceID = generateDeviceID(c).String()
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	// Parse user agent for device details
	deviceName, deviceType, os, osVersion, browser, browserVer, platform := parseUserAgent(userAgent)

	return config.DeviceInfo{
		DeviceID:   uuid.FromStringOrNil(deviceID),
		DeviceName: deviceName,
		DeviceType: deviceType,
		OS:         os,
		OSVersion:  osVersion,
		Browser:    browser,
		BrowserVer: browserVer,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Platform:   platform,
	}
}

// generateDeviceID creates a deterministic device ID from client information
func generateDeviceID(c *gin.Context) uuid.UUID {
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	acceptLanguage := c.GetHeader("Accept-Language")

	// Combine client info to create a unique fingerprint
	fingerprint := ip + "|" + userAgent + "|" + acceptLanguage

	// Create SHA256 hash
	hash := sha256.Sum256([]byte(fingerprint))
	hashStr := hex.EncodeToString(hash[:])

	// Generate UUID from hash (using UUID v5)
	return uuid.NewV5(uuid.NamespaceDNS, hashStr)
}

// parseUserAgent extracts device information from User-Agent string
func parseUserAgent(userAgent string) (deviceName, deviceType, os, osVersion, browser, browserVer, platform string) {
	// Default values
	deviceType = "unknown"
	platform = "unknown"

	if userAgent == "" {
		return "unknown", deviceType, "unknown", "unknown", "unknown", "unknown", platform
	}

	// Detect OS
	switch {
	case strings.Contains(userAgent, "Windows NT"):
		os = "Windows"
		platform = "windows"
		// Extract Windows version
		re := regexp.MustCompile(`Windows NT (\d+\.\d+)`)
		if matches := re.FindStringSubmatch(userAgent); len(matches) > 1 {
			osVersion = matches[1]
			switch osVersion {
			case "10.0":
				osVersion = "10/11"
			case "6.3":
				osVersion = "8.1"
			case "6.2":
				osVersion = "8"
			case "6.1":
				osVersion = "7"
			}
		}

	case strings.Contains(userAgent, "Mac OS X"):
		os = "macOS"
		platform = "macos"
		re := regexp.MustCompile(`Mac OS X (\d+[._]\d+)`)
		if matches := re.FindStringSubmatch(userAgent); len(matches) > 1 {
			osVersion = strings.Replace(matches[1], "_", ".", -1)
		}

	case strings.Contains(userAgent, "iPhone"):
		os = "iOS"
		platform = "ios"
		deviceType = "mobile"
		deviceName = "iPhone"
		re := regexp.MustCompile(`iPhone OS (\d+[._]\d+)`)
		if matches := re.FindStringSubmatch(userAgent); len(matches) > 1 {
			osVersion = strings.Replace(matches[1], "_", ".", -1)
		}

	case strings.Contains(userAgent, "iPad"):
		os = "iOS"
		platform = "ios"
		deviceType = "tablet"
		deviceName = "iPad"
		re := regexp.MustCompile(`iPad OS (\d+[._]\d+)`)
		if matches := re.FindStringSubmatch(userAgent); len(matches) > 1 {
			osVersion = strings.Replace(matches[1], "_", ".", -1)
		}

	case strings.Contains(userAgent, "Android"):
		os = "Android"
		platform = "android"
		deviceType = "mobile"
		re := regexp.MustCompile(`Android (\d+\.\d+)`)
		if matches := re.FindStringSubmatch(userAgent); len(matches) > 1 {
			osVersion = matches[1]
		}
		// Try to get device model
		re = regexp.MustCompile(`; ([^;]+) Build/`)
		if matches := re.FindStringSubmatch(userAgent); len(matches) > 1 {
			deviceName = matches[1]
		}

	case strings.Contains(userAgent, "Linux"):
		os = "Linux"
		platform = "linux"
	}

	// Detect browser
	switch {
	case strings.Contains(userAgent, "Edg/"):
		browser = "Edge"
		re := regexp.MustCompile(`Edg/(\d+\.\d+\.\d+\.\d+)`)
		if matches := re.FindStringSubmatch(userAgent); len(matches) > 1 {
			browserVer = matches[1]
		}

	case strings.Contains(userAgent, "Firefox/"):
		browser = "Firefox"
		re := regexp.MustCompile(`Firefox/(\d+\.\d+)`)
		if matches := re.FindStringSubmatch(userAgent); len(matches) > 1 {
			browserVer = matches[1]
		}

	case strings.Contains(userAgent, "Chrome/"):
		browser = "Chrome"
		re := regexp.MustCompile(`Chrome/(\d+\.\d+\.\d+\.\d+)`)
		if matches := re.FindStringSubmatch(userAgent); len(matches) > 1 {
			browserVer = matches[1]
		}

	case strings.Contains(userAgent, "Safari/"):
		if !strings.Contains(userAgent, "Chrome/") {
			browser = "Safari"
			re := regexp.MustCompile(`Version/(\d+\.\d+)`)
			if matches := re.FindStringSubmatch(userAgent); len(matches) > 1 {
				browserVer = matches[1]
			}
		}

	case strings.Contains(userAgent, "Opera/"):
		browser = "Opera"
		re := regexp.MustCompile(`Opera/(\d+\.\d+)`)
		if matches := re.FindStringSubmatch(userAgent); len(matches) > 1 {
			browserVer = matches[1]
		}
	}

	// Determine device type if not already set
	if deviceType == "unknown" {
		if strings.Contains(userAgent, "Mobile") ||
			strings.Contains(userAgent, "Android") && !strings.Contains(userAgent, "Tablet") {
			deviceType = "mobile"
		} else if strings.Contains(userAgent, "Tablet") ||
			strings.Contains(userAgent, "iPad") {
			deviceType = "tablet"
		} else {
			deviceType = "desktop"
		}
	}

	// Set device name if still unknown
	if deviceName == "" {
		switch platform {
		case "windows":
			deviceName = "Windows PC"
		case "macos":
			deviceName = "Mac"
		case "linux":
			deviceName = "Linux PC"
		case "ios":
			if deviceType == "tablet" {
				deviceName = "iPad"
			} else {
				deviceName = "iPhone"
			}
		case "android":
			if deviceType == "tablet" {
				deviceName = "Android Tablet"
			} else {
				deviceName = "Android Phone"
			}
		default:
			deviceName = deviceType
		}
	}

	return
}

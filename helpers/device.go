package helpers

import (
	"github.com/fasunle/doiitware/config"
	"github.com/gin-gonic/gin"
)

// GetDeviceInfo retrieves the device information from the Gin context.
func GetDeviceInfo(c *gin.Context) (*config.DeviceInfo, bool) {
	deviceInfo, exists := c.Get(config.DeviceContextKey)
	if !exists {
		return nil, false
	}

	if di, ok := deviceInfo.(config.DeviceInfo); ok {
		return &di, true
	}

	return nil, false
}

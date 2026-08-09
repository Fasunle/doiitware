package config

import "github.com/gofrs/uuid/v5"

type DeviceInfo struct {
	DeviceID   uuid.UUID `json:"device_id"`
	DeviceName string    `json:"device_name"`
	DeviceType string    `json:"device_type"`
	OS         string    `json:"os"`
	OSVersion  string    `json:"os_version"`
	Browser    string    `json:"browser"`
	BrowserVer string    `json:"browser_version"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Platform   string    `json:"platform"`
}

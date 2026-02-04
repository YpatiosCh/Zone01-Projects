package session

import (
	uuid "forum/pkg/UUID"
	"time"
)

// GenerateSessionToken creates a new session token/ID
func GenerateSessionToken() string {
	return uuid.GenerateUUID()
}

// CalculateSessionExpiry calculates the expiry time for a session
// Default session lifetime is 24 hours
func CalculateSessionExpiry() time.Time {
	return time.Now().Add(24 * time.Hour)
}

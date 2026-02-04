package uuid

import (
	"github.com/google/uuid"
)

// GenerateUUIDv4 generates a new UUIDv4
// and returns it as a string.
func GenerateUUID() string {
	// Generate a new UUID
	uuid := uuid.New()
	return uuid.String()
}

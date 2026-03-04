package uuid

import "github.com/google/uuid"

// DefaultUUIDV7 returns a new UUID version 7 string using the system clock.
// It is the preferred method for sortable unique identifiers in our system.
func DefaultUUIDV7() uuid.UUID {
	key, err := uuid.NewV7()
	if err != nil {
		panic("critical: failed to generate UUID v7: " + err.Error())
	}

	return key
}

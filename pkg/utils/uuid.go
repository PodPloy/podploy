package utils

import "github.com/google/uuid"

// Default UUID V7 with error management
func DefaultUUIDV7() uuid.UUID {
	key, err := uuid.NewV7()
	if err != nil {
		panic("critical: failed to generate UUID v7: " + err.Error())
	}

	return key
}

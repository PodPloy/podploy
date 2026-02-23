package ports

type contextKey string

// Context keys used to store and retrieve request-scoped values.
const (
	// RequestIDKey is the context key for the unique request identifier.
	RequestIDKey contextKey = "request_id"
	// UserIDKey is the context key for the authenticated user identifier.
	UserIDKey contextKey = "user_id"
)

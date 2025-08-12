package dto

// ContextKey is a custom type for context keys to avoid collisions.
type ContextKey string

// UserIDKey is the key for userID in context.
const UserIDKey ContextKey = "userID"

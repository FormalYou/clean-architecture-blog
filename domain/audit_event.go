package domain

import "time"

// AuditEvent represents a business event that needs to be audited.
type AuditEvent struct {
	Timestamp time.Time              `json:"timestamp"`
	UserID    int64                  `json:"user_id"`
	Action    string                 `json:"action"`
	Entity    string                 `json:"entity"`
	EntityID  int64                 `json:"entity_id"`
	Details   map[string]interface{} `json:"details"`
}

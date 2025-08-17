package usecase

import (
	"encoding/json"

	"github.com/FormalYou/clean-architecture-blog/domain"
	"github.com/FormalYou/clean-architecture-blog/internal/application/contracts"
)

type auditService struct {
	logger contracts.Logger
}

// NewAuditService creates a new audit service.
func NewAuditService(logger contracts.Logger) contracts.AuditService {
	return &auditService{
		// Create a new logger instance with the "audit" type tag.
		logger: logger.With("log_type", "audit"),
	}
}

// RecordEvent records a business event for auditing purposes.
func (s *auditService) RecordEvent(event domain.AuditEvent) {
	// Convert the event details to a JSON string for structured logging.
	details, err := json.Marshal(event.Details)
	if err != nil {
		s.logger.Error("Failed to marshal audit event details", "error", err)
		return
	}

	s.logger.Info("Audit event recorded",
		"timestamp", event.Timestamp,
		"user_id", event.UserID,
		"action", event.Action,
		"entity", event.Entity,
		"entity_id", event.EntityID,
		"details", string(details),
	)
}

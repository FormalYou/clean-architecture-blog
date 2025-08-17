package contracts

import "github.com/FormalYou/clean-architecture-blog/domain"

// AuditService defines the contract for an audit service.
//
//go:generate mockgen -destination=./mocks/mock_audit_service.go -package=mocks . AuditService
type AuditService interface {
	// RecordEvent records a business event for auditing purposes.
	RecordEvent(event domain.AuditEvent)
}

package contracts

import "github.com/formal-you/clean-architecture-blog/domain"

//go:generate mockgen -destination=./mocks/mock_audit_service.go -package=mocks . AuditService
// AuditService defines the contract for an audit service.
type AuditService interface {
	// RecordEvent records a business event for auditing purposes.
	RecordEvent(event domain.AuditEvent)
}

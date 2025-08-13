package usecase

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/formal-you/clean-architecture-blog/domain"
	"github.com/formal-you/clean-architecture-blog/internal/application/contracts/mocks"
	"go.uber.org/mock/gomock"
)

func TestAuditService_RecordEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)
	auditLogger := mocks.NewMockLogger(ctrl)
	mockLogger.EXPECT().With("log_type", "audit").Return(auditLogger)

	auditService := NewAuditService(mockLogger)

	event := domain.AuditEvent{
		Timestamp: time.Now(),
		UserID:    123,
		Action:    "create",
		Entity:    "article",
		EntityID:  456,
		Details:   map[string]interface{}{"title": "New Article"},
	}

	details, _ := json.Marshal(event.Details)

	auditLogger.EXPECT().Info("Audit event recorded",
		"timestamp", gomock.Any(),
		"user_id", event.UserID,
		"action", event.Action,
		"entity", event.Entity,
		"entity_id", event.EntityID,
		"details", string(details),
	)

	auditService.RecordEvent(event)
}
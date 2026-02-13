package usersessions

import (
	"time"

	"github.com/google/uuid"
)

type AuditUserSessions struct {
	ID        uuid.UUID
	SessionID uuid.UUID
	UserID    uuid.UUID
	EventType string
	OldValue  string // "" if null
	NewValue  string // "" if null
	CreatedAt time.Time
}

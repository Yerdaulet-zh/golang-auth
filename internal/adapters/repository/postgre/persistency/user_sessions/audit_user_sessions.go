package usersessions

import (
	"time"

	persistency "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
	domain_sessions "github.com/golang-auth/internal/core/domain/user_sessions"
	"github.com/google/uuid"
)

/*
CREATE TABLE audit_user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL, -- The JTI
    user_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    old_value TEXT NULL,      -- Valid to be NULL
    new_value TEXT NULL,      -- Valid to be NULL
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
*/

type AuditUserSessions struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SessionID uuid.UUID `gorm:"type:uuid;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	EventType string    `gorm:"type:audit_event_type;not null"`
	OldValue  *string   `gorm:"type:text"`
	NewValue  *string   `gorm:"type:text"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:now();not null"`

	User persistency.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func MapAuditToDomain(orm AuditUserSessions) domain_sessions.AuditUserSessions {
	oldVal := ""
	if orm.OldValue != nil {
		oldVal = *orm.OldValue
	}

	newVal := ""
	if orm.NewValue != nil {
		newVal = *orm.NewValue
	}

	return domain_sessions.AuditUserSessions{
		ID:        orm.ID,
		SessionID: orm.SessionID,
		UserID:    orm.UserID,
		EventType: orm.EventType,
		OldValue:  oldVal,
		NewValue:  newVal,
		CreatedAt: orm.CreatedAt,
	}
}

func MapAuditToORM(d domain_sessions.AuditUserSessions) AuditUserSessions {
	// In Audit logs, we often want to store empty strings as NULL in DB
	var oldPtr *string
	if d.OldValue != "" {
		oldPtr = &d.OldValue
	}

	var newPtr *string
	if d.NewValue != "" {
		newPtr = &d.NewValue
	}

	return AuditUserSessions{
		ID:        d.ID,
		SessionID: d.SessionID,
		UserID:    d.UserID,
		EventType: d.EventType,
		OldValue:  oldPtr,
		NewValue:  newPtr,
		CreatedAt: d.CreatedAt,
	}
}

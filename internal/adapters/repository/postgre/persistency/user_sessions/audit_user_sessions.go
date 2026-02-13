package usersessions

import (
	"time"

	"github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
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

	User user.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

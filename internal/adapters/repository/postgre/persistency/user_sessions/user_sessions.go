package repousersessions

import (
	"time"

	repouser "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
	domain_sessions "github.com/golang-auth/internal/core/domain/user_sessions"
	"github.com/google/uuid"
)

type UserSessions struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	IPAddress  string    `gorm:"type:inet;not null"`
	UserAgent  string    `gorm:"type:text;not null"`
	Device     *string   `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"type:timestamptz;default:now();not null"`
	LastActive time.Time `gorm:"type:timestamptz;default:now();not null"`
	ExpiresAt  time.Time `gorm:"type:timestamptz;not null"`

	User repouser.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func MapSessionToDomain(orm UserSessions) domain_sessions.UserSessions {
	device := ""
	if orm.Device != nil {
		device = *orm.Device
	}

	return domain_sessions.UserSessions{
		ID:         orm.ID,
		UserID:     orm.UserID,
		IPAddress:  orm.IPAddress,
		UserAgent:  orm.UserAgent,
		Device:     device,
		CreatedAt:  orm.CreatedAt,
		LastActive: orm.LastActive,
		ExpiresAt:  orm.ExpiresAt,
	}
}

/*
CREATE TABLE user_sessions (
    -- The JTI (JWT ID) serves as the Primary Key
    id UUID PRIMARY KEY,

    -- Link to your users table
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- INET handles both IPv4 and IPv6 efficiently
    ip_address INET NOT NULL,

    -- The full raw string for forensics
    user_agent TEXT NOT NULL,

    -- Human-friendly version for the UI (e.g., "Chrome on MacOS")
    device_name TEXT,

    -- When the session was first created
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- When the user last refreshed their token
    last_active TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Hard expiry (matching your Refresh Token TTL, e.g., 8-24 hours)
    expires_at TIMESTAMPTZ NOT NULL
);

-- Index for fast lookups when a user wants to see their "Active Devices"
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);

-- Optional: Index for cleanup tasks (deleting expired sessions)
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);
*/

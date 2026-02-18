package ports

import (
	"context"
	"time"

	repouser "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
	userverification "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user_verification"
	"github.com/google/uuid"
)

type UserRepoPorts interface {
	GetUserByEmail(ctx context.Context, email string) (*repouser.User, error)
	CreateUserWithCredentials(ctx context.Context, req UserAndCredentialsRequest) error
	GetVerificationByToken(ctx context.Context, token string) (*userverification.UserVerification, error)
	ConfirmVerification(ctx context.Context, userID uuid.UUID, verificationID uuid.UUID) error
	GetVerificationByUserID(ctx context.Context, userID uuid.UUID) (*userverification.UserVerification, error)
	RotateVerificationToken(ctx context.Context, recordID uuid.UUID, status string, req *userverification.UserVerification) error
	GetCountsOfVerificationRecordsByUserID(ctx context.Context, user_id uuid.UUID, timeDuration time.Time) (int64, error)
	UpdateUserVerificationTokenStatus(ctx context.Context, tokenID uuid.UUID, status string) error

	// Login
	// CreateUserSession(ctx context.Context, session *repousersessions.UserSessions) error
	GetUserSessionCountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	CreateAndDeleteOldUserSession(ctx context.Context, sessionReq *CreateUserSessionRequest) (*CreateUserSessionResponse, error)
	CreateUserSession(ctx context.Context, sessionReq *CreateUserSessionRequest) (*CreateUserSessionResponse, error)
	DeleteUserSession(ctx context.Context, session_id uuid.UUID) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

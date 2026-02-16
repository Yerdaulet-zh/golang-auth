package service

import (
	"context"
	"errors"
	"testing"

	repouser "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
	userverification "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user_verification"
	"github.com/golang-auth/internal/core/domain"
	"github.com/golang-auth/internal/core/ports"
	"github.com/golang-auth/internal/testutil"
	"github.com/google/uuid"
)

// Mock Implementation
type mockUserRepo struct {
	getUserByEmailFn       func(ctx context.Context, email string) (*repouser.User, error)
	createUserFn           func(ctx context.Context, req ports.UserAndCredentialsRequest) error
	verifyUserEmail        func(ctx context.Context, token string) error
	getVerificationByToken func(ctx context.Context, token string) (*userverification.UserVerification, error)
	confirmVerification    func(ctx context.Context, userID uuid.UUID, verificationID uuid.UUID) error
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*repouser.User, error) {
	return m.getUserByEmailFn(ctx, email)
}

func (m *mockUserRepo) CreateUserWithCredentials(ctx context.Context, req ports.UserAndCredentialsRequest) error {
	return m.createUserFn(ctx, req)
}

func (m *mockUserRepo) VerifyUserEmail(ctx context.Context, token string) error {
	return m.verifyUserEmail(ctx, token)
}

func (m *mockUserRepo) ConfirmVerification(ctx context.Context, userID uuid.UUID, verificationID uuid.UUID) error {
	return m.confirmVerification(ctx, userID, verificationID)
}

func (m *mockUserRepo) GetVerificationByToken(ctx context.Context, token string) (*userverification.UserVerification, error) {
	return m.getVerificationByToken(ctx, token)
}

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		setupMock   func(m *mockUserRepo)
		wantErr     bool
		expectedErr error
	}{
		{
			name:     "Invalid email format",
			email:    "invalid-email",
			password: "password123",
			setupMock: func(m *mockUserRepo) {
				// Logic shouldn't even reach the repo if validation fails first
			},
			wantErr:     true,
			expectedErr: domain.ErrInvalidEmail,
		},
		{
			name:     "Error while hashing password",
			email:    "test@gmail.com",
			password: string(make([]byte, 100)),
			setupMock: func(m *mockUserRepo) {
				m.getUserByEmailFn = func(ctx context.Context, email string) (*repouser.User, error) {
					return nil, errors.New("Critical Error, hashing algorithm did not react properly")
				}
			},
			wantErr:     true,
			expectedErr: domain.ErrHashingError,
		},
		{
			name:     "Success registration",
			email:    "newuser@gmail.com",
			password: "password123",
			setupMock: func(m *mockUserRepo) {
				// Simulate: User not found in DB
				m.getUserByEmailFn = func(ctx context.Context, email string) (*repouser.User, error) {
					return nil, domain.ErrNotFound
				}
				// Simulate: Successful creation
				m.createUserFn = func(ctx context.Context, req ports.UserAndCredentialsRequest) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:     "User already exists",
			email:    "exists@gmail.com",
			password: "password123",
			setupMock: func(m *mockUserRepo) {
				// Simulate: User already exists in DB
				m.getUserByEmailFn = func(ctx context.Context, email string) (*repouser.User, error) {
					return &repouser.User{Email: email}, nil
				}
			},
			wantErr:     true,
			expectedErr: domain.ErrUserAlreadyExists,
		},
		{
			name:     "Database failure on user lookup",
			email:    "dbfail@gmail.com",
			password: "password123",
			setupMock: func(m *mockUserRepo) {
				m.getUserByEmailFn = func(ctx context.Context, email string) (*repouser.User, error) {
					return nil, errors.New("connection reset by peer")
				}
			},
			wantErr:     true,
			expectedErr: domain.ErrDatabaseInternalError,
		},
		{
			name:     "Repository failure during creation",
			email:    "newuser@gmail.com",
			password: "password123",
			setupMock: func(m *mockUserRepo) {
				m.getUserByEmailFn = func(ctx context.Context, email string) (*repouser.User, error) {
					return nil, domain.ErrNotFound
				}
				m.createUserFn = func(ctx context.Context, req ports.UserAndCredentialsRequest) error {
					return domain.ErrDatabaseInternalError
				}
			},
			wantErr:     true,
			expectedErr: domain.ErrDatabaseInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock and service
			mockRepo := &mockUserRepo{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := NewUserService(mockRepo, &testutil.NoopLogger{}, &testutil.NoPublisher{})

			// Execute
			err := svc.Register(context.Background(), tt.email, tt.password)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Fatalf("Register() unexpected error status: %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, tt.expectedErr) {
				t.Errorf("Register() got = %v, want %v", err, tt.expectedErr)
			}
		})
	}
}

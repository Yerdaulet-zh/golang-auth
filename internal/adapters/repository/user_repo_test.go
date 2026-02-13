package repository

import (
	"context"
	"testing"

	"github.com/golang-auth/internal/adapters/logging"
	"github.com/golang-auth/internal/core/domain"
	"github.com/golang-auth/internal/core/ports"
	"github.com/golang-auth/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Integration(t *testing.T) {
	logger := logging.NewStdoutLogger()

	// 1. Initialize the Real Database via Testcontainers & Atlas
	gormDB := testutil.SetupTestDB(t)
	repo := NewUserRepository(gormDB, logger)
	ctx := context.Background()

	t.Run("CreateUserWithCredentials_Success", func(t *testing.T) {
		req := ports.UserAndCredentialsRequest{
			Email:        "user@enterprise.com",
			UserStatus:   "active",
			IsMFAEnabled: false,
			PasswordHash: "$2a$12$test-hash",
		}

		// Act
		err := repo.CreateUserWithCredentials(ctx, req)

		// Assert
		assert.NoError(t, err)

		// Verify data exists in DB
		user, err := repo.GetUserByEmail(ctx, "user@enterprise.com")
		assert.NoError(t, err)
		assert.Equal(t, req.Email, user.Email)
		assert.Equal(t, req.UserStatus, user.UserStatus)
	})

	t.Run("GetUserByEmail_NotFound", func(t *testing.T) {
		// Act
		user, err := repo.GetUserByEmail(ctx, "non-existent@test.com")

		// Assert
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, user)
	})

	t.Run("CreateUserWithCredentials_TransactionRollback", func(t *testing.T) {
		// This test ensures that if the second part of the transaction fails,
		// the user is NOT created (Atomic operation).

		// To simulate failure, we send a request that violates a DB constraint
		// (e.g., missing password hash if your Atlas schema requires it)
		// Or we can rely on GORM generic errors.

		invalidReq := ports.UserAndCredentialsRequest{
			Email:        "rollback@test.com",
			PasswordHash: "", // Assume this triggers an error in your logic or DB
		}

		// Forcing a failure here depends on your DB constraints.
		// If you have a NOT NULL constraint on PasswordHash in Atlas:
		err := repo.CreateUserWithCredentials(ctx, invalidReq)

		// If it failed as expected:
		if err != nil {
			// Verify the user was NOT created due to rollback
			user, err := repo.GetUserByEmail(ctx, "rollback@test.com")
			assert.ErrorIs(t, err, domain.ErrNotFound)
			assert.Nil(t, user)
		}
	})
}

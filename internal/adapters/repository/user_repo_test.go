package repository

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang-auth/internal/core/domain"
	"github.com/golang-auth/internal/core/ports"
	"github.com/golang-auth/internal/testutil"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	db, container, err := testutil.SetupGlobalTestDB()
	if err != nil {
		panic("Failed to setup repository test DB: " + err.Error())
	}
	testDB = db

	code := m.Run()

	_ = container.Terminate(context.Background())
	os.Exit(code)
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	testutil.TruncateAllTables(testDB)
	repo := NewUserRepository(testDB, &testutil.NoopLogger{})
	ctx := context.Background()

	t.Run("Should return ErrNotFound when user doesn't exist", func(t *testing.T) {
		user, err := repo.GetUserByEmail(ctx, "nonexistent@test.com")
		if user != nil {
			t.Error("Expected nil user")
		}
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Expected domain.ErrNotFound, got %v", err)
		}
	})

	t.Run("Should return user when user exists", func(t *testing.T) {
		email := "exists@test.com"
		repo.CreateUserWithCredentials(ctx, ports.UserAndCredentialsRequest{
			Email:        email,
			PasswordHash: "hashed_pass",
		})

		user, err := repo.GetUserByEmail(ctx, email)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if user.Email != email {
			t.Errorf("Expected email %s, got %s", email, user.Email)
		}
	})
}

func TestUserRepository_CreateUserWithCredentials(t *testing.T) {
	repo := NewUserRepository(testDB, &testutil.NoopLogger{})
	ctx := context.Background()

	t.Run("Should successfully create user and credentials", func(t *testing.T) {
		testutil.TruncateAllTables(testDB)
		req := ports.UserAndCredentialsRequest{
			Email:        "new@test.com",
			PasswordHash: "secret_hash",
		}

		err := repo.CreateUserWithCredentials(ctx, req)
		if err != nil {
			t.Fatalf("Failed to create: %v", err)
		}

		var userCount, credCount int64
		testDB.Table("user").Count(&userCount)
		testDB.Table("user_credentials").Count(&credCount)

		if userCount != 1 || credCount != 1 {
			t.Errorf("Expected 1 user and 1 cred, got %d and %d", userCount, credCount)
		}
	})

	t.Run("Should rollback transaction if second insert fails", func(t *testing.T) {
		testutil.TruncateAllTables(testDB)

		// We trigger an error by creating a credentials record
		// that violates a constraint or by manually injecting a failure.
		// For this test, let's use a very long email that exceeds DB column limits
		// if applicable, OR simply rely on the fact that 'tx' fix now enables rollback.

		// A reliable way: Attempt to create a user with an email that already exists
		// (if the first insert succeeds but second fails for another reason)
		// Since we want to test the TRANSACTION, let's pass an invalid UserID
		// to the second part. To do this, we'd need to mock the tx, but in integration:

		// Use a password hash that is way too long for a VARCHAR(255) to trigger DB error
		longHash := string(make([]byte, 5000))

		req := ports.UserAndCredentialsRequest{
			Email:        "rollback-me@test.com",
			PasswordHash: longHash,
		}

		err := repo.CreateUserWithCredentials(ctx, req)
		if err == nil {
			t.Error("Expected error due to long password hash, but got nil")
		}

		var count int64
		testDB.Table("user").Where("email = ?", "rollback-me@test.com").Count(&count)
		if count != 0 {
			t.Error("Transaction failed to rollback: User record exists even though credentials failed")
		}
	})
}

package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang-auth/internal/adapters/repository"
	"github.com/golang-auth/internal/core/service"
	"github.com/golang-auth/internal/testutil"
	"gorm.io/gorm"
)

var globalDB *gorm.DB

func TestMain(m *testing.M) {
	// 1. Setup the container once for the whole package
	db, container, err := testutil.SetupGlobalTestDB()
	if err != nil {
		panic("Failed to setup integration test DB: " + err.Error())
	}
	globalDB = db

	// 2. Run all tests in this file
	code := m.Run()

	// 3. Cleanup after all tests are done
	_ = container.Terminate(context.Background())
	os.Exit(code)
}

func TestUserHandler_Register_Integration(t *testing.T) {
	// PREPARATION: Reset DB state for this specific test
	if err := testutil.TruncateAllTables(globalDB); err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}

	// Setup real layers with the global DB
	repo := repository.NewUserRepository(globalDB, &testutil.NoopLogger{})
	svc := service.NewUserService(repo, &testutil.NoopLogger{}, &testutil.NoPublisher{})
	handler := NewUserHandler(svc, &testutil.NoopLogger{})

	t.Run("Integration: Successful Registration and Duplicate Check", func(t *testing.T) {
		payload := `{"email": "integ_test@gmail.com", "password": "password123"}`

		// --- Request 1: Should succeed ---
		req1 := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(payload))
		w1 := httptest.NewRecorder()
		handler.Register(w1, req1)

		if w1.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d. Body: %s", w1.Code, w1.Body.String())
		}

		// --- Request 2: Should conflict (Same Email) ---
		req2 := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(payload))
		w2 := httptest.NewRecorder()
		handler.Register(w2, req2)

		if w2.Code != http.StatusConflict {
			t.Errorf("Expected 409 Conflict, got %d", w2.Code)
		}
	})
}

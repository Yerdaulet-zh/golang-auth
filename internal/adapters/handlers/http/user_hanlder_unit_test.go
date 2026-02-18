package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-auth/internal/core/domain"
	"github.com/golang-auth/internal/core/ports"
	"github.com/golang-auth/internal/testutil"
)

// Mock Service
type mockUserService struct {
	registerFn                   func(ctx context.Context, email, password string) error
	verifyUserEmail              func(ctx context.Context, token string) error
	resendEmailVerificationToken func(ctx context.Context, email string) error
	login                        func(ctx context.Context, req *ports.LoginRequest) (*ports.LoginResponse, error)
}

func (m *mockUserService) Register(ctx context.Context, email, password string) error {
	return m.registerFn(ctx, email, password)
}

func (m *mockUserService) VerifyUserEmail(ctx context.Context, token string) error {
	return m.verifyUserEmail(ctx, token)
}

func (m *mockUserService) ResendEmailVerificationToken(ctx context.Context, email string) error {
	return m.resendEmailVerificationToken(ctx, email)
}

func (m *mockUserService) Login(ctx context.Context, req *ports.LoginRequest) (*ports.LoginResponse, error) {
	return m.login(ctx, req)
}

func TestUserHandler_Register_Unit(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{}
		mockReturn     error
		expectedStatus int
	}{
		{
			name: "Successful Registration",
			payload: map[string]string{
				"email":    "test@gmail.com",
				"password": "password123",
			},
			mockReturn:     nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Conflict - User Already Exists",
			payload: map[string]string{
				"email":    "exists@gmail.com",
				"password": "password123",
			},
			mockReturn:     domain.ErrUserAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Invalid JSON Payload",
			payload:        "not-a-json",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Setup Mock
			mockSvc := &mockUserService{
				registerFn: func(ctx context.Context, email, password string) error {
					return tt.mockReturn
				},
			}
			handler := NewUserHandler(mockSvc, &testutil.NoopLogger{})

			// 2. Prepare Request
			var body []byte
			if s, ok := tt.payload.(string); ok {
				body = []byte(s)
			} else {
				body, _ = json.Marshal(tt.payload)
			}

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			// 3. Execute
			handler.Register(w, req)

			// 4. Assert
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

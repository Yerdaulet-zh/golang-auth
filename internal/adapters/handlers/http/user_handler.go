package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-auth/internal/core/domain"
	"github.com/golang-auth/internal/core/ports"
)

var validate = validator.New()

type UserHandler struct {
	userService ports.UserUseCase
	logger      ports.Logger
}

func NewUserHandler(service ports.UserUseCase, logger ports.Logger) *UserHandler {
	return &UserHandler{
		userService: service,
		logger:      logger,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}
	if err := validate.Struct(req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Missing or invalid required fields "+err.Error())
		return
	}

	ctx := r.Context()
	if err := h.userService.Register(ctx, req.Email, req.Password); err != nil {
		h.mapErrorToResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created, verify email and try to login"})
}

func (h *UserHandler) VerifyUserEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Verification token is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := h.userService.VerifyUserEmail(ctx, token)
	if err != nil {
		h.mapErrorToResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user account verified successfully",
	})
}

func (h *UserHandler) writeJSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if message != "" {
		json.NewEncoder(w).Encode(map[string]string{"error": message})
	}
}

func (h *UserHandler) mapErrorToResponse(w http.ResponseWriter, err error) {
	switch {
	// 409 Conflict
	case errors.Is(err, domain.ErrUserAlreadyExists):
		h.writeJSONError(w, http.StatusConflict, "User already registered")

	// 404 Not Found
	case errors.Is(err, domain.ErrNotFound):
		h.writeJSONError(w, http.StatusNotFound, "Resource not found")

	// 500 Internal Server Errors
	case errors.Is(err, domain.ErrDatabaseInternalError):
		h.writeJSONError(w, http.StatusInternalServerError, "")
	case errors.Is(err, domain.ErrRepositoryInternalError):
		h.writeJSONError(w, http.StatusInternalServerError, "")
	case errors.Is(err, domain.ErrHashingError):
		h.writeJSONError(w, http.StatusInternalServerError, "")
	case errors.Is(err, domain.ErrDomainInternalError):
		h.writeJSONError(w, http.StatusInternalServerError, "")
	case errors.Is(err, domain.ErrInvalidTokenState):
		h.writeJSONError(w, http.StatusInternalServerError, "")
	case errors.Is(err, domain.ErrBrokerInternalError):
		h.writeJSONError(w, http.StatusInternalServerError, "")

	// 400 Bad Request Errors
	case errors.Is(err, domain.ErrInvalidEmail):
		h.writeJSONError(w, http.StatusBadRequest, "Invalid Email")
	case errors.Is(err, domain.ErrTokenExpired):
		h.writeJSONError(w, http.StatusBadRequest, "")
	case errors.Is(err, domain.ErrTokenNotFound):
		h.writeJSONError(w, http.StatusInternalServerError, "")

	// 500 Internal Server Error (The Default)
	default:
		// We log the real error here for debugging
		h.logger.Error("Unhandled error", "error", err)
		h.writeJSONError(w, http.StatusInternalServerError, "Internal server error")
	}
}

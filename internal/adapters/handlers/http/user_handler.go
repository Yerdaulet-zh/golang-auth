package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-auth/internal/core/domain"
	"github.com/golang-auth/internal/core/ports"
	"github.com/google/uuid"
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

func (h *UserHandler) ResendVerificationToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ResendVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validate.Struct(req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Valid email is required")
		return
	}

	ctx := r.Context()
	if err := h.userService.ResendEmailVerificationToken(ctx, req.Email); err != nil {
		h.mapErrorToResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If the account exists and is unverified, a new verification email has been sent.",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	ctx := r.Context()
	loginRequest := ports.LoginRequest{
		Email:     req.Email,
		Password:  req.Password,
		IPAddress: "0.0.0.0",
		UserAgent: "Agent",
		Device:    "iPad",
	}
	res, err := h.userService.Login(ctx, &loginRequest)
	if err != nil {
		h.mapErrorToResponse(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "UserID",
		Value:    res.UserID.String(),
		Expires:  res.RefreshTokenExpiresAt,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "SessionID",
		Value:    res.SessionID.String(),
		Expires:  res.AccessTokenExpiresAt,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	// Set Access Token Cookie (Short-lived)
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		Expires:  res.AccessTokenExpiresAt,
		HttpOnly: true,
		Secure:   true, // Set to false only if developing on localhost without HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	// Set Refresh Token Cookie (Long-lived)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Expires:  res.RefreshTokenExpiresAt,
		HttpOnly: true,
		Secure:   true,
		Path:     "/auth/refresh", // Optional: Only send this cookie to the refresh endpoint
		SameSite: http.SameSiteLaxMode,
	})

	// Return success response (usually excluding the tokens from the JSON body for security)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
	})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cookie, err := r.Cookie("SessionID")
	if err != nil {
		h.writeJSONError(w, http.StatusUnauthorized, "No active session found")
		return
	}
	sessionID, err := uuid.Parse(cookie.Value)
	if err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Invalid session format")
		return
	}

	ctx := r.Context()
	err = h.userService.Logout(ctx, sessionID)
	if err != nil {
		h.logger.Debug("error SNIOAB", err)
		h.mapErrorToResponse(w, err)
		return
	}

	h.clearCookie(w, "UserID", "/")
	h.clearCookie(w, "SessionID", "/")
	h.clearCookie(w, "access_token", "/")
	h.clearCookie(w, "refresh_token", "/auth/refresh")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully logged out",
	})
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cookieUserID, err := r.Cookie("UserID")
	if err != nil {
		h.writeJSONError(w, http.StatusUnauthorized, "No user ID found from cookie")
		return
	}
	userID, err := uuid.Parse(cookieUserID.Value)
	if err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Invalid user id format")
		return
	}

	ctx := r.Context()
	if err := h.userService.DeleteAccount(ctx, userID); err != nil {
		h.mapErrorToResponse(w, err)
		return
	}

	h.clearCookie(w, "UserID", "/")
	h.clearCookie(w, "SessionID", "/")
	h.clearCookie(w, "access_token", "/")
	h.clearCookie(w, "refresh_token", "/auth/refresh")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Account and all associated sessions have been successfully deleted",
	})
}

// Helper to clear cookies
func (h *UserHandler) clearCookie(w http.ResponseWriter, name, path string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     path,
		MaxAge:   -1, // Tells browser to delete immediately
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
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
		h.writeJSONError(w, http.StatusConflict, domain.ErrUserAlreadyExists.Error())
	case errors.Is(err, domain.ErrUserAlreadyVerified):
		h.writeJSONError(w, http.StatusConflict, domain.ErrUserAlreadyVerified.Error())

	// 403 Forbidden
	case errors.Is(err, domain.ErrUserNotVerified):
		h.writeJSONError(w, http.StatusForbidden, domain.ErrUserNotVerified.Error())
	case errors.Is(err, domain.ErrUserAccountBanned):
		h.writeJSONError(w, http.StatusForbidden, domain.ErrUserAccountBanned.Error())
	case errors.Is(err, domain.ErrUserAccountSuspended):
		h.writeJSONError(w, http.StatusForbidden, domain.ErrUserAccountSuspended.Error())
	case errors.Is(err, domain.ErrTooManyUserSessions):
		h.writeJSONError(w, http.StatusForbidden, domain.ErrTooManyUserSessions.Error())

	// 404 Not Found
	case errors.Is(err, domain.ErrNotFound):
		h.writeJSONError(w, http.StatusNotFound, domain.ErrNotFound.Error())

	// 500 Internal Server Errors
	case errors.Is(err, domain.ErrDatabaseInternalError):
		h.writeJSONError(w, http.StatusInternalServerError, domain.ErrDatabaseInternalError.Error())
	case errors.Is(err, domain.ErrRepositoryInternalError):
		h.writeJSONError(w, http.StatusInternalServerError, domain.ErrRepositoryInternalError.Error())
	case errors.Is(err, domain.ErrHashingError):
		h.writeJSONError(w, http.StatusInternalServerError, domain.ErrHashingError.Error())
	case errors.Is(err, domain.ErrDomainInternalError):
		h.writeJSONError(w, http.StatusInternalServerError, domain.ErrDomainInternalError.Error())
	case errors.Is(err, domain.ErrInvalidTokenState):
		h.writeJSONError(w, http.StatusInternalServerError, domain.ErrInvalidTokenState.Error())
	case errors.Is(err, domain.ErrBrokerInternalError):
		h.writeJSONError(w, http.StatusInternalServerError, domain.ErrBrokerInternalError.Error())

	// 400 Bad Request Errors
	case errors.Is(err, domain.ErrInvalidEmail):
		h.writeJSONError(w, http.StatusBadRequest, domain.ErrInvalidEmail.Error())
	case errors.Is(err, domain.ErrTokenExpired):
		h.writeJSONError(w, http.StatusBadRequest, domain.ErrTokenExpired.Error())
	case errors.Is(err, domain.ErrTokenNotFound):
		h.writeJSONError(w, http.StatusInternalServerError, domain.ErrTokenNotFound.Error())
	case errors.Is(err, domain.ErrUserNotFound):
		h.writeJSONError(w, http.StatusBadRequest, domain.ErrUserNotFound.Error())
	case errors.Is(err, domain.ErrUsedToken):
		h.writeJSONError(w, http.StatusBadRequest, domain.ErrUsedToken.Error())
	case errors.Is(err, domain.ErrInvaidPassword):
		h.writeJSONError(w, http.StatusBadRequest, domain.ErrInvaidPassword.Error())
	// case errors.Is(err, domain.ErrSessionNotFound):
	// 	h.writeJSONError(w, http.StatusBadRequest, domain.ErrSessionNotFound.Error())

	// 429
	case errors.Is(err, domain.ErrTooManyRequests):
		h.writeJSONError(w, http.StatusTooManyRequests, domain.ErrTooManyRequests.Error())

	// 500 Internal Server Error (The Default)
	default:
		h.logger.Error("Unhandled error", "error", err)
		h.writeJSONError(w, http.StatusInternalServerError, "Internal server error")
	}
}

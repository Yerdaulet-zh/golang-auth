package http

import "github.com/google/uuid"

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=62"`
}

type CreateUserResponse struct {
	UserID       uuid.UUID
	SessionID    uuid.UUID
	RefreshToken string
	JWTToken     string
	JTI          uuid.UUID
}

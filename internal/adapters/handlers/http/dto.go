package http

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=62"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=62"`
}

// type DeleteAccountRequest struct {
// 	Email string `json:"email" validate:"required,email"`
// }

// type CreateUserResponse struct {
// 	UserID       uuid.UUID
// 	SessionID    uuid.UUID
// 	RefreshToken string
// 	JWTToken     string
// 	JTI          uuid.UUID
// }

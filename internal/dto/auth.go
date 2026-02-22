package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" validate:"len=6,numeric"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" validate:"len=6,numeric"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

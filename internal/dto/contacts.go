package dto

import "time"

type CreateContactRequest struct {
	Name        string  `json:"name" validate:"required"`
	PhoneNumber string  `json:"phone_number" validate:"required,e164"`
	ImageURI    *string `json:"image_uri"`
	Type        *string `json:"type" validate:"omitempty,oneof=Trusted Casual"`
	IsPinned    *bool   `json:"is_pinned"`
}

type UpdateContactRequest struct {
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,e164"`
	ImageURI    *string `json:"image_uri"`
	Type        *string `json:"type" validate:"omitempty,oneof=Trusted Casual"`
	IsPinned    *bool   `json:"is_pinned"`
}

type ContactResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	ImageURI    *string   `json:"image_uri"`
	Type        string    `json:"type"`
	IsPinned    bool      `json:"is_pinned"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

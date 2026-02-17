package dto

type UpdateProfileRequest struct {
	Name            *string `json:"name"`
	PhoneNumber     *string `json:"phone_number"`
	ProfileImageURL *string `json:"profile_image_url"`
	BloodGroup      *string `json:"blood_group"`
	Allergies       *string `json:"allergies"`
	Medications     *string `json:"medications"`
}
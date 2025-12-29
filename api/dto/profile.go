package dto

type ChangePasswordRequest struct {
	OldPassword string `form:"oldPassword" binding:"required"`
	NewPassword string `form:"newPassword" binding:"required,min=6"`
}

type ProfileResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

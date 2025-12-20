package domain

import "context"

type Profile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	OldPassword string `form:"oldPassword" binding:"required"`
	NewPassword string `form:"newPassword" binding:"required,min=6"`
}

type ProfileUsecase interface {
	GetProfileByID(c context.Context, userID string) (*Profile, error)
	ChangePassword(c context.Context, userID string, oldPassword string, newPassword string) error
}

package dtos

import (
	"time"
)

type UserRegisterInput struct {
	FullName        string  `form:"full_name" json:"full_name" example:"Nadiah Indrianti"`
	Email           string  `form:"email" json:"email" example:"user@nadiah.com"`
	Password        string  `form:"password" json:"password" example:"Nadiah271,."`
	ConfirmPassword string  `form:"confirm_password" json:"confirm_password" example:"Nadiah271,."`
	NIMNIP          string  `form:"nim_nip" json:"nim_nip" example:"2009189"`
	Role            string  `form:"role" json:"role" example:"user"`
	IsActive        *bool   `form:"is_active" json:"is_active,omitempty" example:"true"`
}

type UserRegisterInputByAdmin struct {
	FullName        string  `form:"full_name" json:"full_name" example:"Nadiah Indrianti"`
	Email           string  `form:"email" json:"email" example:"user@nadiah.com"`
	Password        string  `form:"password" json:"password" example:"Nadiah271,."`
	ConfirmPassword string  `form:"confirm_password" json:"confirm_password" example:"Nadiah271,."`
	NIMNIP          string  `form:"nim_nip" json:"nim_nip" example:"2009189"`
	Role            string  `form:"role" json:"role" example:"user"`
}

type UserRegisterInputUpdateByAdmin struct {
	FullName        string  `form:"full_name" json:"full_name" example:"Nadiah Indrianti"`
	Email           string  `form:"email" json:"email" example:"user@nadiah.com"`
	NIMNIP          string  `form:"nim_nip" json:"nim_nip" example:"2009189"`
	Role            string  `form:"role" json:"role" example:"user"`
	IsActive        *bool   `form:"is_active" json:"is_active,omitempty" example:"true"`
}

type UserLoginInput struct {
	Email           string  `form:"email" json:"email" example:"user@nadiah.com"`
	Password        string  `form:"password" json:"password" example:"Nadiah271,."`
}

type UserUpdatePhotoProfileInput struct {
	ProfilePicture string `form:"file" json:"file" example:"https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg"`
}

type UserUpdatePasswordInput struct {
	OldPassword     string `form:"old_password" json:"old_password" example:"qweqwe123"`
	NewPassword     string `form:"new_password" json:"new_password" example:"asdqwe123"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" example:"asdqwe123"`
}

type UserUpdateProfileInput struct {
	FullName    string  `form:"full_name" json:"full_name" example:"Nadiah Indrianti"`
	Email       string  `form:"email" json:"email" example:"user@nadiah.com"`
	NIMNIP      string  `form:"nim_nip" json:"nim_nip" example:"2009189"`
}

type UserLoginResponse struct {
	FullName    string  `form:"full_name" json:"full_name" example:"Nadiah Indrianti"`
	Email       string  `form:"email" json:"email" example:"user@nadiah.com"`
	NIMNIP      string  `form:"nim_nip" json:"nim_nip" example:"2009189"`
	Role        string    `json:"role" example:"user"`
	CreatedAt   time.Time `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}

type UserRegisterResponse struct {
	FullName    string  `form:"full_name" json:"full_name" example:"Nadiah Indrianti"`
	Email       string  `form:"email" json:"email" example:"user@nadiah.com"`
	NIMNIP      string  `form:"nim_nip" json:"nim_nip" example:"2009189"`
	Role        string    `json:"role" example:"user"`
	CreatedAt   time.Time `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}

type UserInformationResponse struct {
	ID         uint    `form:"user_id" json:"user_id"`
	FullName       string  `form:"full_name" json:"full_name" example:"Nadiah Indrianti"`
	Email          string  `form:"email" json:"email" example:"user@nadiah.com"`
	NIMNIP         string  `form:"nim_nip" json:"nim_nip" example:"2009189"`
	ProfilePicture string    `json:"profile_picture_url" example:"https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg"`
	Role           *string   `json:"role,omitempty" example:"user"`
	Token          *string   `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2ODQ0MDYzMzMsInJvbGUiOiJ1c2VyIiwidXNlcklkIjozfQ.B8vBlMIiU4iZR0YHe4-Mo3DpJ2nwlTV3PuhEJc31pMo"`
	CreatedAt      time.Time `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt      time.Time `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
	DeletedAt      *string   `json:"deleted_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
}

type UserInformationResponses struct {
	ID         uint    `form:"user_id" json:"user_id"`
	FullName       string  `form:"full_name" json:"full_name" example:"Nadiah Indrianti"`
	Email          string  `form:"email" json:"email" example:"user@nadiah.com"`
	NIMNIP         string  `form:"nim_nip" json:"nim_nip" example:"2009189"`
	ProfilePicture string `json:"profile_picture_url" example:"https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg"`
}

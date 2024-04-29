package dtos

import "time"

type NotificationInput struct {
	UserID        uint `json:"user_id" form:"user_id"`
	TemplateID    uint `json:"template_id" form:"template_id"`
	PeminjamanID  uint `json:"peminjaman_id" form:"peminjaman_id"`
}

type NotificationCreateResponse struct {
	NotificationID      uint                    `json:"notification_id" form:"notification_id"`
	UserID              uint                    `json:"user_id" form:"user_id"`
	TemplateID          uint                    `json:"template_id" form:"template_id"`
	PeminjamanID  		uint 					`json:"peminjaman_id" form:"peminjaman_id"`
	NotificationContent TemplateMessageResponse `json:"notification_content" form:"notification_content"`
	CreatedAt           time.Time               `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt           time.Time               `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}

type NotificationResponse struct {
	UserID              uint                              `json:"user_id" form:"user_id"`
	NotificationContent []TemplateMessageByUserIDResponse `json:"notification_content" form:"notification_content"`
}

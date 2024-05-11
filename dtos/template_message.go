package dtos

import "time"

type TemplateMessageInput struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
}

type TemplateMessageResponse struct {
	TemplateID uint      `json:"template_id" form:"template_id"`
	Title      string    `json:"title" form:"title"`
	Content    string    `json:"content" form:"content"`
	CreatedAt  time.Time `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt  time.Time `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}

type TemplateMessageByUserIDResponse struct {
	Title     string    `json:"title" form:"title"`
	Content   string    `json:"content" form:"content"`
	CreatedAt time.Time `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}

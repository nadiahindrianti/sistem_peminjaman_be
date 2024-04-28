package models

import "gorm.io/gorm"

type TemplateMessage struct {
	gorm.Model
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
}

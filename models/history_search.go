package models

import "gorm.io/gorm"

type HistorySearch struct {
	gorm.Model
	UserID uint   `json:"user_id" form:"user_id"`
	Name   string `json:"name" form:"name"`
}

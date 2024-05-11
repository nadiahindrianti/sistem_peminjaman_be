package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID        		 uint            `json:"user_id" form:"user_id"`
	User          		 User            `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TemplateID    		 uint            `json:"template_id" form:"template_id"`
	TemplateMessage      TemplateMessage `gorm:"foreignKey:TemplateID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

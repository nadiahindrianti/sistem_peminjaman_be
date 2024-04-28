package models

import "gorm.io/gorm"

type Lab struct {
	gorm.Model
	Name        string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
	
}
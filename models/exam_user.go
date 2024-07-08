package models

import "gorm.io/gorm"

type ExamUser struct {
	gorm.Model
	FullName       string
	Email          string `gorm:"unique"`
	Password       string
	NIMNIP         string
	ProfilePicture string
	Role           string `gorm:"type:ENUM('user','admin')"`
}
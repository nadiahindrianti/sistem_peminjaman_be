package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName       string
	Email          string `gorm:"unique"`
	Password       string
	NIMNIP         string
	KartuIdentitas string
	ProfilePicture string
	Role           string `gorm:"type:ENUM('user','admin')"`
}



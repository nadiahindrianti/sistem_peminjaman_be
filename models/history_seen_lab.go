package models

import "gorm.io/gorm"

type HistorySeenLab struct {
	gorm.Model
	UserID  uint
	User    User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LabID uint
	Lab   Lab    `gorm:"foreignKey:LabID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

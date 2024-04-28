package models

import (
	"time"

	"gorm.io/gorm"
)

type Peminjaman struct {
	gorm.Model
	UserID         uint      `form:"user_id" json:"user_id"`
	User           User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LabID          uint      `form:"lab_id" json:"lab_id"`
	Lab            Lab       `gorm:"foreignKey:LabID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TanggalPeminjaman time.Time `gorm:"type:DATE"`
	WaktuPeminjaman  string    `gorm:"type:ENUM('09:00', '12:00', '15:00')"`
	Description    string    `form:"description" json:"description"`
	Status         string    `gorm:"type:ENUM('request', 'accept', 'reject')"`
}

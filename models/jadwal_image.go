package models

import "gorm.io/gorm"

type JadwalImage struct {
	gorm.Model
	JadwalID    uint   `form:"jadwal_id" json:"jadwal_id"`
	Jadwal      Jadwal    `gorm:"foreignKey:JadwalID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageUrl    string `form:"image_url" json:"image_url"`
}

package models

import "gorm.io/gorm"

type BeritaAcaraImage struct {
	gorm.Model
	JadwalID    		   uint   `form:"jadwal_id" json:"jadwal_id"`
	Jadwal      		   Jadwal    `gorm:"foreignKey:JadwalID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BeritaAcaraImageUrl    string `form:"beritaacara_image_url" json:"beritaacara_image_url"`
}
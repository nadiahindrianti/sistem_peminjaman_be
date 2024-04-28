package models

import "gorm.io/gorm"

type SuratRekomendasiImage struct {
	gorm.Model
	PeminjamanID    uint   `form:"peminjaman_id" json:"peminjaman_id"`
	Peminjaman      Peminjaman    `gorm:"foreignKey:PeminjamanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageUrl string `form:"image_url" json:"image_url"`
}

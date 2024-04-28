package models

import "gorm.io/gorm"

type LabImage struct {
	gorm.Model
	LabID    uint   `form:"lab_id" json:"lab_id"`
	Lab      Lab    `gorm:"foreignKey:LabID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageUrl string `form:"image_url" json:"image_url"`
}

package models

import "gorm.io/gorm"

type Jadwal struct {
	gorm.Model
  TanggalJadwal          time.Time `gorm:"type:DATE"`
	WaktuJadwal            string    `gorm:"type:ENUM('09:00', '12:00', '15:00')"`
	NameUser               string    `form:"name_user" json:"name_user"`
  NameLaboratorium       string    `form:"name_lab" json:"name_lab"`
	Status                 string    `gorm:"type:ENUM('notused', 'inused', 'finished')"`
	
}

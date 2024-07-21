package dtos

import "time"


type JadwalInput struct {
	TanggalJadwal           *string 					   `form:"tanggal_jadwal" json:"tanggal_jadwal,omitempty" example:"2002-09-12"`
	WaktuJadwal             string    					   `form:"waktu_jadwal" json:"waktu_jadwal" example:"09:00"`
	NameUser                string    					   `form:"name_user" json:"name_user"`
  	NameLaboratorium        string    					   `form:"name_lab" json:"name_lab"`
	BeritaAcaraImage        []BeritaAcaraImageInput        `form:"beritaacara_image" json:"beritaacara_image"`
	Status                  string    					   `form:"status" json:"status" example:"notused"`
}

type StatusJadwalResponse struct {
	Status                 	    string    					   	    `form:"status" json:"status" example:"request"`
}

type JadwalResponse struct {
	JadwalID           		int                      	   `form:"jadwal_id" json:"jadwal_id"`
	TanggalJadwal           string 					   	   `form:"tanggal_jadwal" json:"tanggal_jadwal,omitempty" example:"2002-09-12"`
	WaktuJadwal             string    					   `form:"waktu_jadwal" json:"waktu_jadwal" example:"09:00"`
	NameUser                string    					   `form:"name_user" json:"name_user"`
  	NameLaboratorium        string    					   `form:"name_lab" json:"name_lab"`
	BeritaAcaraImage        []BeritaAcaraImageResponse     `form:"beritaacara_image" json:"beritaacara_image"`
	Status                  string    					   `form:"status" json:"status" example:"notused"`
	CreatedAt       		time.Time                 	   `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt       		time.Time                 	   `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}

type JadwalUserResponses struct {
	JadwalID           		uint                      	   `form:"jadwal_id" json:"jadwal_id"`
	TanggalJadwal           string 					   `form:"tanggal_jadwal" json:"tanggal_jadwal,omitempty" example:"2002-09-12"`
	WaktuJadwal             string    					    `form:"waktu_jadwal" json:"waktu_jadwal" example:"09:00"`
	NameUser                string    					   `form:"name_user" json:"name_user"`
  	NameLaboratorium        string    					   `form:"name_lab" json:"name_lab"`
	Status                  string    					   `form:"status" json:"status" example:"notused"`
	CreatedAt       		*time.Time                	   `json:"created_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt       		*time.Time                     `json:"updated_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
}


package dtos

import "time"

type PeminjamanInput struct {
	LabID           			int                      			`form:"lab_id" json:"lab_id" example:"1"`
	TanggalPeminjaman  			*string 					   		`form:"tanggal_peminjaman" json:"tanggal_peminjaman,omitempty" example:"2002-09-12"`
	JamPeminjaman    			string    							`form:"jam_peminjaman" json:"jam_peminjaman" example:"09:00"`
	SuratRekomendasiImage       []SuratRekomendasiImageInput        `form:"suratrekomendasi_image" json:"suratrekomendasi_image"`
	Description     			string                 				`form:"description" json:"description"`
	Status                 	    string    					   	    `form:"status" json:"status" example:"request"`
}

type PeminjamanResponse struct {
	PeminjamanID     			int                       			`json:"peminjaman_id" example:"1"`
	TanggalPeminjaman  			string 					   			`form:"tanggal_peminjaman" json:"tanggal_peminjaman,omitempty" example:"2002-09-12"`
	JamPeminjaman    			string    							`form:"jam_peminjaman" json:"jam_peminjaman" example:"09:00"`
	SuratRekomendasiImage       []SuratRekomendasiImageResponse     `form:"suratrekomendasi_image" json:"suratrekomendasi_image"`
	Description     			string                 				`form:"description" json:"description"`
	Status                 	    string    					   	    `form:"status" json:"status" example:"request"`
	Lab            				LabByIDResponses        			`json:"lab"`
	User           			   *UserInformationResponses 			`json:"user,omitempty"`
	CreatedAt        			time.Time                			`json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt        			time.Time                 			`json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}



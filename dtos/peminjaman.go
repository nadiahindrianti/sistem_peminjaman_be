package dtos

import "time"

type PeminjamanInput struct {
	TanggalPeminjaman  			string    							`form:"tanggal_peminjaman" json:"tanggal_peminjaman" example:"2023-05-01"`
	WaktuPeminjaman    			string    							`form:"waktu_peminjaman" json:"waktu_peminjaman" example:"2023-05-01"`
	SuratRekomendasiImage       []SuratRekomendasiImageInput        `form:"surat_rekomendasi_image" json:"surat_rekomendasi_image"`
	Description     			string                 				`form:"description" json:"description"`
}

type PeminjamanResponse struct {
	PeminjamanID     			int                       			`json:"peminjaman_id" example:"1"`
	TanggalPeminjaman  			string    							`form:"tanggal_peminjaman" json:"tanggal_peminjaman" example:"2023-05-01"`
	WaktuPeminjaman    			string    							`form:"waktu_peminjaman" json:"waktu_peminjaman" example:"2023-05-01"`
	SuratRekomendasiImage       []SuratRekomendasiImageResponse     `form:"surat_rekomendasi_image" json:"surat_rekomendasi_image"`
	Description     			string                 				`form:"description" json:"description"`
	Status           			string                    			`json:"status" example:"unpaid"`
	Lab            				LabByIDResponses        			`json:"lab"`
	User           			   *UserInformationResponses 			`json:"user,omitempty"`
	CreatedAt        			time.Time                			`json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt        			time.Time                 			`json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}



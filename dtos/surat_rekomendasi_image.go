package dtos

type SuratRekomendasiImageInput struct {
	SuratRekomendasiImageUrl string `form:"suratrekomendasi_image_url" json:"suratrekomendasi_image_url"`
}

type SuratRekomendasiImageResponse struct {
	PeminjamanID  			 uint   `form:"peminjaman_id" json:"peminjaman_id"`
	SuratRekomendasiImageUrl string `form:"suratrekomendasi_image_url" json:"suratrekomendasi_image_url"`
}

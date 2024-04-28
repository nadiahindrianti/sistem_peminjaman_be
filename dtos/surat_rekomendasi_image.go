package dtos

type SuratRekomendasiImageInput struct {
	ImageUrl string `form:"image_url" json:"image_url"`
}

type SuratRekomendasiImageResponse struct {
	PeminjamanIDID  uint   `form:"peminjaman_id" json:"peminjaman_id"`
	ImageUrl string `form:"image_url" json:"image_url"`
}

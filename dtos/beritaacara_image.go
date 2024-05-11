package dtos

type BeritaAcaraImageInput struct {
	BeritaAcaraImageUrl string `form:"beritaacara_image_url" json:"beritaacara_image_url"`
}

type BeritaAcaraImageResponse struct {
	JadwalID  			uint   `form:"jadwal_id" json:"jadwal_id"`
	BeritaAcaraImageUrl string `form:"beritaacara_image_url" json:"beritaacara_image_url"`
}

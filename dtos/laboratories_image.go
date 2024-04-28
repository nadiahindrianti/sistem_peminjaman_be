package dtos

type LabImageInput struct {
	ImageUrl string `form:"image_url" json:"image_url"`
}

type LabImageResponse struct {
	LabID  uint   `form:"lab_id" json:"lab_id"`
	ImageUrl string `form:"image_url" json:"image_url"`
}

package dtos

import "time"

type LabInput struct {
	Name            string                 `form:"name" json:"name"`
	LabImage        []LabImageInput        `form:"lab_image" json:"lab_image"`
	Description     string                 `form:"description" json:"description"`
}

type LabResponse struct {
	LabID           uint                      `form:"lab_id" json:"lab_id"`
	Name            string                    `form:"name" json:"name"`  
	LabImage        []LabImageResponse        `form:"lab_image" json:"lab_image"`
	Description     string                    `form:"description" json:"description"`
	CreatedAt       time.Time                 `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt       time.Time                 `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}
type LabByIDResponse struct {
	LabID           uint                       `form:"lab_id" json:"lab_id"`
	Name            string                     `form:"name" json:"name"`
	LabImage        []LabImageResponse         `form:"lab_image" json:"lab_image"`
	Description     string                     `form:"description" json:"description"`
	CreatedAt       time.Time                  `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt       time.Time                  `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}

type LabByIDResponses struct {
	LabID           uint                      `form:"lab_id" json:"lab_id"`
	Name            string                    `form:"name" json:"name"`
	Description     string                    `form:"description" json:"description"`
	CreatedAt       *time.Time                `json:"created_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt       *time.Time                `json:"updated_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
}

type LabByIDSimply struct {
	LabID           uint                      `form:"lab_id" json:"lab_id"`
	Name            string                    `form:"name" json:"name"`
	LabImage        []LabImageResponse        `form:"lab_image" json:"lab_image"`
	Description     string                    `form:"description" json:"description"`
	CreatedAt       *time.Time                `json:"created_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt       *time.Time                `json:"updated_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
}

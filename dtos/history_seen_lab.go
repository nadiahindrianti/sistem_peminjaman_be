package dtos

import "time"

type HistorySeenLabInput struct {
	LabID uint `json:"lab_id" form:"lab_id" example:"1"`
}

type HistorySeenLabResponse struct {
	ID        uint            `json:"id" example:"1"`
	Lab       LabByIDSimply `json:"lab"`
	CreatedAt time.Time       `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt time.Time       `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}

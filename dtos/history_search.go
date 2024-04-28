package dtos

type HistorySearchInput struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type HistorySearchResponse struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
}

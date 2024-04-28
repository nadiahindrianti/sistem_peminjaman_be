package repositories

import (
	"sistem_peminjaman_be/models"

	"gorm.io/gorm"
)

type HistorySearchRepository interface {
	HistorySearchGetById(userId, id uint) (models.HistorySearch, error)
	HistorySearchGetByUserId(userId uint, page, limit int) ([]models.HistorySearch, int, error)
	HistorySearchCreate(historySearch models.HistorySearch) (models.HistorySearch, error)
	HistorySearchDelete(userId, id uint) error
}

type historySearchRepository struct {
	db *gorm.DB
}

func NewHistorySearchRepository(db *gorm.DB) HistorySearchRepository {
	return &historySearchRepository{db}
}

func (r *historySearchRepository) HistorySearchGetById(userId, id uint) (models.HistorySearch, error) {
	var historySearch models.HistorySearch
	err := r.db.Where("id = ? AND user_id = ?", id, userId).First(&historySearch).Error
	return historySearch, err
}

func (r *historySearchRepository) HistorySearchGetByUserId(userId uint, page, limit int) ([]models.HistorySearch, int, error) {
	var historySearch []models.HistorySearch
	var count int64
	err := r.db.Where("user_id = ?", userId).Find(&historySearch).Count(&count).Error
	if err != nil {
		return historySearch, 0, err
	}

	offset := (page - 1) * limit

	err = r.db.Where("user_id = ?", userId).Limit(limit).Offset(offset).Find(&historySearch).Error

	return historySearch, int(count), err
}

func (r *historySearchRepository) HistorySearchCreate(historySearch models.HistorySearch) (models.HistorySearch, error) {
	err := r.db.Create(&historySearch).Error
	return historySearch, err
}

func (r *historySearchRepository) HistorySearchDelete(userId, id uint) error {
	var historySearch models.HistorySearch
	err := r.db.Where("id = ? AND user_id = ?", id, userId).Delete(&historySearch).Error
	return err
}

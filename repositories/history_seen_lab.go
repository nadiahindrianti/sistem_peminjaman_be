package repositories

import (
	"sistem_peminjaman_be/models"

	"gorm.io/gorm"
)

type HistorySeenLabRepository interface {
	GetAllHistorySeenLab(page, limit int, userId uint) ([]models.HistorySeenLab, int, error)
	GetHistorySeenLabByID(labId, userId uint) (models.HistorySeenLab, error)
	CreateHistorySeenLab(HistorySeenLab models.HistorySeenLab) (models.HistorySeenLab, error)
	UpdateHistorySeenLab(HistorySeenLab models.HistorySeenLab) (models.HistorySeenLab, error)
	DeleteHistorySeenLab(HistorySeenLab models.HistorySeenLab) (models.HistorySeenLab, error)
}

type historySeenLabRepository struct {
	db *gorm.DB
}

func NewHistorySeenLabRepository(db *gorm.DB) HistorySeenLabRepository {
	return &historySeenLabRepository{db}
}

func (r *historySeenLabRepository) GetAllHistorySeenLab(page, limit int, userId uint) ([]models.HistorySeenLab, int, error) {
	var (
		histories []models.HistorySeenLab
		count     int64
	)
	err := r.db.Where("user_id = ?", userId).Order("id DESC").Find(&histories).Count(&count).Error
	if err != nil {
		return histories, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Where("user_id = ?", userId).Order("id DESC").Limit(limit).Offset(offset).Find(&histories).Error

	return histories, int(count), err
}

func (r *historySeenLabRepository) GetHistorySeenLabByID(labId, userId uint) (models.HistorySeenLab, error) {
	var historySeenLab models.HistorySeenLab
	err := r.db.Where("lab_id = ? AND user_id = ?", labId, userId).First(&historySeenLab).Error
	return historySeenLab, err
}

func (r *historySeenLabRepository) CreateHistorySeenLab(historySeenLab models.HistorySeenLab) (models.HistorySeenLab, error) {
	err := r.db.Create(&historySeenLab).Error
	return historySeenLab, err
}

func (r *historySeenLabRepository) UpdateHistorySeenLab(historySeenLab models.HistorySeenLab) (models.HistorySeenLab, error) {
	err := r.db.Save(&historySeenLab).Error
	return historySeenLab, err
}

func (r *historySeenLabRepository) DeleteHistorySeenLab(historySeenLab models.HistorySeenLab) (models.HistorySeenLab, error) {
	err := r.db.Unscoped().Delete(&historySeenLab).Error
	return historySeenLab, err
}

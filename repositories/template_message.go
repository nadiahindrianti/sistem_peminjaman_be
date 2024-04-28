package repositories

import (
	"sistem_peminjaman_be/models"

	"gorm.io/gorm"
)

type TemplateMessageRepository interface {
	GetAllTemplateMessages(page, limit int) ([]models.TemplateMessage, int, error)
	GetTemplateMessageByID(id uint) (models.TemplateMessage, error)
	CreateTemplateMessage(template models.TemplateMessage) (models.TemplateMessage, error)
	UpdateTemplate(template models.TemplateMessage) (models.TemplateMessage, error)
	DeleteTemplate(id uint) error
}

type templateMessageRepository struct {
	db *gorm.DB
}

func NewTemplateMessageRepository(db *gorm.DB) TemplateMessageRepository {
	return &templateMessageRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *templateMessageRepository) GetAllTemplateMessages(page, limit int) ([]models.TemplateMessage, int, error) {
	var (
		templates []models.TemplateMessage
		count     int64
	)
	err := r.db.Find(&templates).Count(&count).Error
	if err != nil {
		return templates, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Order("id DESC").Limit(limit).Offset(offset).Find(&templates).Error

	return templates, int(count), err
}

func (r *templateMessageRepository) GetTemplateMessageByID(id uint) (models.TemplateMessage, error) {
	var template models.TemplateMessage
	err := r.db.Where("id = ?", id).First(&template).Error
	return template, err
}

func (r *templateMessageRepository) CreateTemplateMessage(template models.TemplateMessage) (models.TemplateMessage, error) {
	err := r.db.Create(&template).Error
	return template, err
}

func (r *templateMessageRepository) UpdateTemplate(template models.TemplateMessage) (models.TemplateMessage, error) {
	err := r.db.Save(&template).Error
	return template, err
}

func (r *templateMessageRepository) DeleteTemplate(id uint) error {
	var template models.TemplateMessage
	err := r.db.Where("id = ?", id).Delete(&template).Error
	return err
}

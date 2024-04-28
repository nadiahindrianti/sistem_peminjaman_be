package repositories

import (
	"sistem_peminjaman_be/models"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	GetNotificationByUserID(id uint) ([]models.Notification, error)
	CreateNotification(notification models.Notification) (models.Notification, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r notificationRepository) GetNotificationByUserID(id uint) ([]models.Notification, error) {
	var notification []models.Notification
	err := r.db.Where("user_id = ?", id).Order("created_at DESC").Find(&notification).Error
	return notification, err
}

func (r notificationRepository) CreateNotification(notification models.Notification) (models.Notification, error) {
	err := r.db.Create(&notification).Error
	return notification, err
}

package repositories

import (
	"sistem_peminjaman_be/models"

	"gorm.io/gorm"
)

type LabImageRepository interface {
	GetAllLabImages(page, limit int) ([]models.LabImage, int, error)
	GetAllLabImageByID(id uint) ([]models.LabImage, error)
	GetLabImageByID(id uint) (models.LabImage, error)
	CreateLabImage(labImage models.LabImage) (models.LabImage, error)
	UpdateLabImage(labImage models.LabImage) (models.LabImage, error)
	DeleteLabImage(id uint) error
}

type labImageRepository struct {
	db *gorm.DB
}

func NewLabImageRepository(db *gorm.DB) LabImageRepository {
	return &labImageRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *labImageRepository) GetAllLabImages(page, limit int) ([]models.LabImage, int, error) {
	var (
		labs []models.LabImage
		count  int64
	)
	err := r.db.Find(&labs).Count(&count).Error
	if err != nil {
		return labs, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&labs).Error

	return labs, int(count), err
}

func (r *labImageRepository) GetAllLabImageByID(id uint) ([]models.LabImage, error) {
	var labImage []models.LabImage
	err := r.db.Where("lab_id = ?", id).Find(&labImage).Error
	return labImage, err
}

func (r *labImageRepository) GetLabImageByID(id uint) (models.LabImage, error) {
	var labImage models.LabImage
	err := r.db.Where("id = ?", id).First(&labImage).Error
	return labImage, err
}

func (r *labImageRepository) CreateLabImage(labImage models.LabImage) (models.LabImage, error) {
	err := r.db.Create(&labImage).Error
	return labImage, err
}

func (r *labImageRepository) UpdateLabImage(labImage models.LabImage) (models.LabImage, error) {
	err := r.db.Save(&labImage).Error
	return labImage, err
}

func (r *labImageRepository) DeleteLabImage(id uint) error {
	var labImage models.LabImage
	err := r.db.Unscoped().Where("lab_id = ?", id).Delete(&labImage).Error
	return err
}

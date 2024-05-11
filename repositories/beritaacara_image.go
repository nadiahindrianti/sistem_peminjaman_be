package repositories

import (
	"sistem_peminjaman_be/models"

	"gorm.io/gorm"
)

type BeritaAcaraImageRepository interface {
	GetAllBeritaAcaraImages(page, limit int) ([]models.BeritaAcaraImage, int, error)
	GetAllBeritaAcaraImageByID(id uint) ([]models.BeritaAcaraImage, error)
	GetBeritaAcaraImageByID(id uint) (models.BeritaAcaraImage, error)
	CreateBeritaAcaraImage(beritaacaraImage models.BeritaAcaraImage) (models.BeritaAcaraImage, error)
	UpdateBeritaAcaraImage(beritaacaraImage models.BeritaAcaraImage) (models.BeritaAcaraImage, error)
	DeleteBeritaAcaraImage(id uint) error
}

type beritaAcaraImageRepository struct {
	db *gorm.DB
}

func NewBeritaAcaraImageRepository(db *gorm.DB) BeritaAcaraImageRepository {
	return &beritaAcaraImageRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *beritaAcaraImageRepository) GetAllBeritaAcaraImages(page, limit int) ([]models.BeritaAcaraImage, int, error) {
	var (
		jadwals []models.BeritaAcaraImage
		count  int64
	)
	err := r.db.Find(&jadwals).Count(&count).Error
	if err != nil {
		return jadwals, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&jadwals).Error

	return jadwals, int(count), err
}

func (r *beritaAcaraImageRepository) GetAllBeritaAcaraImageByID(id uint) ([]models.BeritaAcaraImage, error) {
	var beritaAcaraImage []models.BeritaAcaraImage
	err := r.db.Where("jadwal_id = ?", id).Find(&beritaAcaraImage).Error
	return beritaAcaraImage, err
}

func (r *beritaAcaraImageRepository) GetBeritaAcaraImageByID(id uint) (models.BeritaAcaraImage, error) {
	var beritaAcaraImage models.BeritaAcaraImage
	err := r.db.Where("id = ?", id).First(&beritaAcaraImage).Error
	return beritaAcaraImage, err
}

func (r *beritaAcaraImageRepository) CreateBeritaAcaraImage(beritaAcaraImage models.BeritaAcaraImage) (models.BeritaAcaraImage, error) {
	err := r.db.Create(&beritaAcaraImage).Error
	return beritaAcaraImage, err
}

func (r *beritaAcaraImageRepository) UpdateBeritaAcaraImage(beritaAcaraImage models.BeritaAcaraImage) (models.BeritaAcaraImage, error) {
	err := r.db.Save(&beritaAcaraImage).Error
	return beritaAcaraImage, err
}

func (r *beritaAcaraImageRepository) DeleteBeritaAcaraImage(id uint) error {
	var beritaAcaraImage models.BeritaAcaraImage
	err := r.db.Unscoped().Where("jadwal_id = ?", id).Delete(&beritaAcaraImage).Error
	return err
}

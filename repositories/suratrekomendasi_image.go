package repositories

import (
	"sistem_peminjaman_be/models"

	"gorm.io/gorm"
)

type SuratRekomendasiImageRepository interface {
	GetAllSuratRekomendasiImages(page, limit int) ([]models.SuratRekomendasiImage, int, error)
	GetAllSuratRekomendasiImageByID(id uint) ([]models.SuratRekomendasiImage, error)
	GetSuratRekomendasiImageByID(id uint) (models.SuratRekomendasiImage, error)
	CreateSuratRekomendasiImage(suratrekomendasiImage models.SuratRekomendasiImage) (models.SuratRekomendasiImage, error)
	UpdateSuratRekomendasiImage(suratrekomendasiImage models.SuratRekomendasiImage) (models.SuratRekomendasiImage, error)
	DeleteSuratRekomendasiImage(id uint) error
}

type suratRekomendasiImageRepository struct {
	db *gorm.DB
}

func NewSuratRekomendasiImageRepository(db *gorm.DB) SuratRekomendasiImageRepository {
	return &suratRekomendasiImageRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *suratRekomendasiImageRepository) GetAllSuratRekomendasiImages(page, limit int) ([]models.SuratRekomendasiImage, int, error) {
	var (
		peminjamans []models.SuratRekomendasiImage
		count  int64
	)
	err := r.db.Find(&peminjamans).Count(&count).Error
	if err != nil {
		return peminjamans, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&peminjamans).Error

	return peminjamans, int(count), err
}

func (r *suratRekomendasiImageRepository) GetAllSuratRekomendasiImageByID(id uint) ([]models.SuratRekomendasiImage, error) {
	var suratRekomendasiImage []models.SuratRekomendasiImage
	err := r.db.Where("peminjaman_id = ?", id).Find(&suratRekomendasiImage).Error
	return suratRekomendasiImage, err
}

func (r *suratRekomendasiImageRepository) GetSuratRekomendasiImageByID(id uint) (models.SuratRekomendasiImage, error) {
	var suratRekomendasiImage models.SuratRekomendasiImage
	err := r.db.Where("id = ?", id).First(&suratRekomendasiImage).Error
	return suratRekomendasiImage, err
}

func (r *suratRekomendasiImageRepository) CreateSuratRekomendasiImage(suratRekomendasiImage models.SuratRekomendasiImage) (models.SuratRekomendasiImage, error) {
	err := r.db.Create(&suratRekomendasiImage).Error
	return suratRekomendasiImage, err
}

func (r *suratRekomendasiImageRepository) UpdateSuratRekomendasiImage(suratRekomendasiImage models.SuratRekomendasiImage) (models.SuratRekomendasiImage, error) {
	err := r.db.Save(&suratRekomendasiImage).Error
	return suratRekomendasiImage, err
}

func (r *suratRekomendasiImageRepository) DeleteSuratRekomendasiImage(id uint) error {
	var suratRekomendasiImage models.SuratRekomendasiImage
	err := r.db.Unscoped().Where("peminjaman_id = ?", id).Delete(&suratRekomendasiImage).Error
	return err
}

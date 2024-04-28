package repositories

import (
	"sistem_peminjaman_be/models"

	"gorm.io/gorm"
)

type LabRepository interface {
	GetAllLabs(page, limit int) ([]models.Lab, int, error)
	GetLabByID(id uint) (models.Lab, error)
	GetLabByID2(id uint) (models.Lab, error)
	CreateLab(Lab models.Lab) (models.Lab, error)
	UpdateLab(Lab models.Lab) (models.Lab, error)
	DeleteLab(id uint) error
	SearchLabAvailable(page, limit int, name string) ([]models.Lab, int, error)
}

type labRepository struct {
	db *gorm.DB
}

func NewLabRepository(db *gorm.DB) LabRepository {
	return &labRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *labRepository) GetAllLabs(page, limit int) ([]models.Lab, int, error) {
	var (
		labs []models.Lab
		count  int64
	)
	err := r.db.Find(&labs).Count(&count).Error
	if err != nil {
		return labs, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Order("id DESC").Limit(limit).Offset(offset).Find(&labs).Error

	return labs, int(count), err
}

func (r *labRepository) GetLabByID(id uint) (models.Lab, error) {
	var lab models.Lab
	err := r.db.Where("id = ?", id).First(&lab).Error
	return lab, err
}

func (r *labRepository) GetLabByID2(id uint) (models.Lab, error) {
	var lab models.Lab
	err := r.db.Unscoped().Where("id = ?", id).First(&lab).Error
	return lab, err
}

func (r *labRepository) CreateLab(lab models.Lab) (models.Lab, error) {
	err := r.db.Create(&lab).Error
	return lab, err
}

func (r *labRepository) UpdateLab(lab models.Lab) (models.Lab, error) {
	err := r.db.Save(&lab).Error
	return lab, err
}

func (r *labRepository) DeleteLab(id uint) error {
	var lab models.Lab
	err := r.db.Where("id = ?", id).Delete(&lab).Error
	return err
}

func (r *labRepository) SearchLabAvailable(page, limit int, name string) ([]models.Lab, int, error) {
	var (
		labs []models.Lab
		count  int64
		err    error
	)

	if  name == "" {
		err = r.db.Find(&labs).Count(&count).Error
	}
	if  name != "" {
		err = r.db.Where("name LIKE ?", "%"+name+"%").Find(&labs).Count(&count).Error
	}

	if err != nil {
		return labs, int(count), err
	}

	offset := (page - 1) * limit

	if  name == "" {
		err = r.db.Order("id DESC").Limit(limit).Offset(offset).Find(&labs).Error
	}
	if  name != "" {
		err = r.db.Where("name LIKE ?", "%"+name+"%").Order("id DESC").Limit(limit).Offset(offset).Find(&labs).Error
	}

	return labs, int(count), err

}

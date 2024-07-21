package repositories

import (
	"sistem_peminjaman_be/models"

	"gorm.io/gorm"
)

type JadwalRepository interface {
	GetAllJadwals(page, limit int) ([]models.Jadwal, int, error)
	GetJadwalByID(id uint) (models.Jadwal, error)
	GetJadwalByID2(id uint) (models.Jadwal, error)
	GetJadwalByID3(id uint, userID uint) (models.Jadwal, error)
	CreateJadwal(jadwal models.Jadwal) (models.Jadwal, error)
	UpdateJadwal(jadwal models.Jadwal) (models.Jadwal, error)
	DeleteJadwal(id uint) error
	SearchJadwalAvailable(page, limit int, name_laboratorium string) ([]models.Jadwal, int, error)
}

type jadwalRepository struct {
	db *gorm.DB
}

func NewJadwalRepository(db *gorm.DB) JadwalRepository {
	return &jadwalRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *jadwalRepository) GetAllJadwals(page, limit int) ([]models.Jadwal, int, error) {
	var (
		jadwals []models.Jadwal
		count  int64
	)
	err := r.db.Find(&jadwals).Count(&count).Error
	if err != nil {
		return jadwals, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Order("id DESC").Limit(limit).Offset(offset).Find(&jadwals).Error

	return jadwals, int(count), err
}

func (r *jadwalRepository) GetJadwalByID(id uint) (models.Jadwal, error) {
	var jadwal models.Jadwal
	err := r.db.Where("id = ?", id).First(&jadwal).Error
	return jadwal, err
}

func (r *jadwalRepository) GetJadwalByID2(id uint) (models.Jadwal, error) {
	var jadwal models.Jadwal
	err := r.db.Unscoped().Where("id = ?", id).First(&jadwal).Error
	return jadwal, err
}

func (r *jadwalRepository) GetJadwalByID3(id uint, userID uint) (models.Jadwal, error) {
	var jadwal models.Jadwal
	if userID == 1 {
		err := r.db.Where("id = ?", id).First(&jadwal).Error
		return jadwal, err
	}
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&jadwal).Error
	return jadwal, err
}

func (r *jadwalRepository) CreateJadwal(jadwal models.Jadwal) (models.Jadwal, error) {
	err := r.db.Create(&jadwal).Error
	return jadwal, err
}

func (r *jadwalRepository) UpdateJadwal(jadwal models.Jadwal) (models.Jadwal, error) {
	err := r.db.Save(&jadwal).Error
	return jadwal, err
}

func (r *jadwalRepository) DeleteJadwal(id uint) error {
	var jadwal models.Jadwal
	err := r.db.Where("id = ?", id).Delete(&jadwal).Error
	return err
}

func (r *jadwalRepository) SearchJadwalAvailable(page, limit int, name_laboratorium string) ([]models.Jadwal, int, error) {
	var (
		jadwals []models.Jadwal
		count  int64
		err    error
	)

	if  name_laboratorium == "" {
		err = r.db.Find(&jadwals).Count(&count).Error
	}
	if  name_laboratorium != "" {
		err = r.db.Where("name_laboratorium LIKE ?", "%"+name_laboratorium+"%").Find(&jadwals).Count(&count).Error
	}

	if err != nil {
		return jadwals, int(count), err
	}

	offset := (page - 1) * limit

	if  name_laboratorium == "" {
		err = r.db.Order("id DESC").Limit(limit).Offset(offset).Find(&jadwals).Error
	}
	if  name_laboratorium != "" {
		err = r.db.Where("name_laboratorium LIKE ?", "%"+name_laboratorium+"%").Order("id DESC").Limit(limit).Offset(offset).Find(&jadwals).Error
	}

	return jadwals, int(count), err

}
package repositories

import (
	"sistem_peminjaman_be/models"
	"time"

	"gorm.io/gorm"
)

type PeminjamanRepository interface {
	GetPeminjamans(page, limit int, userID uint, status string) ([]models.Peminjaman, int, error)
	GetPeminjamanByStatusAndID(id, userID uint, status string) (models.Peminjaman, error)
	GetPeminjamanByID(id, userID uint) (models.Peminjaman, error)
	CreatePeminjaman(peminjaman models.Peminjaman) (models.Peminjaman, error)
	UpdatePeminjaman(peminjaman models.Peminjaman) (models.Peminjaman, error)
	DeletePeminjaman(peminjaman models.Peminjaman) (models.Peminjaman, error)
}

type peminjamanRepository struct {
	db *gorm.DB
}

func NewPeminjamanRepository(db *gorm.DB) PeminjamanRepository {
	return &peminjamanRepository{db}
}

func (r *peminjamanRepository) GetPeminjamans(page, limit int, userID uint, status string) ([]models.Peminjaman, int, error) {
	var (
		peminjamans []models.Peminjaman
		count       int64
		err         error
	)
	if userID == 1 {
		if status == "" {
			err = r.db.Find(&peminjamans).Count(&count).Error
		} else {
			err = r.db.Where("status = ?", status).Find(&peminjamans).Count(&count).Error
		}
	} else {
		if status == "" {
			err = r.db.Where("user_id = ?", userID).Find(&peminjamans).Count(&count).Error
		} else {
			err = r.db.Where("user_id = ? AND status = ?", userID, status).Find(&peminjamans).Count(&count).Error
		}
	}
	if err != nil {
		return peminjamans, int(count), err
	}

	offset := (page - 1) * limit

	if userID == 1 {
		if status == "" {
			err = r.db.Find(&peminjamans).Count(&count).Error
		} else {
			err = r.db.Where("status = ?", status).Limit(limit).Offset(offset).Find(&peminjamans).Error
		}
	} else {
		if status == "" {
			err = r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&peminjamans).Error
		} else {
			err = r.db.Where("user_id = ? AND status = ?", userID, status).Limit(limit).Offset(offset).Find(&peminjamans).Error
		}
	}

	return peminjamans, int(count), err
}

func (r *peminjamanRepository) GetPeminjamanByStatusAndID(id, userID uint, status string) (models.Peminjaman, error) {
	var peminjaman models.Peminjaman
	if userID == 1 {
		err := r.db.Where("id = ? AND status = ?", id, status).First(&peminjaman).Error
		return peminjaman, err
	}
	err := r.db.Where("id = ? AND user_id = ? AND status = ?", id, userID, status).First(&peminjaman).Error
	return peminjaman, err
}

func (r *peminjamanRepository) GetPeminjamanByID(id, userID uint) (models.Peminjaman, error) {
	var peminjaman models.Peminjaman
	if userID == 1 {
		err := r.db.Where("id = ?", id).First(&peminjaman).Error
		return peminjaman, err
	}
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&peminjaman).Error
	return peminjaman, err
}

func (r *peminjamanRepository) CreatePeminjaman(peminjaman models.Peminjaman) (models.Peminjaman, error) {
	err := r.db.Create(&peminjaman).Error
	return peminjaman, err
}

func (r *peminjamanRepository) UpdatePeminjaman(peminjaman models.Peminjaman) (models.Peminjaman, error) {
	err := r.db.Save(peminjaman).Error
	return peminjaman, err
}

func (r *peminjamanRepository) DeletePeminjaman(peminjaman models.Peminjaman) (models.Peminjaman, error) {
	err := r.db.Unscoped().Delete(&peminjaman).Error
	return peminjaman, err
}

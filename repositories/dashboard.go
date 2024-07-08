package repositories

import (
	"sistem_peminjaman_be/models"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	DashboardGetAll() (int, int, int, int, int, []models.Peminjaman, []models.User, []models.Jadwal, int, int, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db}
}

func (r *dashboardRepository) DashboardGetAll() (int, int, int, int, int, []models.Peminjaman, []models.User, []models.Jadwal, int, int, error) {
	var countUser int64
	var countUserToday int64

	user := models.User{}
	err := r.db.Unscoped().Where("role = 'user'").Order("id DESC").Find(&user).Count(&countUser).Error
	if err != nil {
		return  0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, err
	}

	newUser := []models.User{}
	err = r.db.Unscoped().Where("role = 'user'").Order("id DESC").Limit(10).Find(&newUser).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, err
	}

	// Get the start and end of the current day
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Count the number of users created today
	err = r.db.Model(&models.User{}).Unscoped().Where("role = 'user' AND created_at BETWEEN ? AND ?", startOfDay, endOfDay).Count(&countUserToday).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, err
	}



	var countLab int64
	lab := models.Lab{}
	err = r.db.Unscoped().Find(&lab).Count(&countLab).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, err
	}



	peminjaman := []models.Peminjaman{}

	var countPeminjaman int64
	var countPeminjamanToday int64
	err = r.db.Find(&peminjaman).Count(&countPeminjaman).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, err
	}
	newPeminjaman := []models.Peminjaman{}
	err = r.db.Order("id DESC").Limit(10).Find(&newPeminjaman).Error

	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, err
	}

	// Count the number of users created today
	err = r.db.Model(&models.Peminjaman{}).Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).Count(&countPeminjamanToday).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, err
	}



	jadwal := []models.Jadwal{}

	var countJadwal int64
	var countJadwalToday int64
	err = r.db.Find(&jadwal).Count(&countJadwal).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, err
	}
	newJadwal := []models.Jadwal{}
	err = r.db.Order("id DESC").Limit(10).Find(&newJadwal).Error

	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, err
	}

	// Count the number of users created today
	err = r.db.Model(&models.Jadwal{}).Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).Count(&countJadwalToday).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, err
	}

	return int(countUser), int(countUserToday), int(countLab),  int(countPeminjaman), int(countPeminjamanToday), newPeminjaman, newUser, newJadwal, int(countJadwal), int(countJadwalToday), nil
}

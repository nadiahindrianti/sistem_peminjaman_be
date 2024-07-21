package repositories

import (
	"sistem_peminjaman_be/models"
	"sistem_peminjaman_be/dtos"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	DashboardGetAll() (int, int, int, int, int, []models.Peminjaman, []models.User, []models.Jadwal, int, int, []dtos.UserTeraktifMeminjam, error)
	DashboardGetByMonth(month, year int) (int, int, int, int, int, int, int, int, int, int, error)
	
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db}
}

func (r *dashboardRepository) DashboardGetAll() (int, int, int, int, int, []models.Peminjaman, []models.User, []models.Jadwal, int, int, []dtos.UserTeraktifMeminjam, error) {
	var countUser int64
	var countUserToday int64
	var userTeraktif []dtos.UserTeraktifMeminjam

	user := models.User{}
	err := r.db.Unscoped().Where("role = 'user'").Order("id DESC").Find(&user).Count(&countUser).Error
	if err != nil {
		return  0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, nil, err
	}

	newUser := []models.User{}
	err = r.db.Unscoped().Where("role = 'user'").Order("id DESC").Limit(10).Find(&newUser).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, nil, err
	}

	// Get the start and end of the current day
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Count the number of users created today
	err = r.db.Model(&models.User{}).Unscoped().Where("role = 'user' AND created_at BETWEEN ? AND ?", startOfDay, endOfDay).Count(&countUserToday).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, nil, err
	}



	var countLab int64
	lab := models.Lab{}
	err = r.db.Unscoped().Find(&lab).Count(&countLab).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, nil, err
	}



	peminjaman := []models.Peminjaman{}

	var countPeminjaman int64
	var countPeminjamanToday int64
	err = r.db.Find(&peminjaman).Count(&countPeminjaman).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, nil, err
	}
	newPeminjaman := []models.Peminjaman{}
	err = r.db.Order("id DESC").Limit(10).Find(&newPeminjaman).Error

	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, nil, err
	}

	// Count the number of users created today
	err = r.db.Model(&models.Peminjaman{}).Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).Count(&countPeminjamanToday).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, nil, err
	}



	jadwal := []models.Jadwal{}

	var countJadwal int64
	var countJadwalToday int64
	err = r.db.Find(&jadwal).Count(&countJadwal).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, nil, err
	}
	newJadwal := []models.Jadwal{}
	err = r.db.Order("id DESC").Limit(10).Find(&newJadwal).Error

	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, nil, err
	}

	// Count the number of users created today
	err = r.db.Model(&models.Jadwal{}).Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).Count(&countJadwalToday).Error
	if err != nil {
		return 0, 0, 0, 0, 0, []models.Peminjaman{}, []models.User{}, []models.Jadwal{}, 0, 0, nil, err
	}

	// ntuk mendapatkan user teraktif meminjam (top 3)
	r.db.Raw(`
		SELECT u.full_name, COUNT(p.id) AS jumlah_peminjaman
		FROM peminjamen p
		JOIN users u ON p.user_id = u.id
		GROUP BY u.full_name
		ORDER BY jumlah_peminjaman DESC
		LIMIT 3
	`).Scan(&userTeraktif)

	return int(countUser), int(countUserToday), int(countLab),  int(countPeminjaman), int(countPeminjamanToday), newPeminjaman, newUser, newJadwal, int(countJadwal), int(countJadwalToday), userTeraktif, nil
}



func (r *dashboardRepository) DashboardGetByMonth(month, year int) (int, int, int, int, int, int, int, int, int, int, error) {
	var countUser int64
	var countLab int64
	var countPeminjaman int64
	var countJadwal int64

	// Count total users
	err := r.db.Model(&models.User{}).Where("role = 'user' AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", month, year).Count(&countUser).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	// Count total labs
	err = r.db.Model(&models.Lab{}).Count(&countLab).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	// Count total peminjaman
	err = r.db.Model(&models.Peminjaman{}).Where("EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", month, year).Count(&countPeminjaman).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	// Count total jadwal
	err = r.db.Model(&models.Jadwal{}).Where("EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", month, year).Count(&countJadwal).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	// Count peminjaman per lab
	var countPeminjamanLabElektronika int64
	var countPeminjamanLabTransmisi int64
	var countPeminjamanLabJaringan int64

	err = r.db.Model(&models.Peminjaman{}).Where("lab_id = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", 1, month, year).Count(&countPeminjamanLabElektronika).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	err = r.db.Model(&models.Peminjaman{}).Where("lab_id = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", 2, month, year).Count(&countPeminjamanLabTransmisi).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	err = r.db.Model(&models.Peminjaman{}).Where("lab_id = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", 3, month, year).Count(&countPeminjamanLabJaringan).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	// Count jadwal per lab
	var countJadwalLabElektronika int64
	var countJadwalLabTransmisi int64
	var countJadwalLabJaringan int64

	err = r.db.Model(&models.Jadwal{}).Where("name_laboratorium = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", "Laboratorium Elektronika", month, year).Count(&countJadwalLabElektronika).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	err = r.db.Model(&models.Jadwal{}).Where("name_laboratorium = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", "Laboratorium Transmisi", month, year).Count(&countJadwalLabTransmisi).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	err = r.db.Model(&models.Jadwal{}).Where("name_laboratorium = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", "Laboratorium Jaringan", month, year).Count(&countJadwalLabJaringan).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	return int(countUser), int(countLab), int(countPeminjaman), int(countJadwal), int(countPeminjamanLabElektronika), int(countPeminjamanLabTransmisi), int(countPeminjamanLabJaringan), int(countJadwalLabElektronika), int(countJadwalLabTransmisi), int(countJadwalLabJaringan), nil
}


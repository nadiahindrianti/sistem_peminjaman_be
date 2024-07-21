package usecases

import (
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/repositories"
	"time"

	"github.com/labstack/echo/v4"
)

type DashboardUsecase interface {
	DashboardGetAll() (dtos.DashboardResponse, error)
	DashboardGetByMonth(month, year int) (dtos.DashboardfilterResponse, error)
}

type dashboardUsecase struct {
	dashboardRepository repositories.DashboardRepository
	userRepo            repositories.UserRepository
	peminjamanRepo          repositories.PeminjamanRepository
	jadwalRepo           repositories.JadwalRepository
	labRepo                repositories.LabRepository
}

func NewDashboardUsecase(dashboardRepository repositories.DashboardRepository, userRepo repositories.UserRepository, peminjamanRepo  repositories.PeminjamanRepository, jadwalRepo repositories.JadwalRepository, labRepo repositories.LabRepository) DashboardUsecase {
	return &dashboardUsecase{dashboardRepository, userRepo, peminjamanRepo, jadwalRepo, labRepo}
}


func (u *dashboardUsecase) DashboardGetAll() (dtos.DashboardResponse, error) {
	var dashboardResponse dtos.DashboardResponse
	countUser, countUserToday, countLab,  countPeminjaman, countPeminjamanToday, newPeminjaman, newUser, newJadwal, countJadwal, countJadwalToday, userTeraktif, err := u.dashboardRepository.DashboardGetAll()
	if err != nil {
		return dashboardResponse, err
	}
	var newPeminjamanResponses []map[string]interface{}
	for _, peminjaman := range newPeminjaman {

		getPeminjaman, err := u.peminjamanRepo.GetPeminjamanID2(peminjaman.ID, 1)
		if err != nil {
			return dashboardResponse, err
		}

		getLab, err := u.labRepo.GetLabByID(uint(peminjaman.LabID))
		if err != nil {
			return dashboardResponse, err
		}

		newPeminjamanResponse := map[string]interface{}{
			"id":             getPeminjaman.ID,
			"peminjaman_name":     getLab.Name,
			"kegiatan":           "Peminjaman",
			"created_at":     getPeminjaman.CreatedAt,
			"updated_at":     getPeminjaman.UpdatedAt,
			"deleted_at":     getPeminjaman.DeletedAt,
		}
		newPeminjamanResponses = append(newPeminjamanResponses, newPeminjamanResponse)
	}

	var newJadwalResponses []map[string]interface{}
	for _, jadwal := range newJadwal {

		getJadwal, err := u.jadwalRepo.GetJadwalByID3(jadwal.ID, 1)
		if err != nil {
			return dashboardResponse, err
		}

		newJadwalResponse := map[string]interface{}{
			"id":             getJadwal.ID,
			"name_laboratorium":     getJadwal.NameLaboratorium,
			"kegiatan":           "Jadwal Peminjaman",
			"created_at":     getJadwal.CreatedAt,
			"updated_at":     getJadwal.UpdatedAt,
			"deleted_at":     getJadwal.DeletedAt,
		}
		newJadwalResponses = append(newJadwalResponses, newJadwalResponse)
	}

	var newUserResponses []map[string]interface{}
	for _, user := range newUser {
		newUserResponse := map[string]interface{}{
			"id":              user.ID,
			"full_name":       user.FullName,
			"nim_nip":			   user.NIMNIP,
			"profile_picture": user.ProfilePicture,
			"created_at":      user.CreatedAt,
			"updated_at":      user.UpdatedAt,
			"deleted_at":      user.DeletedAt,
		}
		newUserResponses = append(newUserResponses, newUserResponse)
	}

	
	for i := 0; i < len(newPeminjamanResponses)-1; i++ {
		for j := 0; j < len(newPeminjamanResponses)-i-1; j++ {
			timeI := newPeminjamanResponses[j]["created_at"].(time.Time)
			timeJ := newPeminjamanResponses[j+1]["created_at"].(time.Time)
			if timeI.Before(timeJ) {
				newPeminjamanResponses[j], newPeminjamanResponses[j+1] = newPeminjamanResponses[j+1],newPeminjamanResponses[j]
			}
		}
	}

	for i := 0; i < len(newJadwalResponses)-1; i++ {
		for j := 0; j < len(newJadwalResponses)-i-1; j++ {
			timeI := newJadwalResponses[j]["created_at"].(time.Time)
			timeJ := newJadwalResponses[j+1]["created_at"].(time.Time)
			if timeI.Before(timeJ) {
				newJadwalResponses[j], newJadwalResponses[j+1] = newJadwalResponses[j+1],newJadwalResponses[j]
			}
		}
	}

	// Batasi menjadi 10 data terbaru
	var limitedPeminjamanResponses []map[string]interface{}
	if len(newPeminjamanResponses) > 10 {
		limitedPeminjamanResponses = newPeminjamanResponses[:10]
	} else {
		limitedPeminjamanResponses = newPeminjamanResponses
	}

	var limitedJadwalResponses []map[string]interface{}
	if len(newJadwalResponses) > 10 {
		limitedJadwalResponses = newJadwalResponses[:10]
	} else {
		limitedJadwalResponses = newJadwalResponses
	}

	var userTeraktifResponses []dtos.UserTeraktifMeminjam
	for _, user := range userTeraktif {
		userTeraktifResponse := dtos.UserTeraktifMeminjam{
			FullName:         user.FullName,
			JumlahPeminjaman: user.JumlahPeminjaman,
		}
		userTeraktifResponses = append(userTeraktifResponses, userTeraktifResponse)
	}

	dashboardResponse = dtos.DashboardResponse{
		CountUser: echo.Map{
			"total_user":       countUser,
			"total_user_today": countUserToday,
		},
		CountLab: echo.Map{
			"total_lab": countLab,
		},
		
		CountPeminjaman: echo.Map{
			"total_peminjaman":       countPeminjaman,
			"total_peminjaman_today": countPeminjamanToday,
		},

		CountJadwal: echo.Map{
			"total_jadwal":       countJadwal,
			"total_jadwal_today": countJadwalToday,
		},

		NewPeminjaman: limitedPeminjamanResponses,
		NewJadwal:           limitedJadwalResponses,
		NewUser:             newUserResponses,
		UserTeraktifMeminjam: userTeraktifResponses,
	}
	return dashboardResponse, nil
}



func (u *dashboardUsecase) DashboardGetByMonth(month, year int) (dtos.DashboardfilterResponse, error) {
	var dashboardResponse dtos.DashboardfilterResponse
	countUser, countLab, countPeminjaman, countJadwal, countPeminjamanLabElektronika, countPeminjamanLabTransmisi, countPeminjamanLabJaringan, countJadwalLabElektronika, countJadwalLabTransmisi, countJadwalLabJaringan, err := u.dashboardRepository.DashboardGetByMonth(month, year)
	if err != nil {
		return dashboardResponse, err
	}

	dashboardResponse = dtos.DashboardfilterResponse{
		CountUser: echo.Map{
			"total_user": countUser,
		},
		CountLab: echo.Map{
			"total_lab": countLab,
		},
		CountPeminjaman: echo.Map{
			"total_peminjaman":               countPeminjaman,
			"total_peminjaman_lab_elektronika": countPeminjamanLabElektronika,
			"total_peminjaman_lab_transmisi":   countPeminjamanLabTransmisi,
			"total_peminjaman_lab_jaringan":    countPeminjamanLabJaringan,
		},
		CountJadwal: echo.Map{
			"total_jadwal":               countJadwal,
			"total_jadwal_lab_elektronika": countJadwalLabElektronika,
			"total_jadwal_lab_transmisi":   countJadwalLabTransmisi,
			"total_jadwal_lab_jaringan":    countJadwalLabJaringan,
		},
	}

	return dashboardResponse, nil
}

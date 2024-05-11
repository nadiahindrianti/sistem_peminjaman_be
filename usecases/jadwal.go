package usecases

import (
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/models"
	"sistem_peminjaman_be/repositories"
	"sistem_peminjaman_be/helpers"
	"time"
	"errors"
)

type JadwalUsecase interface {
	GetAllJadwals(page, limit int, name_laboratirum string) ([]dtos.JadwalResponse, int, error)
	GetJadwalByID(userId, id uint) (dtos.JadwalResponse, error)
	CreateJadwal(jadwal *dtos.JadwalInput) (dtos.JadwalResponse, error)
	UpdateJadwal(id uint, jadwalInput dtos.JadwalInput) (dtos.JadwalResponse, error)
	DeleteJadwal(id uint) error

	SearchJadwalAvailable(userId, page, limit int, name_laboratorium string) ([]dtos.JadwalResponse, int, error)
}

type jadwalUsecase struct {
	jadwalRepo               	  repositories.JadwalRepository
	beritaAcaraImageRepo          repositories.BeritaAcaraImageRepository
	userRepo                      repositories.UserRepository
}

func NewJadwalUsecase(jadwalRepo repositories.JadwalRepository, beritaAcaraImageRepo repositories.BeritaAcaraImageRepository, userRepo repositories.UserRepository) JadwalUsecase {
	return &jadwalUsecase{jadwalRepo, beritaAcaraImageRepo, userRepo}
}


func (u *jadwalUsecase) GetAllJadwals(page, limit int, name_laboratirum string) ([]dtos.JadwalResponse, int, error) {
	jadwals, count, err := u.jadwalRepo.SearchJadwalAvailable(page, limit, name_laboratirum)
	if err != nil {
		return nil, 0, err
	}

	var jadwalResponses []dtos.JadwalResponse

	for _, jadwal := range jadwals {
		getBeritaAcaraImage, err := u.beritaAcaraImageRepo.GetAllBeritaAcaraImageByID(jadwal.ID)
		if err != nil {
			continue
		}

		var beritaAcaraImageResponses []dtos.BeritaAcaraImageResponse
		for _, beritaacaraimage := range getBeritaAcaraImage {
			beritaAcaraImageResponse := dtos.BeritaAcaraImageResponse{
				JadwalID:  beritaacaraimage.JadwalID,
				BeritaAcaraImageUrl: beritaacaraimage.BeritaAcaraImageUrl,
			}
			beritaAcaraImageResponses = append(beritaAcaraImageResponses, beritaAcaraImageResponse)
		}

		jadwalResponse := dtos.JadwalResponse{
			JadwalID:           jadwal.ID,
			TanggalJadwal:      helpers.FormatDateToYMD(jadwal.TanggalJadwal),
			WaktuJadwal:        jadwal.WaktuJadwal,
			NameUser:           jadwal.NameUser,
			NameLaboratorium:   jadwal.NameLaboratorium,
			BeritaAcaraImage:   beritaAcaraImageResponses,
			Status:             jadwal.Status,
			CreatedAt:       	jadwal.CreatedAt,
			UpdatedAt:       	jadwal.UpdatedAt,
		}

		jadwalResponses = append(jadwalResponses, jadwalResponse)

	}

	return jadwalResponses, count, nil
}

func (u *jadwalUsecase) GetJadwalByID(userId, id uint) (dtos.JadwalResponse, error) {
	var jadwalResponses dtos.JadwalResponse
	jadwal, err := u.jadwalRepo.GetJadwalByID(id)
	if err != nil {
		return jadwalResponses, err
	}

	getBeritaAcaraImage, err := u.beritaAcaraImageRepo.GetAllBeritaAcaraImageByID(jadwal.ID)
	if err != nil {
		return jadwalResponses, err
	}

	var beritaAcaraImageResponses []dtos.BeritaAcaraImageResponse
	for _, beritaacaraimage := range getBeritaAcaraImage {
		beritaAcaraImageResponse := dtos.BeritaAcaraImageResponse{
			JadwalID:    beritaacaraimage.JadwalID,
			BeritaAcaraImageUrl: beritaacaraimage.BeritaAcaraImageUrl,
		}
		beritaAcaraImageResponses = append(beritaAcaraImageResponses, beritaAcaraImageResponse)
	}

	jadwalResponse := dtos.JadwalResponse{
		JadwalID:           jadwal.ID,
		TanggalJadwal:      helpers.FormatDateToYMD(jadwal.TanggalJadwal),
		WaktuJadwal:        jadwal.WaktuJadwal,
		NameUser:           jadwal.NameUser,
		NameLaboratorium:   jadwal.NameLaboratorium,
		BeritaAcaraImage:   beritaAcaraImageResponses,
		Status:             jadwal.Status,
		CreatedAt:          jadwal.CreatedAt,
		UpdatedAt:          jadwal.UpdatedAt,
	}
	return jadwalResponse, nil
}


func (u *jadwalUsecase) CreateJadwal(jadwal *dtos.JadwalInput) (dtos.JadwalResponse, error) {
	var jadwalResponse dtos.JadwalResponse

	// Cek apakah tanggal jadwal valid
	if *jadwal.TanggalJadwal < time.Now().Format("2006-01-02") {
		return jadwalResponse, errors.New("Tanggal Jadwal invalid")
	}

	// Mengambil tanggal sekarang
	dateNow := time.Now()

	// Parsing tanggal jadwal
	tanggalJadwalParse, err := time.Parse("2006-01-02", *jadwal.TanggalJadwal)
	if err != nil {
		return jadwalResponse, errors.New("Failed to parse Tanggal Jadwal")
	}

	// Membandingkan tanggal jadwal dengan tanggal sekarang
	if tanggalJadwalParse.Before(dateNow) {
		return jadwalResponse, errors.New("Tanggal Jadwal harus setelah tanggal sekarang")
	}

	// Membuat struktur Jadwal dari data input
	createJadwal := models.Jadwal{
		TanggalJadwal:    &tanggalJadwalParse,
		WaktuJadwal:      jadwal.WaktuJadwal,
		NameUser:         jadwal.NameUser,
		NameLaboratorium: jadwal.NameLaboratorium,
		Status:           "notused",
	}

	// Menyimpan data jadwal ke repository
	createdJadwal, err := u.jadwalRepo.CreateJadwal(createJadwal)
	if err != nil {
		return jadwalResponse, err
	}


	// Memproses gambar berita acara jika ada
	for _, beritaAcaraImage := range jadwal.BeritaAcaraImage {
		if beritaAcaraImage.BeritaAcaraImageUrl == "" {
			return jadwalResponse, errors.New("failed to create jadwal: berita acara image URL is empty")
		}
		jadwalImage := models.BeritaAcaraImage{
			JadwalID:            createdJadwal.ID,
			BeritaAcaraImageUrl: beritaAcaraImage.BeritaAcaraImageUrl,
		}
		_, err := u.beritaAcaraImageRepo.CreateBeritaAcaraImage(jadwalImage)
		if err != nil {
			return jadwalResponse, err
		}
	}

	// Mengambil semua gambar berita acara untuk jadwal yang baru saja dibuat
	getBeritaAcaraImage, err := u.beritaAcaraImageRepo.GetAllBeritaAcaraImageByID(createdJadwal.ID)
	if err != nil {
		return jadwalResponse, err
	}

	// Mengonversi data gambar berita acara ke DTO
	var beritaAcaraImageResponses []dtos.BeritaAcaraImageResponse
	for _, beritaAcaraImage := range getBeritaAcaraImage {
		beritaAcaraImageResponse := dtos.BeritaAcaraImageResponse{
			JadwalID:            beritaAcaraImage.JadwalID,
			BeritaAcaraImageUrl: beritaAcaraImage.BeritaAcaraImageUrl,
		}
		beritaAcaraImageResponses = append(beritaAcaraImageResponses, beritaAcaraImageResponse)
	}

	// Menyiapkan respons jadwal yang akan dikembalikan
	jadwalResponse = dtos.JadwalResponse{
		JadwalID:           createdJadwal.ID,
		TanggalJadwal:      helpers.FormatDateToYMD(createdJadwal.TanggalJadwal),
		WaktuJadwal:        createdJadwal.WaktuJadwal,
		NameUser:           createdJadwal.NameUser,
		NameLaboratorium:   createdJadwal.NameLaboratorium,
		BeritaAcaraImage:   beritaAcaraImageResponses,
		Status:             createdJadwal.Status,
		CreatedAt:          createdJadwal.CreatedAt,
		UpdatedAt:          createdJadwal.UpdatedAt,
	}

	return jadwalResponse, nil
}



func (u *jadwalUsecase) UpdateJadwal(id uint, jadwal dtos.JadwalInput) (dtos.JadwalResponse, error) {
	var jadwals models.Jadwal
	var jadwalResponse dtos.JadwalResponse

	jadwals, err := u.jadwalRepo.GetJadwalByID(id)
	if err != nil {
		return jadwalResponse, err
	}

	if jadwal.WaktuJadwal == "" || jadwal.NameUser == "" || jadwal.NameLaboratorium == "" || jadwal.BeritaAcaraImage == nil || jadwal.Status == "" {
		return jadwalResponse, errors.New("failed to update jadwal")
	}

	if *jadwal.TanggalJadwal > time.Now().Format("2006-01-02") {
		return jadwalResponse, errors.New("Tanggal Jadwal invalid")
	}

	dateNow := "2006-01-02"
	tanggalJadwalParse, err := time.Parse(dateNow, *jadwal.TanggalJadwal)
	if err != nil {
		return jadwalResponse, errors.New("Failed to parse tanggal jadwal")
	}


	jadwals.TanggalJadwal    = &tanggalJadwalParse
	jadwals.WaktuJadwal      =  jadwal.WaktuJadwal
	jadwals.NameUser         = jadwal.NameUser
	jadwals.NameLaboratorium = jadwal.NameLaboratorium
	jadwals.Status           = jadwal.Status

	updatedJadwal, err := u.jadwalRepo.UpdateJadwal(jadwals)
	if err != nil {
		return jadwalResponse, err
	}

	u.beritaAcaraImageRepo.DeleteBeritaAcaraImage(id)

	for _, beritaAcaraImage := range jadwal.BeritaAcaraImage {
		if beritaAcaraImage.BeritaAcaraImageUrl == "" {
			return jadwalResponse, errors.New("failed to update jadwal")
		}
		beritaAcaraImagee := models.BeritaAcaraImage{
			JadwalID:  updatedJadwal.ID,
			BeritaAcaraImageUrl: beritaAcaraImage.BeritaAcaraImageUrl,
		}
		_, err = u.beritaAcaraImageRepo.UpdateBeritaAcaraImage(beritaAcaraImagee)
		if err != nil {
			return jadwalResponse, err
		}
	}

	getBeritaAcaraImage, err := u.beritaAcaraImageRepo.GetAllBeritaAcaraImageByID(updatedJadwal.ID)
	if err != nil {
		return jadwalResponse, err
	}

	var beritaAcaraImageResponses []dtos.BeritaAcaraImageResponse
	for _, beritaAcaraimage := range getBeritaAcaraImage {
		beritaAcaraImageResponse := dtos.BeritaAcaraImageResponse{
			JadwalID:            beritaAcaraimage.JadwalID,
			BeritaAcaraImageUrl: beritaAcaraimage.BeritaAcaraImageUrl,
		}
		beritaAcaraImageResponses = append(beritaAcaraImageResponses, beritaAcaraImageResponse)
	}

	jadwalResponse = dtos.JadwalResponse{
		JadwalID:           updatedJadwal.ID,
		TanggalJadwal:      helpers.FormatDateToYMD(updatedJadwal.TanggalJadwal),
		WaktuJadwal:        updatedJadwal.WaktuJadwal,
		NameUser:           updatedJadwal.NameUser,
		NameLaboratorium:   updatedJadwal.NameLaboratorium,
		BeritaAcaraImage:   beritaAcaraImageResponses,
		Status:             updatedJadwal.Status,
		CreatedAt:       	updatedJadwal.CreatedAt,
		UpdatedAt:       	updatedJadwal.UpdatedAt,
	}
	return jadwalResponse, nil
}


func (u *jadwalUsecase) DeleteJadwal(id uint) error {
	return u.jadwalRepo.DeleteJadwal(id)
}


func (u *jadwalUsecase) SearchJadwalAvailable(userId, page, limit int, name_laboratorium string) ([]dtos.JadwalResponse, int, error) {
	jadwals, count, err := u.jadwalRepo.SearchJadwalAvailable(page, limit, name_laboratorium)
	if err != nil {
		return nil, 0, err
	}

	var jadwalResponses []dtos.JadwalResponse

	for _, jadwal := range jadwals {
		getBeritaAcaraImage, err := u.beritaAcaraImageRepo.GetAllBeritaAcaraImageByID(jadwal.ID)
		if err != nil {
			continue
		}

		var beritaAcaraImageResponses []dtos.BeritaAcaraImageResponse
		for _, beritaAcaraimage := range getBeritaAcaraImage {
			beritaAcaraImageResponse := dtos.BeritaAcaraImageResponse{
				JadwalID:    beritaAcaraimage.JadwalID,
				BeritaAcaraImageUrl: beritaAcaraimage.BeritaAcaraImageUrl,
			}
			beritaAcaraImageResponses = append(beritaAcaraImageResponses, beritaAcaraImageResponse)
		}

		jadwalResponse := dtos.JadwalResponse{
			JadwalID:           jadwal.ID,
			TanggalJadwal:      helpers.FormatDateToYMD(jadwal.TanggalJadwal),
			WaktuJadwal:        jadwal.WaktuJadwal,
			NameUser:           jadwal.NameUser,
			NameLaboratorium:   jadwal.NameLaboratorium,
			BeritaAcaraImage:   beritaAcaraImageResponses,
			Status:             jadwal.Status,
			CreatedAt:       	jadwal.CreatedAt,
			UpdatedAt:       	jadwal.UpdatedAt,
		}

		jadwalResponses = append(jadwalResponses, jadwalResponse)


	}

	return jadwalResponses, count, nil

}
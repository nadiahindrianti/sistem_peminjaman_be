package usecases

import (
	"errors"
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/helpers"
	"sistem_peminjaman_be/models"
	"sistem_peminjaman_be/repositories"
	"strings"
	"time"
)

type PeminjamanUsecase interface {
	GetPeminjamans(page, limit int, userID uint, nameLaboratorium, status string) ([]dtos.PeminjamanResponse, int, error)
	GetPeminjamanByID(userId, peminjamanId uint) (dtos.PeminjamanResponse, error)
	GetPeminjamansByAdmin(page, limit int, search, status string) ([]dtos.PeminjamanResponse, int, error)
	GetPeminjamansDetailByAdmin(peminjamanId uint) (dtos.PeminjamanResponse, error)
	CreatePeminjaman(userID uint, peminjaman *dtos.PeminjamanInput) (dtos.PeminjamanResponse, error)
	UpdatePeminjaman(id uint, userID uint, peminjaman dtos.PeminjamanInput) (dtos.PeminjamanResponse, error)
	DeletePeminjaman(id uint) error
}

type peminjamanUsecase struct {
	peminjamanRepo            repositories.PeminjamanRepository
	suratRekomendasiImageRepo repositories.SuratRekomendasiImageRepository
	labRepo                   repositories.LabRepository
	labImageRepo              repositories.LabImageRepository
	userRepo                  repositories.UserRepository
}

func NewPeminjamanUsecase(peminjamanRepo repositories.PeminjamanRepository, suratRekomendasiImageRepo repositories.SuratRekomendasiImageRepository, labRepo repositories.LabRepository, labImageRepo repositories.LabImageRepository, userRepo repositories.UserRepository) PeminjamanUsecase {
	return &peminjamanUsecase{peminjamanRepo, suratRekomendasiImageRepo, labRepo, labImageRepo, userRepo}
}

func (u *peminjamanUsecase) GetPeminjamans(page, limit int, userID uint, nameLaboratorium, status string) ([]dtos.PeminjamanResponse, int, error) {
	var peminjamanResponses []dtos.PeminjamanResponse

	peminjamans, _, err := u.peminjamanRepo.GetPeminjamans(page, limit, userID, status)
	if err != nil {
		return peminjamanResponses, 0, err
	}

	for _, peminjaman := range peminjamans {

		getSuratRekomendasiImage, err := u.suratRekomendasiImageRepo.GetAllSuratRekomendasiImageByID(peminjaman.ID)
		if err != nil {
			continue
		}

		var suratRekomendasiImageResponses []dtos.SuratRekomendasiImageResponse
		for _, suratRekomendasiImage := range getSuratRekomendasiImage {
			suratRekomendasiImageResponse := dtos.SuratRekomendasiImageResponse{
				PeminjamanID:             suratRekomendasiImage.PeminjamanID,
				SuratRekomendasiImageUrl: suratRekomendasiImage.SuratRekomendasiImageUrl,
			}
			suratRekomendasiImageResponses = append(suratRekomendasiImageResponses, suratRekomendasiImageResponse)
		}

		getLab, err := u.labRepo.GetLabByID2(peminjaman.LabID)
		if err != nil {
			return peminjamanResponses, 0, err
		}

		if nameLaboratorium != "" && !strings.Contains(strings.ToLower(getLab.Name), strings.ToLower(nameLaboratorium)) {
			continue
		}

		getLabImage, err := u.labImageRepo.GetAllLabImageByID(peminjaman.LabID)
		if err != nil {
			continue
		}

		var labImageResponses []dtos.LabImageResponse
		for _, labImage := range getLabImage {
			labImageResponse := dtos.LabImageResponse{
				LabID:    labImage.LabID,
				ImageUrl: labImage.ImageUrl,
			}
			labImageResponses = append(labImageResponses, labImageResponse)
		}

		peminjamanResponse := dtos.PeminjamanResponse{
			PeminjamanID:          int(peminjaman.ID),
			TanggalPeminjaman:     helpers.FormatDateToYMD(peminjaman.TanggalPeminjaman),
			JamPeminjaman:         peminjaman.JamPeminjaman,
			SuratRekomendasiImage: suratRekomendasiImageResponses,
			Description:           peminjaman.Description,
			Status:                peminjaman.Status,
			Lab: dtos.LabByIDResponses{
				LabID:       getLab.ID,
				Name:        getLab.Name,
				LabImage:    labImageResponses,
				Description: getLab.Description,
			},
			CreatedAt: peminjaman.CreatedAt,
			UpdatedAt: peminjaman.UpdatedAt,
		}
		peminjamanResponses = append(peminjamanResponses, peminjamanResponse)
	}

	return peminjamanResponses, len(peminjamanResponses), nil
}

func (u *peminjamanUsecase) GetPeminjamanByID(userId, peminjamanId uint) (dtos.PeminjamanResponse, error) {
	peminjamanResponses := dtos.PeminjamanResponse{}

	// Mendapatkan data peminjaman berdasarkan ID
	peminjaman, err := u.peminjamanRepo.GetPeminjamanByID(userId, peminjamanId)
	if err != nil {
		return peminjamanResponses, err
	}

	// Mendapatkan data surat rekomendasi image berdasarkan ID peminjaman
	suratRekomendasiImagee, err := u.suratRekomendasiImageRepo.GetAllSuratRekomendasiImageByID(peminjamanId)
	if err != nil {
		return peminjamanResponses, err
	}

	var suratRekomendasiImageResponses []dtos.SuratRekomendasiImageResponse
	for _, suratRekomendasiImageee := range suratRekomendasiImagee {
		suratRekomendasiImageResponse := dtos.SuratRekomendasiImageResponse{
			PeminjamanID:             suratRekomendasiImageee.PeminjamanID,
			SuratRekomendasiImageUrl: suratRekomendasiImageee.SuratRekomendasiImageUrl,
		}
		suratRekomendasiImageResponses = append(suratRekomendasiImageResponses, suratRekomendasiImageResponse)
	}

	// Mendapatkan data lab berdasarkan ID
	getLab, err := u.labRepo.GetLabByID2(peminjaman.LabID)
	if err != nil {
		return peminjamanResponses, err
	}

	// Mendapatkan data lab image berdasarkan ID lab
	getLabImage, err := u.labImageRepo.GetAllLabImageByID(peminjaman.LabID)
	if err != nil {
		return peminjamanResponses, err
	}

	var labImageResponses []dtos.LabImageResponse
	for _, labImage := range getLabImage {
		labImageResponse := dtos.LabImageResponse{
			LabID:    labImage.LabID,
			ImageUrl: labImage.ImageUrl,
		}
		labImageResponses = append(labImageResponses, labImageResponse)
	}

	// Membuat respons peminjaman
	peminjamanResponse := dtos.PeminjamanResponse{
		PeminjamanID:          int(peminjaman.ID),
		TanggalPeminjaman:     helpers.FormatDateToYMD(peminjaman.TanggalPeminjaman),
		JamPeminjaman:         peminjaman.JamPeminjaman,
		SuratRekomendasiImage: suratRekomendasiImageResponses,
		Description:           peminjaman.Description,
		Status:                peminjaman.Status,
		Lab: dtos.LabByIDResponses{
			LabID:       getLab.ID,
			Name:        getLab.Name,
			LabImage:    labImageResponses,
			Description: getLab.Description,
		},
		CreatedAt: peminjaman.CreatedAt,
		UpdatedAt: peminjaman.UpdatedAt,
	}

	return peminjamanResponse, nil
}

func (u *peminjamanUsecase) GetPeminjamansByAdmin(page, limit int, search, status string) ([]dtos.PeminjamanResponse, int, error) {
	var peminjamanResponses []dtos.PeminjamanResponse

	peminjamans, _, err := u.peminjamanRepo.GetPeminjamans(page, limit, 1, status)
	if err != nil {
		return peminjamanResponses, 0, err
	}

	for _, peminjaman := range peminjamans {

		getSuratRekomendasiImage, err := u.suratRekomendasiImageRepo.GetAllSuratRekomendasiImageByID(peminjaman.ID)
		if err != nil {
			continue
		}

		var suratRekomendasiImageResponses []dtos.SuratRekomendasiImageResponse
		for _, suratRekomendasiImage := range getSuratRekomendasiImage {
			suratRekomendasiImageResponse := dtos.SuratRekomendasiImageResponse{
				PeminjamanID:             suratRekomendasiImage.PeminjamanID,
				SuratRekomendasiImageUrl: suratRekomendasiImage.SuratRekomendasiImageUrl,
			}
			suratRekomendasiImageResponses = append(suratRekomendasiImageResponses, suratRekomendasiImageResponse)
		}

		getLab, err := u.labRepo.GetLabByID2(peminjaman.LabID)
		if err != nil {
			return peminjamanResponses, 0, err
		}

		if search != "" &&
			!strings.Contains(strings.ToLower(getLab.Name), strings.ToLower(search)) {
			continue // Skip hotel order if it doesn't match the search query
		}

		getUser, err := u.userRepo.UserGetById2(peminjaman.UserID)
		if err != nil {
			return peminjamanResponses, 0, err
		}

		getLabImage, err := u.labImageRepo.GetAllLabImageByID(peminjaman.LabID)
		if err != nil {
			continue
		}

		var labImageResponses []dtos.LabImageResponse
		for _, labImage := range getLabImage {
			labImageResponse := dtos.LabImageResponse{
				LabID:    labImage.LabID,
				ImageUrl: labImage.ImageUrl,
			}
			labImageResponses = append(labImageResponses, labImageResponse)
		}

		// Membuat respons peminjaman
		peminjamanResponse := dtos.PeminjamanResponse{
			PeminjamanID:          int(peminjaman.ID),
			TanggalPeminjaman:     helpers.FormatDateToYMD(peminjaman.TanggalPeminjaman),
			JamPeminjaman:         peminjaman.JamPeminjaman,
			SuratRekomendasiImage: suratRekomendasiImageResponses,
			Description:           peminjaman.Description,
			Status:                peminjaman.Status,
			Lab: dtos.LabByIDResponses{
				LabID:       getLab.ID,
				Name:        getLab.Name,
				LabImage:    labImageResponses,
				Description: getLab.Description,
			},
			User: &dtos.UserInformationResponses{
				ID:             getUser.ID,
				FullName:       getUser.FullName,
				Email:          getUser.Email,
				NIMNIP:         getUser.NIMNIP,
				ProfilePicture: getUser.ProfilePicture,
			},
			CreatedAt: peminjaman.CreatedAt,
			UpdatedAt: peminjaman.UpdatedAt,
		}

		peminjamanResponses = append(peminjamanResponses, peminjamanResponse)
	}

	return peminjamanResponses, len(peminjamanResponses), nil

}

func (u *peminjamanUsecase) GetPeminjamansDetailByAdmin(peminjamanId uint) (dtos.PeminjamanResponse, error) {
	var peminjamanResponses dtos.PeminjamanResponse

	// Mendapatkan data peminjaman berdasarkan ID
	peminjaman, err := u.peminjamanRepo.GetPeminjamanByID(peminjamanId, 1)
	if err != nil {
		return peminjamanResponses, err
	}

	//Mendapatkan data surat rekomendasi image berdasarkan ID peminjaman
	suratRekomendasiImagee, err := u.suratRekomendasiImageRepo.GetAllSuratRekomendasiImageByID(peminjamanId)
	if err != nil {
		return peminjamanResponses, err
	}

	var suratRekomendasiImageResponses []dtos.SuratRekomendasiImageResponse
	for _, suratRekomendasiImageee := range suratRekomendasiImagee {
		suratRekomendasiImageResponse := dtos.SuratRekomendasiImageResponse{
			PeminjamanID:             suratRekomendasiImageee.PeminjamanID,
			SuratRekomendasiImageUrl: suratRekomendasiImageee.SuratRekomendasiImageUrl,
		}
		suratRekomendasiImageResponses = append(suratRekomendasiImageResponses, suratRekomendasiImageResponse)
	}

	// Mendapatkan data lab berdasarkan ID
	getLab, err := u.labRepo.GetLabByID2(peminjaman.LabID)
	if err != nil {
		return peminjamanResponses, err
	}

	// Mendapatkan data lab image berdasarkan ID lab
	getLabImage, err := u.labImageRepo.GetAllLabImageByID(peminjaman.LabID)
	if err != nil {
		return peminjamanResponses, err
	}

	var labImageResponses []dtos.LabImageResponse
	for _, labImage := range getLabImage {
		labImageResponse := dtos.LabImageResponse{
			LabID:    labImage.LabID,
			ImageUrl: labImage.ImageUrl,
		}
		labImageResponses = append(labImageResponses, labImageResponse)
	}

	// Membuat respons peminjaman
	peminjamanResponse := dtos.PeminjamanResponse{
		PeminjamanID:          int(peminjaman.ID),
		TanggalPeminjaman:     helpers.FormatDateToYMD(peminjaman.TanggalPeminjaman),
		JamPeminjaman:         peminjaman.JamPeminjaman,
		SuratRekomendasiImage: suratRekomendasiImageResponses,
		Description:           peminjaman.Description,
		Status:                peminjaman.Status,
		Lab: dtos.LabByIDResponses{
			LabID:       getLab.ID,
			Name:        getLab.Name,
			LabImage:    labImageResponses,
			Description: getLab.Description,
		},
		CreatedAt: peminjaman.CreatedAt,
		UpdatedAt: peminjaman.UpdatedAt,
	}

	return peminjamanResponse, nil
}

// error func
func (u *peminjamanUsecase) CreatePeminjaman(userID uint, PeminjamanInput *dtos.PeminjamanInput) (dtos.PeminjamanResponse, error) {
	var peminjamanResponse dtos.PeminjamanResponse

	// Mendapatkan data lab berdasarkan ID yang diberikan di input
	getUsers, err := u.userRepo.UserGetById(peminjamanResponse.User.ID)
	if err != nil {
		return peminjamanResponse, err
	}

	// Mendapatkan data lab berdasarkan ID yang diberikan di input
	getLabs, err := u.labRepo.GetLabByID(peminjamanResponse.Lab.LabID)
	if err != nil {
		return peminjamanResponse, err
	}

	// Cek apakah tanggal peminjaman valid
	if *peminjaman.TanggalPeminjaman < time.Now().Format("2006-01-02") {
		return peminjamanResponse, errors.New("Tanggal Peminjaman invalid")
	}

	// Mengambil tanggal sekarang
	dateNow := time.Now()

	// Parsing tanggal peminjaman
	tanggalPeminjamanParse, err := time.Parse("2006-01-02", *peminjaman.TanggalPeminjaman)
	if err != nil {
		return peminjamanResponse, errors.New("Failed to parse Tanggal Peminjaman")
	}

	// Membandingkan tanggal peminjaman dengan tanggal sekarang
	if tanggalPeminjamanParse.Before(dateNow) {
		return peminjamanResponse, errors.New("Tanggal Peminjaman harus setelah tanggal sekarang")
	}

	// Membuat struktur Peminjaman dari data input
	createPeminjaman := models.Peminjaman{
		UserID:            getUsers.ID,
		LabID:             getLabs.ID,
		TanggalPeminjaman: &tanggalPeminjamanParse,
		JamPeminjaman:     peminjaman.JamPeminjaman,
		Description:       peminjaman.Description,
		Status:            "request",
	}

	// Menyimpan data peminjaman ke repository
	createdPeminjaman, err := u.peminjamanRepo.CreatePeminjaman(createPeminjaman)
	if err != nil {
		return peminjamanResponse, err
	}

	// Memproses gambar surat peminjaman jika ada
	for _, suratRekomendasiImage := range peminjaman.SuratRekomendasiImage {
		if suratRekomendasiImage.SuratRekomendasiImageUrl == "" {
			return peminjamanResponse, errors.New("failed to create peminjaman: surat rekomendasi image URL is empty")
		}
		peminjamanImage := models.SuratRekomendasiImage{
			PeminjamanID:             createdPeminjaman.ID,
			SuratRekomendasiImageUrl: suratRekomendasiImage.SuratRekomendasiImageUrl,
		}
		_, err := u.suratRekomendasiImageRepo.CreateSuratRekomendasiImage(peminjamanImage)
		if err != nil {
			return peminjamanResponse, err
		}
	}

	// Mengonversi data gambar surat rekomendasi ke DTO
	var suratRekomendasiImageResponses []dtos.SuratRekomendasiImageResponse
	for _, suratRekomendasiImage := range peminjaman.SuratRekomendasiImage {
		suratRekomendasiImageResponse := dtos.SuratRekomendasiImageResponse{
			PeminjamanID:             createdPeminjaman.ID,
			SuratRekomendasiImageUrl: suratRekomendasiImage.SuratRekomendasiImageUrl,
		}
		suratRekomendasiImageResponses = append(suratRekomendasiImageResponses, suratRekomendasiImageResponse)
	}

	peminjaman, err := u.peminjamanRepo.GetPeminjamanByID(createdPeminjaman.ID, userID)
	if err != nil {
		return peminjamanResponse, err
	}

	getLab, err := u.labRepo.GetLabByID(peminjaman.LabID)
	if err != nil {
		return peminjamanResponse, err
	}

	// Mendapatkan data lab image berdasarkan ID lab
	getLabImage, err := u.labImageRepo.GetAllLabImageByID(peminjaman.LabID)
	if err != nil {
		return peminjamanResponse, err
	}

	// Mengonversi data gambar lab image ke DTO
	var labImageResponses []dtos.LabImageResponse
	for _, labImage := range getLabImage {
		labImageResponse := dtos.LabImageResponse{
			LabID:    labImage.LabID,
			ImageUrl: labImage.ImageUrl,
		}
		labImageResponses = append(labImageResponses, labImageResponse)
	}

	// Membuat respons peminjaman
	peminjamanResponse = dtos.PeminjamanResponse{
		PeminjamanID:          int(peminjaman.ID),
		TanggalPeminjaman:     helpers.FormatDateToYMD(createdPeminjaman.TanggalPeminjaman),
		JamPeminjaman:         createPeminjaman.JamPeminjaman,
		SuratRekomendasiImage: suratRekomendasiImageResponses,
		Description:           createPeminjaman.Description,
		Status:                createPeminjaman.Status,
		Lab: dtos.LabByIDResponses{
			LabID:       getLabs.ID,
			Name:        getLabs.Name,
			LabImage:    labImageResponses,
			Description: getLabs.Description,
		},
		CreatedAt: createdPeminjaman.CreatedAt,
		UpdatedAt: createdPeminjaman.UpdatedAt,
	}

	return peminjamanResponse, nil
}

func (u *peminjamanUsecase) UpdatePeminjaman(id uint, userID uint, peminjaman dtos.PeminjamanInput) (dtos.PeminjamanResponse, error) {
	var peminjamans models.Peminjaman
	var peminjamanResponse dtos.PeminjamanResponse

	peminjamans, err := u.peminjamanRepo.GetPeminjamanByID(id, userID)
	if err != nil {
		return peminjamanResponse, err
	}

	if peminjaman.JamPeminjaman == "" || peminjaman.Description == "" || peminjaman.SuratRekomendasiImage == nil || peminjaman.Status == "" {
		return peminjamanResponse, errors.New("failed to update peminjaman")
	}

	if *peminjaman.TanggalPeminjaman > time.Now().Format("2006-01-02") {
		return peminjamanResponse, errors.New("Tanggal Peminjaman invalid")
	}

	dateNow := "2006-01-02"
	tanggalPeminjamanParse, err := time.Parse(dateNow, *peminjaman.TanggalPeminjaman)
	if err != nil {
		return peminjamanResponse, errors.New("Failed to parse tanggal peminjaman")
	}

	peminjamans.TanggalPeminjaman = &tanggalPeminjamanParse
	peminjamans.JamPeminjaman = peminjaman.JamPeminjaman
	peminjamans.Description = peminjaman.Description
	peminjamans.Status = peminjaman.Status

	updatedPeminjaman, err := u.peminjamanRepo.UpdatePeminjaman(peminjamans)
	if err != nil {
		return peminjamanResponse, err
	}

	u.suratRekomendasiImageRepo.DeleteSuratRekomendasiImage(id)

	for _, suratRekomendasiImage := range peminjaman.SuratRekomendasiImage {
		if suratRekomendasiImage.SuratRekomendasiImageUrl == "" {
			return peminjamanResponse, errors.New("failed to update peminjaman")
		}
		suratRekomendasiImagee := models.SuratRekomendasiImage{
			PeminjamanID:             updatedPeminjaman.ID,
			SuratRekomendasiImageUrl: suratRekomendasiImage.SuratRekomendasiImageUrl,
		}
		_, err = u.suratRekomendasiImageRepo.UpdateSuratRekomendasiImage(suratRekomendasiImagee)
		if err != nil {
			return peminjamanResponse, err
		}
	}

	getSuratRekomendasiImage, err := u.suratRekomendasiImageRepo.GetAllSuratRekomendasiImageByID(updatedPeminjaman.ID)
	if err != nil {
		return peminjamanResponse, err
	}

	var suratRekomendasiImageResponses []dtos.SuratRekomendasiImageResponse
	for _, suratRekomendasiimage := range getSuratRekomendasiImage {
		suratRekomendasiImageResponse := dtos.SuratRekomendasiImageResponse{
			PeminjamanID:             suratRekomendasiimage.PeminjamanID,
			SuratRekomendasiImageUrl: suratRekomendasiimage.SuratRekomendasiImageUrl,
		}
		suratRekomendasiImageResponses = append(suratRekomendasiImageResponses, suratRekomendasiImageResponse)
	}

	getLab, err := u.labRepo.GetLabByID2(peminjamans.LabID)
	if err != nil {
		return peminjamanResponse, err
	}

	// Mendapatkan data lab image berdasarkan ID lab
	getLabImage, err := u.labImageRepo.GetAllLabImageByID(peminjamans.LabID)
	if err != nil {
		return peminjamanResponse, err
	}

	// Mengonversi data gambar lab image ke DTO
	var labImageResponses []dtos.LabImageResponse
	for _, labImage := range getLabImage {
		labImageResponse := dtos.LabImageResponse{
			LabID:    labImage.LabID,
			ImageUrl: labImage.ImageUrl,
		}
		labImageResponses = append(labImageResponses, labImageResponse)
	}

	peminjamanResponse = dtos.PeminjamanResponse{
		PeminjamanID:          int(updatedPeminjaman.ID),
		TanggalPeminjaman:     helpers.FormatDateToYMD(updatedPeminjaman.TanggalPeminjaman),
		JamPeminjaman:         updatedPeminjaman.JamPeminjaman,
		SuratRekomendasiImage: suratRekomendasiImageResponses,
		Description:           updatedPeminjaman.Description,
		Status:                updatedPeminjaman.Status,
		Lab: dtos.LabByIDResponses{
			LabID:       getLab.ID,
			Name:        getLab.Name,
			LabImage:    labImageResponses,
			Description: getLab.Description,
		},
		CreatedAt: updatedPeminjaman.CreatedAt,
		UpdatedAt: updatedPeminjaman.UpdatedAt,
	}

	return peminjamanResponse, nil
}

func (u *peminjamanUsecase) DeletePeminjaman(id uint) error {
	return u.peminjamanRepo.DeletePeminjaman(id)
}

package usecases

import (
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/models"
	"sistem_peminjaman_be/repositories"
	"errors"
)

type LabUsecase interface {
	// admin
	GetAllLabs(page, limit int, name string) ([]dtos.LabResponse, int, error)
	GetLabByID(userId, id uint) (dtos.LabByIDResponse, error)
	CreateLab(lab *dtos.LabInput) (dtos.LabResponse, error)
	UpdateLab(id uint, labInput dtos.LabInput) (dtos.LabResponse, error)
	DeleteLab(id uint) error

	SearchLabAvailable(userId, page, limit int, name string) ([]dtos.LabResponse, int, error)
}

type labUsecase struct {
	labRepo               repositories.LabRepository
	labImageRepo          repositories.LabImageRepository
	historySearchRepo       repositories.HistorySearchRepository
	userRepo                repositories.UserRepository
	historySeenLabUsecase   HistorySeenLabUsecase
}

func NewLabUsecase(labRepo repositories.LabRepository, labImageRepo repositories.LabImageRepository, historySearchRepo repositories.HistorySearchRepository, userRepo repositories.UserRepository, historySeenLabUsecase HistorySeenLabUsecase) LabUsecase {
	return &labUsecase{labRepo, labImageRepo, historySearchRepo, userRepo, historySeenLabUsecase}
}


func (u *labUsecase) GetAllLabs(page, limit int, name string) ([]dtos.LabResponse, int, error) {
	labs, count, err := u.labRepo.SearchLabAvailable(page, limit, name)
	if err != nil {
		return nil, 0, err
	}

	var labResponses []dtos.LabResponse

	for _, lab := range labs {
		getImage, err := u.labImageRepo.GetAllLabImageByID(lab.ID)
		if err != nil {
			continue
		}

		var labImageResponses []dtos.LabImageResponse
		for _, image := range getImage {
			labImageResponse := dtos.LabImageResponse{
				LabID:  image.LabID,
				ImageUrl: image.ImageUrl,
			}
			labImageResponses = append(labImageResponses, labImageResponse)
		}

		labResponse := dtos.LabResponse{
			LabID:           lab.ID,
			Name:            lab.Name,
			LabImage:        labImageResponses,
			Description:     lab.Description,
			CreatedAt:       lab.CreatedAt,
			UpdatedAt:       lab.UpdatedAt,
		}

		labResponses = append(labResponses, labResponse)

	}

	return labResponses, count, nil
}

func (u *labUsecase) GetLabByID(userId, id uint) (dtos.LabByIDResponse, error) {
	var labResponses dtos.LabByIDResponse
	lab, err := u.labRepo.GetLabByID(id)
	if err != nil {
		return labResponses, err
	}

	historySeenLabInput := dtos.HistorySeenLabInput{
		LabID: lab.ID,
	}

	_, err = u.historySeenLabUsecase.CreateHistorySeenLab(userId, historySeenLabInput)
	if err != nil {
		return labResponses, err
	}

	getImage, err := u.labImageRepo.GetAllLabImageByID(lab.ID)
	if err != nil {
		return labResponses, err
	}


	var labImageResponses []dtos.LabImageResponse
	for _, image := range getImage {
		labImageResponse := dtos.LabImageResponse{
			LabID:    image.LabID,
			ImageUrl: image.ImageUrl,
		}
		labImageResponses = append(labImageResponses, labImageResponse)
	}

	labResponse := dtos.LabByIDResponse{
		LabID:           lab.ID,
		Name:            lab.Name,
		LabImage:        labImageResponses,
		Description:     lab.Description,
		CreatedAt:       lab.CreatedAt,
		UpdatedAt:       lab.UpdatedAt,
	}
	return labResponse, nil
}


func (u *labUsecase) CreateLab(lab *dtos.LabInput) (dtos.LabResponse, error) {
	var labResponse dtos.LabResponse
	if lab.Name == "" || lab.LabImage == nil || lab.Description == "" {
		return labResponse, errors.New("failed to create lab")
	}

	createLab := models.Lab{
		Name:        lab.Name,
		Description: lab.Description,
	}

	createdLab, err := u.labRepo.CreateLab(createLab)
	if err != nil {
		return labResponse, err
	}

	for _, labImage := range lab.LabImage {
		if labImage.ImageUrl == "" {
			return labResponse, errors.New("failed to create lab")
		}
		labImagee := models.LabImage{
			LabID:  createdLab.ID,
			ImageUrl: labImage.ImageUrl,
		}
		_, err = u.labImageRepo.CreateLabImage(labImagee)
		if err != nil {
			return labResponse, err
		}
	}

	getImage, err := u.labImageRepo.GetAllLabImageByID(createdLab.ID)
	if err != nil {
		return labResponse, err
	}

	var labImageResponses []dtos.LabImageResponse
	for _, image := range getImage {
		labImageResponse := dtos.LabImageResponse{
			LabID:    image.LabID,
			ImageUrl: image.ImageUrl,
		}
		labImageResponses = append(labImageResponses, labImageResponse)
	}

	labResponse = dtos.LabResponse{
		LabID:           createdLab.ID,
		Name:            createdLab.Name,
		LabImage:        labImageResponses,
		Description:     createdLab.Description,
		CreatedAt:       createdLab.CreatedAt,
		UpdatedAt:       createdLab.UpdatedAt,
	}
	return labResponse, nil
}


func (u *labUsecase) UpdateLab(id uint, lab dtos.LabInput) (dtos.LabResponse, error) {
	var labs models.Lab
	var labResponse dtos.LabResponse

	if lab.Name == "" || lab.LabImage == nil || lab.Description == "" {
		return labResponse, errors.New("failed to update lab")
	}

	labs, err := u.labRepo.GetLabByID(id)
	if err != nil {
		return labResponse, err
	}

	labs.Name = lab.Name
	labs.Description = lab.Description

	updatedLab, err := u.labRepo.UpdateLab(labs)
	if err != nil {
		return labResponse, err
	}

	u.labImageRepo.DeleteLabImage(id)

	for _, labImage := range lab.LabImage {
		if labImage.ImageUrl == "" {
			return labResponse, errors.New("failed to update lab")
		}
		labImagee := models.LabImage{
			LabID:  updatedLab.ID,
			ImageUrl: labImage.ImageUrl,
		}
		_, err = u.labImageRepo.UpdateLabImage(labImagee)
		if err != nil {
			return labResponse, err
		}
	}

	getImage, err := u.labImageRepo.GetAllLabImageByID(updatedLab.ID)
	if err != nil {
		return labResponse, err
	}

	var labImageResponses []dtos.LabImageResponse
	for _, image := range getImage {
		labImageResponse := dtos.LabImageResponse{
			LabID:    image.LabID,
			ImageUrl: image.ImageUrl,
		}
		labImageResponses = append(labImageResponses, labImageResponse)
	}

	labResponse = dtos.LabResponse{
		LabID:           updatedLab.ID,
		Name:            updatedLab.Name,
		LabImage:        labImageResponses,
		Description:     updatedLab.Description,
		CreatedAt:       updatedLab.CreatedAt,
		UpdatedAt:       updatedLab.UpdatedAt,
	}
	return labResponse, nil
}


func (u *labUsecase) DeleteLab(id uint) error {
	return u.labRepo.DeleteLab(id)
}


func (u *labUsecase) SearchLabAvailable(userId, page, limit int, name string) ([]dtos.LabResponse, int, error) {
	labs, count, err := u.labRepo.SearchLabAvailable(page, limit, name)
	if err != nil {
		return nil, 0, err
	}
	if len(labs) > 0 && name != "" && userId > 1 {
		historySearches := models.HistorySearch{
			UserID: uint(userId),
			Name:   name,
		}
		_, err := u.historySearchRepo.HistorySearchCreate(historySearches)
		if err != nil {
			return nil, 0, err
		}
	}

	var labResponses []dtos.LabResponse

	for _, lab := range labs {
		getImage, err := u.labImageRepo.GetAllLabImageByID(lab.ID)
		if err != nil {
			continue
		}

		var labImageResponses []dtos.LabImageResponse
		for _, image := range getImage {
			labImageResponse := dtos.LabImageResponse{
				LabID:    image.LabID,
				ImageUrl: image.ImageUrl,
			}
			labImageResponses = append(labImageResponses, labImageResponse)
		}

		labResponse := dtos.LabResponse{
			LabID:           lab.ID,
			Name:            lab.Name,
			LabImage:        labImageResponses,
			Description:     lab.Description,
			CreatedAt:       lab.CreatedAt,
			UpdatedAt:       lab.UpdatedAt,
		}

		labResponses = append(labResponses, labResponse)


	}

	return labResponses, count, nil

}
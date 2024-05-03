package usecases

import (
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/models"
	"sistem_peminjaman_be/repositories"
	"errors"
)

type HistorySeenLabUsecase interface {
	GetAllHistorySeenLabs(page, limit int, userId uint) ([]dtos.HistorySeenLabResponse, int, error)
	CreateHistorySeenLab(userId uint, historySeenLabInput dtos.HistorySeenLabInput) (dtos.HistorySeenLabResponse, error)
}

type historySeenLabUsecase struct {
	historySeenLabRepo repositories.HistorySeenLabRepository
	labRepo            repositories.LabRepository
	labImageRepo       repositories.LabImageRepository
}

func NewHistorySeenLabUsecase(historySeenLabRepo repositories.HistorySeenLabRepository, labRepo repositories.LabRepository, labImageRepo repositories.LabImageRepository) HistorySeenLabUsecase {
	return &historySeenLabUsecase{historySeenLabRepo, labRepo, labImageRepo}
}


func (u *historySeenLabUsecase) GetAllHistorySeenLabs(page, limit int, userId uint) ([]dtos.HistorySeenLabResponse, int, error) {
	historySeenLabs, count, err := u.historySeenLabRepo.GetAllHistorySeenLab(page, limit, userId)
	if err != nil {
		return nil, 0, err
	}

	var historySeenLabResponses []dtos.HistorySeenLabResponse
	for _, historySeenLab := range historySeenLabs {
		getLab, err := u.labRepo.GetLabByID(historySeenLab.LabID)
		if err != nil {
			continue
		}
		getLabImage, err := u.labImageRepo.GetAllLabImageByID(getLab.ID)
		if err != nil {
			continue
		}
		var labImageResponses []dtos.LabImageResponse
		for _, lab := range getLabImage {
			labImageResponse := dtos.LabImageResponse{
				LabID:  lab.LabID,
				ImageUrl: lab.ImageUrl,
			}
			labImageResponses = append(labImageResponses, labImageResponse)
		}
		
		historySeenLabResponse := dtos.HistorySeenLabResponse{
			ID: historySeenLab.ID,
			Lab: dtos.LabByIDSimply{
				LabID:           getLab.ID,
				Name:            getLab.Name,
				LabImage:        labImageResponses,
				Description:     getLab.Description,
			},
			CreatedAt: historySeenLab.CreatedAt,
			UpdatedAt: historySeenLab.UpdatedAt,
		}
		historySeenLabResponses = append(historySeenLabResponses, historySeenLabResponse)
	}

	return historySeenLabResponses, count, nil
}

func (u *historySeenLabUsecase) CreateHistorySeenLab(userId uint, historySeenLabInput dtos.HistorySeenLabInput) (dtos.HistorySeenLabResponse, error) {
	var historySeenLabResponses dtos.HistorySeenLabResponse

	if historySeenLabInput.LabID < 1 {
		return historySeenLabResponses, errors.New("Failed to create history seen lab")
	}

	createHistorySeenLab := models.HistorySeenLab{
		UserID:  userId,
		LabID: historySeenLabInput.LabID,
	}

	getHistorySeenLab, _ := u.historySeenLabRepo.GetHistorySeenLabByID(historySeenLabInput.LabID, userId)
	if getHistorySeenLab.ID > 0 {
		createHistorySeenLab, _ = u.historySeenLabRepo.UpdateHistorySeenLab(getHistorySeenLab)
	} else {
		createHistorySeenLab, _ = u.historySeenLabRepo.CreateHistorySeenLab(createHistorySeenLab)
	}

	getLab, _ := u.labRepo.GetLabByID(createHistorySeenLab.LabID)
	getLabImage, _ := u.labImageRepo.GetAllLabImageByID(getLab.ID)
	var labImageResponses []dtos.LabImageResponse
	for _, lab := range getLabImage {
		labImageResponse := dtos.LabImageResponse{
			LabID:  lab.LabID,
			ImageUrl: lab.ImageUrl,
		}
		labImageResponses = append(labImageResponses, labImageResponse)
	}
	
	historySeenLabResponse := dtos.HistorySeenLabResponse{
		ID: createHistorySeenLab.ID,
		Lab: dtos.LabByIDSimply{
			LabID:         getLab.ID,
			Name:            getLab.Name,
			LabImage:      labImageResponses,
			Description:     getLab.Description,
		},
		CreatedAt: createHistorySeenLab.CreatedAt,
		UpdatedAt: createHistorySeenLab.UpdatedAt,
	}
	return historySeenLabResponse, nil
}

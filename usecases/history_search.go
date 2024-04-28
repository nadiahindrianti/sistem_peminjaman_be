package usecases

import (
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/models"
	"sistem_peminjaman_be/repositories"
)

type HistorySearchUseCase interface {
	HistorySearchGetAll(userId uint, page, limit int) ([]dtos.HistorySearchResponse, int, error)
	HistorySearchCreate(userId uint, input dtos.HistorySearchInput) (dtos.HistorySearchResponse, error)
	HistorySearchDelete(userId, id uint) error
}

type historySearchUsecase struct {
	historySearchRepository repositories.HistorySearchRepository
	userRepo                repositories.UserRepository
}

func NewHistorySearchUsecase(historySearchRepository repositories.HistorySearchRepository, userRepo repositories.UserRepository) HistorySearchUseCase {
	return &historySearchUsecase{historySearchRepository, userRepo}
}

// HistorySearchGetAll godoc
// @Summary      Get all history search
// @Description  Get all history search
// @Tags         User - History Search
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success      200 {object} dtos.HistorySearchStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/history-search [get]
// @Security BearerAuth
func (u *historySearchUsecase) HistorySearchGetAll(userId uint, page, limit int) ([]dtos.HistorySearchResponse, int, error) {
	var historySearchResponses []dtos.HistorySearchResponse

	histories, count, err := u.historySearchRepository.HistorySearchGetByUserId(userId, page, limit)
	if err != nil {
		return historySearchResponses, count, err
	}

	for _, history := range histories {
		historySearchResponse := dtos.HistorySearchResponse{
			ID:     history.ID,
			UserID: history.UserID,
			Name:   history.Name,
		}
		historySearchResponses = append(historySearchResponses, historySearchResponse)
	}

	return historySearchResponses, count, nil
}

// HistorySearchCreate godoc
// @Summary      Create a new history search
// @Description  Create a new history search
// @Tags         User - History Search
// @Accept       json
// @Produce      json
// @Param        request body dtos.HistorySearchInput true "Payload Body [RAW]"
// @Success      201 {object} dtos.HistorySearchCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/history-search [post]
// @Security BearerAuth
func (u *historySearchUsecase) HistorySearchCreate(userId uint, input dtos.HistorySearchInput) (dtos.HistorySearchResponse, error) {
	var (
		history         models.HistorySearch
		historyResponse dtos.HistorySearchResponse
	)

	history.UserID = userId
	history.Name = input.Name

	history, err := u.historySearchRepository.HistorySearchCreate(history)
	if err != nil {
		return historyResponse, err
	}

	historyResponse.ID = history.ID
	historyResponse.UserID = history.UserID
	historyResponse.Name = history.Name

	return historyResponse, nil
}

// HistorySearchDelete godoc
// @Summary      Delete a history search
// @Description  Delete a history search
// @Tags         User - History Search
// @Accept       json
// @Produce      json
// @Param id path integer true "ID history search"
// @Success      200 {object} dtos.StatusOKDeletedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/history-search/{id} [delete]
// @Security BearerAuth
func (u *historySearchUsecase) HistorySearchDelete(userId, id uint) error {
	err := u.historySearchRepository.HistorySearchDelete(userId, id)
	if err != nil {
		return err
	}

	return nil
}

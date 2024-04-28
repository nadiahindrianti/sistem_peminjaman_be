package usecases

import (
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/repositories"
	"strings"
)

type NotificationUsecase interface {
	GetNotificationByUserID(id uint) (dtos.NotificationResponse, error)
}

type notificationUsecase struct {
	notificationRepo    repositories.NotificationRepository
	templateMessageRepo repositories.TemplateMessageRepository
	userRepo            repositories.UserRepository
}

func NewNotificationUsecase(notificationRepo repositories.NotificationRepository, templateMessageRepo repositories.TemplateMessageRepository, userRepo repositories.UserRepository) NotificationUsecase {
	return &notificationUsecase{notificationRepo, templateMessageRepo, userRepo}
}

// GetNotificationByUserID godoc
// @Summary      Get notification by user id
// @Description  Get notification by user id
// @Tags         User - Notification
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.GetNotificationByUserIDStatusOKResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/notification [get]
// @Security BearerAuth
func (u *notificationUsecase) GetNotificationByUserID(id uint) (dtos.NotificationResponse, error) {
	var notificationResponsee dtos.NotificationResponse
	notifications, err := u.notificationRepo.GetNotificationByUserID(id)
	if err != nil {
		return notificationResponsee, err
	}

	getUser, err := u.userRepo.UserGetById(id)
	if err != nil {
		return notificationResponsee, err
	}

	var templateContentResponses []dtos.TemplateMessageByUserIDResponse

	for _, notification := range notifications {
		getTemplate, err := u.templateMessageRepo.GetTemplateMessageByID(notification.TemplateID)
		if err != nil {
			return notificationResponsee, err
		}

		newTitle := strings.Replace(getTemplate.Title, "[Nama Pengguna]", getUser.FullName, -1)
		newContent := strings.Replace(getTemplate.Content, "[Nama Pengguna]", getUser.FullName, -1)

		templateContentResponse := dtos.TemplateMessageByUserIDResponse{
			Title:     newTitle,
			Content:   newContent,
			CreatedAt: notification.CreatedAt,
			UpdatedAt: notification.UpdatedAt,
		}

		templateContentResponses = append(templateContentResponses, templateContentResponse)
	}

	notificationResponse := dtos.NotificationResponse{
		UserID:              id,
		NotificationContent: templateContentResponses,
	}

	return notificationResponse, nil
}
package usecases

import (
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/models"
	"sistem_peminjaman_be/repositories"
	"errors"
)

type TemplateMessageUsecase interface {
	GetAllTemplateMessages(page, limit int) ([]dtos.TemplateMessageResponse, int, error)
	GetTemplateMessageByID(id uint) (dtos.TemplateMessageResponse, error)
	CreateTemplateMessage(template *dtos.TemplateMessageInput) (dtos.TemplateMessageResponse, error)
	UpdateTemplateMessage(id uint, template dtos.TemplateMessageInput) (dtos.TemplateMessageResponse, error)
	DeleteTemplateMessage(id uint) error
}

type templateMessageUsecase struct {
	templateMessageRepo repositories.TemplateMessageRepository
}

func NewTemplateMessageUsecase(TemplateMessageRepo repositories.TemplateMessageRepository) TemplateMessageUsecase {
	return &templateMessageUsecase{TemplateMessageRepo}
}

// GetAllArticles godoc
// @Summary      Get all articles
// @Description  Get all articles
// @Tags         Admin - Article
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success      200 {object} dtos.GetAllArticleStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/article [get]
func (u *templateMessageUsecase) GetAllTemplateMessages(page, limit int) ([]dtos.TemplateMessageResponse, int, error) {
	templates, count, err := u.templateMessageRepo.GetAllTemplateMessages(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var templateMessageResponses []dtos.TemplateMessageResponse
	for _, template := range templates {
		templateResponse := dtos.TemplateMessageResponse{
			TemplateID: template.ID,
			Title:             template.Title,
			Content:           template.Content,
			CreatedAt:         template.CreatedAt,
			UpdatedAt:         template.UpdatedAt,
		}
		templateMessageResponses = append(templateMessageResponses, templateResponse)
	}

	return templateMessageResponses, count, nil
}

// GetArticleByID godoc
// @Summary      Get article by ID
// @Description  Get article by ID
// @Tags         Admin - Article
// @Accept       json
// @Produce      json
// @Param id path integer true "ID article"
// @Success      200 {object} dtos.ArticleStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/article/{id} [get]
func (u *templateMessageUsecase) GetTemplateMessageByID(id uint) (dtos.TemplateMessageResponse, error) {
	var templateResponses dtos.TemplateMessageResponse
	template, err := u.templateMessageRepo.GetTemplateMessageByID(id)
	if err != nil {
		return templateResponses, err
	}
	templateResponse := dtos.TemplateMessageResponse{
		TemplateID: template.ID,
		Title:             template.Title,
		Content:           template.Content,
		CreatedAt:         template.CreatedAt,
		UpdatedAt:         template.UpdatedAt,
	}
	return templateResponse, nil
}

// CreateArticle godoc
// @Summary      Create a new article
// @Description  Create a new article
// @Tags         Admin - Article
// @Accept       json
// @Produce      json
// @Param        file formData file true "Image file"
// @Param		 title formData string true "Title article"
// @Param		 description formData string true "Description article"
// @Param		 label formData string true "Label article"
// @Success      201 {object} dtos.ArticleCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/article [post]
// @Security BearerAuth
func (u *templateMessageUsecase) CreateTemplateMessage(templateInput *dtos.TemplateMessageInput) (dtos.TemplateMessageResponse, error) {
	var templateResponses dtos.TemplateMessageResponse

	if templateInput.Title == "" || templateInput.Content == "" {
		return templateResponses, errors.New("Failed to create template message")
	}

	createTemplate := models.TemplateMessage{
		Title:   templateInput.Title,
		Content: templateInput.Content,
	}

	createdTemplate, err := u.templateMessageRepo.CreateTemplateMessage(createTemplate)
	if err != nil {
		return templateResponses, err
	}

	templateResponse := dtos.TemplateMessageResponse{
		TemplateID: createdTemplate.ID,
		Title:             createdTemplate.Title,
		Content:           createdTemplate.Content,
		CreatedAt:         createdTemplate.CreatedAt,
		UpdatedAt:         createdTemplate.UpdatedAt,
	}
	return templateResponse, nil
}

// UpdateArticle godoc
// @Summary      Update article
// @Description  Update article
// @Tags         Admin - Article
// @Accept       json
// @Produce      json
// @Param		 id path integer true "ID article"
// @Param        file formData file true "Image file"
// @Param		 title formData string true "Title article"
// @Param		 description formData string true "Description article"
// @Param		 label formData string true "Label article"
// @Success      200 {object} dtos.ArticleStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/article/{id} [put]
// @Security BearerAuth
func (u *templateMessageUsecase) UpdateTemplateMessage(id uint, templateInput dtos.TemplateMessageInput) (dtos.TemplateMessageResponse, error) {
	var template models.TemplateMessage
	var templateResponse dtos.TemplateMessageResponse

	if templateInput.Title == "" || templateInput.Content == "" {
		return templateResponse, errors.New("Failed to update template message")
	}

	template, err := u.templateMessageRepo.GetTemplateMessageByID(id)
	if err != nil {
		return templateResponse, err
	}

	template.Title = templateInput.Title
	template.Content = templateInput.Content

	template, err = u.templateMessageRepo.UpdateTemplate(template)

	if err != nil {
		return templateResponse, err
	}

	templateResponse.TemplateID = template.ID
	templateResponse.Title = template.Title
	templateResponse.Content = template.Content
	templateResponse.CreatedAt = template.CreatedAt
	templateResponse.UpdatedAt = template.UpdatedAt

	return templateResponse, nil

}

// DeleteArticle godoc
// @Summary      Delete a article
// @Description  Delete a article
// @Tags         Admin - Article
// @Accept       json
// @Produce      json
// @Param id path integer true "ID article"
// @Success      200 {object} dtos.StatusOKDeletedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/article/{id} [delete]
// @Security BearerAuth
func (u *templateMessageUsecase) DeleteTemplateMessage(id uint) error {
	return u.templateMessageRepo.DeleteTemplate(id)
}

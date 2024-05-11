package controllers

import (
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/helpers"
	"sistem_peminjaman_be/usecases"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TemplateMessageController interface {
	GetAllTemplateMessages(c echo.Context) error
	GetTemplateMessageByID(c echo.Context) error
	CreateTemplateMessage(c echo.Context) error
	UpdateTemplateMessage(c echo.Context) error
	DeleteTemplateMessage(c echo.Context) error
}

type templateMessageController struct {
	templateMessageUsecase usecases.TemplateMessageUsecase
}

func NewTemplateMessageController(templateMessageUsecase usecases.TemplateMessageUsecase) TemplateMessageController {
	return &templateMessageController{templateMessageUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *templateMessageController) GetAllTemplateMessages(ctx echo.Context) error {
	pageParam := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limitParam := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}

	templates, count, err := c.templateMessageUsecase.GetAllTemplateMessages(page, limit)
	if err != nil {

		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed fetching template messages",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all template message",
			templates,
			page,
			limit,
			count,
		),
	)
}

func (c *templateMessageController) GetTemplateMessageByID(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	template, err := c.templateMessageUsecase.GetTemplateMessageByID(uint(id))

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get template message by id",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get template message by id",
			template,
		),
	)

}

func (c *templateMessageController) CreateTemplateMessage(ctx echo.Context) error {
	var templateInput dtos.TemplateMessageInput
	if err := ctx.Bind(&templateInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding template message",
				helpers.GetErrorData(err),
			),
		)
	}

	template, err := c.templateMessageUsecase.CreateTemplateMessage(&templateInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a template message",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a template message",
			template,
		),
	)
}

func (c *templateMessageController) UpdateTemplateMessage(ctx echo.Context) error {

	var templateInput dtos.TemplateMessageInput
	if err := ctx.Bind(&templateInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding template message",
				helpers.GetErrorData(err),
			),
		)
	}

	id, _ := strconv.Atoi(ctx.Param("id"))

	template, err := c.templateMessageUsecase.GetTemplateMessageByID(uint(id))
	if template.TemplateID == 0 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get template message by id",
				helpers.GetErrorData(err),
			),
		)
	}

	templateResp, err := c.templateMessageUsecase.UpdateTemplateMessage(uint(id), templateInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding template message",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully updated template message",
			templateResp,
		),
	)
}

func (c *templateMessageController) DeleteTemplateMessage(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := c.templateMessageUsecase.DeleteTemplateMessage(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to delete template message",
				helpers.GetErrorData(err),
			),
		)
	}
	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully deleted template message",
			nil,
		),
	)
}

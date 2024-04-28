package controllers

import (
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/helpers"
	"sistem_peminjaman_be/middlewares"
	"sistem_peminjaman_be/usecases"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type HistorySeenLabController interface {
	GetAllHistorySeenLabs(c echo.Context) error
	CreateHistorySeenLab(c echo.Context) error
}

type historySeenLabController struct {
	historySeenLabUsecase usecases.HistorySeenLabUsecase
}

func NewHistorySeenLabController(historySeenLabUsecase usecases.HistorySeenLabUsecase) HistorySeenLabController {
	return &historySeenLabController{historySeenLabUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *historySeenLabController) GetAllHistorySeenLabs(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				"Unauthorized",
			),
		)
	}

	userId, err := middlewares.GetUserIdFromToken(tokenString)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				helpers.GetErrorData(err),
			),
		)
	}

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

	historySeenLabs, count, err := c.historySeenLabUsecase.GetAllHistorySeenLabs(page, limit, userId)
	if err != nil {

		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed fetching historySeenLabs",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all historySeenLab",
			historySeenLabs,
			page,
			limit,
			count,
		),
	)
}

func (c *historySeenLabController) CreateHistorySeenLab(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				"Unauthorized",
			),
		)
	}

	userId, err := middlewares.GetUserIdFromToken(tokenString)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				helpers.GetErrorData(err),
			),
		)
	}
	var historySeenLabInput dtos.HistorySeenLabInput
	if err := ctx.Bind(&historySeenLabInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding historySeenLab",
				helpers.GetErrorData(err),
			),
		)
	}

	historySeenLab, err := c.historySeenLabUsecase.CreateHistorySeenLab(userId, historySeenLabInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a historySeenLab",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a historySeenLab",
			historySeenLab,
		),
	)
}

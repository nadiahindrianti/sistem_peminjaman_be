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

type JadwalController interface {
	GetAllJadwals(c echo.Context) error
	GetJadwalByID(c echo.Context) error
	CreateJadwal(c echo.Context) error
	UpdateJadwal(c echo.Context) error
	DeleteJadwal(c echo.Context) error
	SearchJadwalAvailable(c echo.Context) error
}

type jadwalController struct {
	jadwalUsecase usecases.JadwalUsecase
}

func NewJadwalController(jadwalUsecase usecases.JadwalUsecase) JadwalController {
	return &jadwalController{jadwalUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *jadwalController) GetAllJadwals(ctx echo.Context) error {
	pageParam := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limitParam := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 1000
	}

	name_laboratoriumParam := ctx.QueryParam("name_laboratorium")

	jadwals, count, err := c.jadwalUsecase.GetAllJadwals(page, limit, name_laboratoriumParam)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get all jadwal",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all jadwals",
			jadwals,
			page,
			limit,
			count,
		),
	)
}

func (c *jadwalController) GetJadwalByID(ctx echo.Context) error {
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

	id, _ := strconv.Atoi(ctx.Param("id"))
	jadwal, err := c.jadwalUsecase.GetJadwalByID(userId, uint(id))

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get jadwal by id",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get jadwal by id",
			jadwal,
		),
	)

}

func (c *jadwalController) CreateJadwal(ctx echo.Context) error {
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
	var jadwalDTO dtos.JadwalInput
	if err := ctx.Bind(&jadwalDTO); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding jadwal",
				helpers.GetErrorData(err),
			),
		)
	}

	jadwal, err := c.jadwalUsecase.CreateJadwal(&jadwalDTO)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a jadwal",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a jadwal",
			jadwal,
		),
	)
}

func (c *jadwalController) UpdateJadwal(ctx echo.Context) error {
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
	var jadwalInput dtos.JadwalInput
	if err := ctx.Bind(&jadwalInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding jadwal",
				helpers.GetErrorData(err),
			),
		)
	}

	id, _ := strconv.Atoi(ctx.Param("id"))

	jadwal, err := c.jadwalUsecase.GetJadwalByID(userId, uint(id))
	if jadwal.JadwalID == 0 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get jadwal by id",
				helpers.GetErrorData(err),
			),
		)
	}

	jadwalResp, err := c.jadwalUsecase.UpdateJadwal(uint(id), jadwalInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to updated a jadwal",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully updated jadwal",
			jadwalResp,
		),
	)
}

func (c *jadwalController) DeleteJadwal(ctx echo.Context) error {
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
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := c.jadwalUsecase.DeleteJadwal(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to delete jadwal",
				helpers.GetErrorData(err),
			),
		)
	}
	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully deleted jadwal",
			nil,
		),
	)
}



func (c *jadwalController) SearchJadwalAvailable(ctx echo.Context) error {
	userId := uint(1)
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		userId = 1
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
		limit = 1000
	}

	name_laboratoriumParam := ctx.QueryParam("name_laboratorium")

	jadwals, count, err := c.jadwalUsecase.SearchJadwalAvailable(int(userId), page, limit, name_laboratoriumParam)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get all jadwal",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all jadwals",
			jadwals,
			page,
			limit,
			count,
		),
	)
}

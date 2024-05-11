package controllers

import (
	"net/http"
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/helpers"
	"sistem_peminjaman_be/middlewares"
	"sistem_peminjaman_be/usecases"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PeminjamanController interface {
	GetAllPeminjamans(c echo.Context) error
	GetPeminjamanByID(c echo.Context) error
	GetPeminjamansByAdmin(c echo.Context) error
	CreatePeminjaman(c echo.Context) error
	UpdatePeminjaman(c echo.Context) error
	DeletePeminjaman(c echo.Context) error
}

type peminjamanController struct {
	peminjamanUsecase usecases.PeminjamanUsecase
}

func NewPeminjamanController(peminjamanUsecase usecases.PeminjamanUsecase) PeminjamanController {
	return &peminjamanController{peminjamanUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *peminjamanController) GetAllPeminjamans(ctx echo.Context) error {
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

	nameLaboratoriumParam := ctx.QueryParam("name")
	statusParam := ctx.QueryParam("status")

	peminjamans, count, err := c.peminjamanUsecase.GetPeminjamans(page, limit, userId, nameLaboratoriumParam, statusParam)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get all peminjaman",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all peminjamans",
			peminjamans,
			page,
			limit,
			count,
		),
	)
}

func (c *peminjamanController) GetPeminjamansByAdmin(ctx echo.Context) error {
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

	searchParam := ctx.QueryParam("search")
	statusParam := ctx.QueryParam("status")


	peminjaman, count, err := c.peminjamanUsecase.GetPeminjamansByAdmin(pageParam, limitParam, searchParam, statusParam)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get a peminjaman",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully to get peminjamans",
			peminjaman,
			page,
			limit,
			count,
		),
	)
}

func (c *peminjamanController) GetPeminjamanByID(ctx echo.Context) error {
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
	peminjaman, err := c.peminjamanUsecase.GetPeminjamanByID(userId, uint(id))

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get peminjaman by id",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get peminjaman by id",
			peminjaman,
		),
	)

}

func (c *peminjamanController) CreatePeminjaman(ctx echo.Context) error {
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

	var peminjamanDTO dtos.PeminjamanInput
	if err := ctx.Bind(&peminjamanDTO); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding peminjaman",
				helpers.GetErrorData(err),
			),
		)
	}

	peminjaman, err := c.peminjamanUsecase.CreatePeminjaman(userId, &peminjamanDTO)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a peminjaman",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a peminjaman",
			peminjaman,
		),
	)
}

func (c *peminjamanController) UpdatePeminjaman(ctx echo.Context) error {
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
	var peminjamanInput dtos.PeminjamanInput
	if err := ctx.Bind(&peminjamanInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding peminjaman",
				helpers.GetErrorData(err),
			),
		)
	}

	id, _ := strconv.Atoi(ctx.Param("id"))

	peminjaman, err := c.peminjamanUsecase.GetPeminjamanByID(userId, uint(id))
	if peminjaman.PeminjamanID == 0 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get peminjaman by id",
				helpers.GetErrorData(err),
			),
		)
	}

	peminjamanResp, err := c.peminjamanUsecase.UpdatePeminjaman(uint(id), userId, peminjamanInput)
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
			"Successfully updated peminjaman",
			peminjamanResp,
		),
	)
}

func (c *peminjamanController) DeletePeminjaman(ctx echo.Context) error {
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

	err := c.peminjamanUsecase.DeletePeminjaman(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to delete peminjaman",
				helpers.GetErrorData(err),
			),
		)
	}
	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully deleted peminjaman",
			nil,
		),
	)
}

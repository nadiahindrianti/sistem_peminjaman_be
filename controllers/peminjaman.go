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
	AdminGetPeminjamanByID(c echo.Context) error
	GetPeminjamansDetailByAdmin(c echo.Context) error
	CreatePeminjaman(c echo.Context) error
	AdminUpdatePeminjaman(c echo.Context) error
	UpdatePeminjaman(c echo.Context) error
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


	peminjaman, count, err := c.peminjamanUsecase.GetPeminjamansByAdmin(page, limit, searchParam, statusParam)
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

func (c *peminjamanController) GetPeminjamansDetailByAdmin(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	peminjaman, err := c.peminjamanUsecase.GetPeminjamansDetailByAdmin(uint(id))

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

func (c *peminjamanController) AdminGetPeminjamanByID(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	peminjaman, err := c.peminjamanUsecase.AdminGetPeminjamanByID( uint(id))

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

//error func

func (c *peminjamanController) UpdatePeminjaman(ctx echo.Context) error {

	// Binding input peminjaman
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

	// Ambil ID dari parameter URL
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Dapatkan detail peminjaman berdasarkan ID
	peminjaman, err := c.peminjamanUsecase.AdminGetPeminjamanByID(uint(id))
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

	// Perbarui peminjaman
	peminjamanResp, err := c.peminjamanUsecase.UpdatePeminjaman(uint(id), peminjamanInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to update peminjaman",
				helpers.GetErrorData(err),
			),
		)
	}

	// Berhasil memperbarui peminjaman
	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully updated peminjaman",
			peminjamanResp,
		),
	)
}


func (c *peminjamanController) AdminUpdatePeminjaman(ctx echo.Context) error {

	// Binding input peminjaman
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

	// Ambil ID dari parameter URL
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Dapatkan detail peminjaman berdasarkan ID
	peminjaman, err := c.peminjamanUsecase.AdminGetPeminjamanByID(uint(id))
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

	// Perbarui peminjaman
	peminjamanResp, err := c.peminjamanUsecase.AdminUpdatePeminjaman(uint(id), peminjamanInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to update peminjaman",
				helpers.GetErrorData(err),
			),
		)
	}

	// Berhasil memperbarui peminjaman
	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully updated peminjaman",
			peminjamanResp,
		),
	)
}


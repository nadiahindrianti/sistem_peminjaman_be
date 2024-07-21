package controllers

import (
	"net/http"
	"sistem_peminjaman_be/helpers"
	"sistem_peminjaman_be/middlewares"
	"sistem_peminjaman_be/usecases"
	"strconv"

	"github.com/labstack/echo/v4"
)

type DashboardController interface {
	DashboardGetAll(c echo.Context) error
	DashboardGetByMonth(c echo.Context) error
}

type dashboardController struct {
	dashboardUsecase usecases.DashboardUsecase
}

func NewDashboardController(dashboardUsecase usecases.DashboardUsecase) DashboardController {
	return &dashboardController{dashboardUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *dashboardController) DashboardGetAll(ctx echo.Context) error {
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
	dashboards, err := c.dashboardUsecase.DashboardGetAll()
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed fetching dashboard",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully get all dashboards",
			dashboards,
		),
	)
}

// dashboard laporan bulanan

func (c *dashboardController) DashboardGetByMonth(ctx echo.Context) error {
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

	monthParam := ctx.QueryParam("month")
	yearParam := ctx.QueryParam("year")
	month, err := strconv.Atoi(monthParam)
	if err != nil || month < 1 || month > 12 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Invalid month parameter",
				"Bad Request",
			),
		)
	}

	year, err := strconv.Atoi(yearParam)
	if err != nil || year < 1 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Invalid year parameter",
				"Bad Request",
			),
		)
	}

	dashboardResponse, err := c.dashboardUsecase.DashboardGetByMonth(month, year)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed fetching dashboard by month",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully get dashboards by month",
			dashboardResponse,
		),
	)
}

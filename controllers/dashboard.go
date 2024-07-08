package controllers

import (
	"sistem_peminjaman_be/helpers"
	"sistem_peminjaman_be/middlewares"
	"sistem_peminjaman_be/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DashboardController interface {
	DashboardGetAll(c echo.Context) error
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

package main

import (
	"sistem_peminjaman_be/configs"
	"sistem_peminjaman_be/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// create a new echo instance
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := configs.ConnectDB()
	if err != nil {
		panic(err)
	}

	err = configs.MigrateDB(db)
	if err != nil {
		panic(err)
	}

	routes.Init(e, db)

	e.Logger.Fatal(e.Start(":8000"))
}

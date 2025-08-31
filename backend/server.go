package main

import (
	"OJ-backend/config"
	model "OJ-backend/models"
	"OJ-backend/routes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.LoadEnv()

	e := echo.New()
	_, err := config.ConnectRabbitMQ()
	if err != nil {
		e.Logger.Fatal("Failed to connect to RabbitMQ:", err)
	}
	db, err := config.ConnectDB()
	if err != nil {
		e.Logger.Fatal("Failed to connect to the database:", err)
	} else {
		e.Logger.Info("Successfully connected to the database", db.Name())
	}
	defer func() {
		config.CloseDB()
		config.CloseRabbitMQ()
	}()
	
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Connect to the database
	db.AutoMigrate(model.User{}, model.Contest{}, model.Problem{}, model.Submission{}, model.TestCase{}, model.Language{})

	// Register routes
	routes.RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}

package main

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
    "OJ-backend/config"
	"OJ-backend/routes"
	"OJ-backend/models"
	rabbitmq "OJ-backend/services/rabbitmq"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	//Close rabbitmq connection when the application exits
	defer rabbitmq.CloseRabbitMQ()

	// Initialize Echo instance
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Connect to the database
	db, err := config.ConnectDB()
	if err != nil {
		e.Logger.Fatal("Failed to connect to the database:", err)
	} else {
		e.Logger.Info("Successfully connected to the database", db.Name())
	}
	db.AutoMigrate(model.User{}, model.Contest{}, model.Problem{}, model.Submission{},model.TestCase{})

	// Register routes
	routes.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
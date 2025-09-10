package main

import (
	"OJ-backend/config"
	model "OJ-backend/models"
	"OJ-backend/routes"
	"OJ-backend/seed"
	"flag"
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

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Backend is healthy!")
	})

	// Connect to the database
	db.AutoMigrate(model.User{}, model.Contest{}, model.Problem{}, model.Submission{}, model.TestCase{}, model.Language{})

	seedDB := flag.Bool("seed", false, "Seed the database with initial data")
	flag.Parse()

	if *seedDB {
		if err := seed.SeedDB(); err != nil {
			e.Logger.Fatal("Failed to seed the database:", err)
		} else {
			e.Logger.Info("Database seeded successfully")
		}
	}
	// Register routes
	routes.RegisterRoutes(e)
	port := config.GetEnv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}

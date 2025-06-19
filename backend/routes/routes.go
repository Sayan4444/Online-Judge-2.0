package routes

import (
	"github.com/labstack/echo/v4"	
	"OJ-backend/controllers"
)

func RegisterRoutes(e *echo.Echo) {
	// Public routes
	e.POST("/login", handler.Login)
	e.POST("/admin/login", handler.AdminLogin)

	// Protected routes
	api := e.Group("/api")
	api.Use(handler.JWTMiddleware())
	api.GET("/profile", handler.GetProfile)
	api.PUT("/profile", handler.UpdateProfile)

	// Admin routes	
	admin := e.Group("/admin")
	admin.Use(handler.AdminJWTMiddleware())
	admin.POST("/create-contest", handler.CreateContest)
	admin.GET("/contests", handler.GetAllContests)
	admin.PUT("/contest/:id", handler.UpdateContest)
	admin.DELETE("/contest/:id", handler.DeleteContest)
}
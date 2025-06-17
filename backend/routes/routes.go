package routes

import (
	"github.com/labstack/echo/v4"	
	"OJ-backend/controllers"
)

func RegisterRoutes(e *echo.Echo) {
	// Public routes
	e.POST("/login", handler.Login)

	// Protected routes
	protected := e.Group("/api")
	protected.Use(handler.JWTMiddleware())
	protected.GET("/profile", handler.GetProfile)
	// protected.PUT("/profile", profile.UpdateProfile)
	// protected.DELETE("/profile", profile.DeleteProfile)
}
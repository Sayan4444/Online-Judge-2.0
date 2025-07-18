package routes

import (
	"github.com/labstack/echo/v4"	
	"OJ-backend/controllers"
)

func RegisterRoutes(e *echo.Echo) {
	// Public routes
	e.POST("/login", handler.Login)
	e.POST("/admin/login", handler.AdminLogin)
	e.GET("/contests",handler.GetAllContests)
	
	// Protected routes
	api := e.Group("/api")
	api.Use(handler.JWTMiddleware())
	api.GET("/profile", handler.GetProfile)
	api.PUT("/profile", handler.UpdateProfile)
	api.GET("/problems/:id", handler.GetAllProblemsByContestID)
	api.GET("/problem/:id", handler.GetProblemByID)
	api.GET("/testcases/:id", handler.GetAllTestCasesByProblemID)
	api.POST("/submit/:user_id/:problem_id", handler.HandleSubmission)
	api.GET("/leaderboard/:contest_id", handler.GetLeaderboardByContestID)

	// Admin routes	
	admin := e.Group("/admin")
	admin.Use(handler.AdminJWTMiddleware())
	//contest routes
	admin.POST("/create-contest", handler.CreateContest)
	admin.PUT("/contest/:id", handler.UpdateContest)
	admin.DELETE("/contest/:id", handler.DeleteContest)
	//problem routes
	admin.POST("/create-problem/:id", handler.CreateProblem)
	admin.GET("/problems/:id", handler.GetAllProblemsByContestID)
	admin.PUT("/problem/:id", handler.UpdateProblem)
	admin.DELETE("/problem/:id", handler.DeleteProblem)
	//test case routes
	admin.POST("/create-testcase/:id", handler.CreateTestCase)
	admin.GET("/testcases/:id", handler.GetAllTestCasesByProblemID)
	admin.PUT("/testcase/:id", handler.UpdateTestCase)
	admin.DELETE("/testcase/:id", handler.DeleteTestCase)
}
package handler

import (
	"OJ-backend/config"
	models "OJ-backend/models"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"OJ-backend/services/rabbitmq"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var jwtSecret = []byte(config.GetEnv("JWT_SECRET"))
var adminSecret = []byte(config.GetEnv("ADMIN_SECRET"))

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// JWT middleware for user authentication
func JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: jwtSecret,
		ContextKey: "user",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(Claims)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			fmt.Println("JWT error:", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
		},
	})
}

func GetUserFromContext(c echo.Context) (username, email string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*Claims)
	return claims.Username, claims.Email
}

func Login(c echo.Context) error {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := c.Bind(&body); err != nil || body.Email == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}

	claims := &Claims{
		Username: body.Username,
		Email:    body.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	// Check if the user exists in the database
	var user models.User
	db := config.DB
	if err := db.Where("email = ?", body.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User does not exist, create a new user
			user = models.User{
				ID:       uuid.New(),
				Username: body.Username,
				Email:    body.Email,
			}
			if err := db.Create(&user).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create user"})
			}
		} else {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
		}
	} else {
		// User exists, update the username if it has changed
		if user.Username != body.Username {
			user.Username = body.Username
			if err := db.Save(&user).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update user"})
			}
		}
		claims.Username = user.Username
		claims.Email = user.Email
		claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Now())
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(72 * time.Hour))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// retrieve the user's profile information
func GetProfile(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*Claims)
	db := config.DB
	var userModel models.User

	result := db.Where("email = ?", claims.Email).First(&userModel)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}

	profile := map[string]interface{}{
		"username":   userModel.Username,
		"email":      userModel.Email,
		"created_at": userModel.CreatedAt,
	}
	return c.JSON(http.StatusOK, profile)
}

// Update the user's profile information
func UpdateProfile(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*Claims)
	db := config.DB
	var userModel models.User

	if err := db.Where("email = ?", claims.Email).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}

	var body struct {
		Username string `json:"username"`
	}

	if err := c.Bind(&body); err != nil || body.Username == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}

	userModel.Username = body.Username
	if err := db.Save(&userModel).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update profile"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "profile updated successfully"})
}

// Admin login and JWT middleware
type AdminClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func AdminLogin(c echo.Context) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	adminEmail := config.GetEnv("ADMIN_EMAIL")
	adminPassword := config.GetEnv("ADMIN_PASSWORD")

	if body.Email != adminEmail || body.Password != adminPassword {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}

	claims := &AdminClaims{
		Email: body.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(adminSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": signedToken})
}

func AdminJWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: adminSecret,
		ContextKey: "admin",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(AdminClaims)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			fmt.Println("Admin JWT error:", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired admin token")
		},
	})
}

// Get all contests
func GetAllContests(c echo.Context) error {
	db := config.DB
	var contests []models.Contest

	if err := db.Find(&contests).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to retrieve contests"})
	}

	return c.JSON(http.StatusOK, contests)
}

// Create Contest
func CreateContest(c echo.Context) error {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
	}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	startTime, err := time.Parse(time.RFC3339, body.StartTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid start_time format"})
	}

	endTime, err := time.Parse(time.RFC3339, body.EndTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid end_time format"})
	}

	contest := models.Contest{
		ID:          uuid.New(),
		Name:        body.Name,
		Description: body.Description,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	db := config.DB

	if err := db.Create(&contest).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not create contest"})
	}

	return c.JSON(http.StatusCreated, contest)
}

func UpdateContest(c echo.Context) error {
	contestID := c.Param("id")
	db := config.DB
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
	}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	startTime, err := time.Parse(time.RFC3339, body.StartTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid start_time format"})
	}

	endTime, err := time.Parse(time.RFC3339, body.EndTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid end_time format"})
	}

	var contest models.Contest

	if err := db.First(&contest, "id = ?", contestID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "contest not found"})
	}

	contest.Name = body.Name
	contest.Description = body.Description
	contest.StartTime = startTime
	contest.EndTime = endTime

	if err := db.Save(&contest).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not update contest"})
	}

	return c.JSON(http.StatusOK, contest)
}

// Delete Contest
func DeleteContest(c echo.Context) error {
	contestID := c.Param("id")
	db := config.DB

	// Check if the contest exists
	var contest models.Contest
	if err := db.First(&contest, "id = ?", contestID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "contest not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}

	// Delete the contest
	if err := db.Delete(&contest).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not delete contest"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "contest deleted successfully"})
}

// Get Problems by Contest ID
func GetAllProblemsByContestID(c echo.Context) error {
	contestID := c.Param("id")
	db := config.DB
	var problems []models.Problem

	if err := db.Where("contest_id = ?", contestID).Find(&problems).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to retrieve problems"})
	}

	if len(problems) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "no problems found for this contest"})
	}

	return c.JSON(http.StatusOK, problems)
}

// Create Problems in a Contest
func CreateProblem(c echo.Context) error {
	contestID := c.Param("id")
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	db := config.DB
	var contest models.Contest

	if err := db.First(&contest, "id = ?", contestID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "contest not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}

	problem := models.Problem{
		ID:          uuid.New(),
		ContestID:   contest.ID,
		Title:       body.Title,
		Description: body.Description,
	}

	if err := db.Create(&problem).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not create problem"})
	}

	return c.JSON(http.StatusCreated, problem)
}

// Update Problem in a Contest
func UpdateProblem(c echo.Context) error {
	problemID := c.Param("id")
	db := config.DB
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := c.Bind(&body); err != nil {		
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}
	var problem models.Problem
	if err := db.First(&problem, "id = ?", problemID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "problem not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}
	problem.Title = body.Title
	problem.Description = body.Description
	if err := db.Save(&problem).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not update problem"})
	}
	return c.JSON(http.StatusOK, problem)
}

//Delete problem in a contest
func DeleteProblem(c echo.Context) error {
	problemID := c.Param("id")
	db := config.DB

	// Check if the problem exists
	var problem models.Problem
	if err := db.First(&problem, "id = ?", problemID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "problem not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}

	// Delete the problem
	if err := db.Delete(&problem).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not delete problem"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "problem deleted successfully"})
}

// Get all test cases for a problem
func GetAllTestCasesByProblemID(c echo.Context) error {
	problemID := c.Param("id")
	db := config.DB
	var testCases []models.TestCase

	if err := db.Where("problem_id = ?", problemID).Find(&testCases).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to retrieve test cases"})
	}

	if len(testCases) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "no test cases found for this problem"})
	}

	return c.JSON(http.StatusOK, testCases)
}

//Create Testcase for a problem
func CreateTestCase(c echo.Context) error {
	problemID := c.Param("id")
	var body struct {
		Input  string `json:"input"`
		Output string `json:"output"`
	}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	db := config.DB
	var problem models.Problem

	if err := db.First(&problem, "id = ?", problemID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "problem not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}

	testCase := models.TestCase{
		ID:        uuid.New(),
		ProblemID: problem.ID,
		Input:     body.Input,
		Output:    body.Output,
	}

	if err := db.Create(&testCase).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not create test case"})
	}

	return c.JSON(http.StatusCreated, testCase)
}

//Update Testcase for a problem
func UpdateTestCase(c echo.Context) error {
	testCaseID := c.Param("id")
	db := config.DB
	var body struct {
		Input  string `json:"input"`
		Output string `json:"output"`
	}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	var testCase models.TestCase
	if err := db.First(&testCase, "id = ?", testCaseID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "test case not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}

	testCase.Input = body.Input
	testCase.Output = body.Output

	if err := db.Save(&testCase).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not update test case"})
	}

	return c.JSON(http.StatusOK, testCase)
}

// Delete Testcase for a problem
func DeleteTestCase(c echo.Context) error {
	testCaseID := c.Param("id")
	db := config.DB

	// Check if the test case exists
	var testCase models.TestCase
	if err := db.First(&testCase, "id = ?", testCaseID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "test case not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}

	// Delete the test case
	if err := db.Delete(&testCase).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not delete test case"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "test case deleted successfully"})
}

// Get all submissions for a problem
func GetAllSubmissionsByProblemID(c echo.Context) error {
	problemID := c.Param("id")
	db := config.DB
	var submissions []models.Submission

	if err := db.Where("problem_id = ?", problemID).Find(&submissions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to retrieve submissions"})
	}

	if len(submissions) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "no submissions found for this problem"})
	}

	return c.JSON(http.StatusOK, submissions)
}

// Handle submission for a problem
func HandleSubmission(c echo.Context) error {
	userID := c.Param("user_id")
	problemID := c.Param("problem_id")
	db := config.DB
	if userID == "" || problemID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user_id and problem_id are required"})
	}
	var body struct {
		SourceCode string `json:"source_code"`
		Language   string `json:"language"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	// Validate user exists
	var user models.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}
	// Validate problem exists
	var problem models.Problem
	if err := db.First(&problem, "id = ?", problemID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "problem not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}

	submission := models.Submission{
		ID:         uuid.New(),
		ProblemID:  problem.ID,
		UserID:     user.ID,
		SubmittedAt: time.Now(),
		Result:     "pending", // Initial status
		SourceCode: body.SourceCode,
		Language:  body.Language,
		Score:     0, // Initial score
	}
	if err := db.Create(&submission).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not create submission"})
	}
	// Send submission to RabbitMQ for processing
	if err := rabbitmq.SendSubmissionToQueue(submission); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to send submission to queue"})
	}
	return c.JSON(http.StatusCreated, submission)
}
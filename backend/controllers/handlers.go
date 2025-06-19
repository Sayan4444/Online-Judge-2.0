package handler

import (
	"OJ-backend/config"
	models "OJ-backend/models"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
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

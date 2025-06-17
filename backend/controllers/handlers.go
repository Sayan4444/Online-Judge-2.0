package handler

import (
    "fmt"
    "time"
    "net/http"
    "github.com/labstack/echo/v4"
    echojwt "github.com/labstack/echo-jwt/v4"
    "github.com/golang-jwt/jwt/v5"
    "OJ-backend/config"
)

var jwtSecret = []byte(config.GetEnv("JWT_SECRET")) 

type Claims struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}

func GenerateToken(username, email string) (string, error) {
    claims := Claims{
        Username: username,
        Email:    email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}


func JWTMiddleware() echo.MiddlewareFunc {
    return echojwt.WithConfig(echojwt.Config{
        SigningKey:  jwtSecret,
        ContextKey:  "user",
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
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
    }

    claims := &Claims{
        Username: body.Username,
        Email:    body.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
        },
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
	username, email := GetUserFromContext(c)

	profile := map[string]string{
		"username": username,
		"email":    email,
	}

	return c.JSON(http.StatusOK, profile)
}

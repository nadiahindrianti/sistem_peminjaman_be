package middlewares

import (
	"sistem_peminjaman_be/helpers"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func CreateToken(userID uint, role string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userID
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // token expires after 24 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}

func GetTokenFromHeader(req *http.Request) string {
	authHeader := req.Header.Get("Authorization")
	if authHeader != "" {
		// The header value should be in the format "Bearer <token>"
		splitHeader := strings.Split(authHeader, " ")
		if len(splitHeader) == 2 {
			return splitHeader[1]
		}
		return ""
	}
	return ""
}

func GetUserIdFromToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_JWT")), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid claims")
	}

	getUserId, ok := claims["userId"].(float64)
	if !ok {
		return 0, errors.New("userId claim not found")
	}
	userId := uint(getUserId)
	return userId, nil
}

func JWTErrorHandler(err error, c echo.Context) error {
	// Customize the JWT error response
	customError := helpers.ErrorResponse{
		StatusCode: http.StatusUnauthorized,
		Message:    "Unauthorized",
		Errors:     err.Error(),
	}
	return c.JSON(http.StatusUnauthorized, customError)
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_JWT")), nil
		})
		if err != nil || !token.Valid {
			return JWTErrorHandler(err, c)
		}

		// Set the validated token in the context
		c.Set("user", token)

		return next(c)
	}
}

func RoleMiddleware(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			userRole := claims["role"].(string)

			// fmt.Println("1", role)
			// fmt.Println("2", userRole)

			// Check if the user's role matches the required role
			if userRole != role {
				// Return an error response indicating unauthorized access
				errorResponse := helpers.ErrorResponse{
					StatusCode: http.StatusForbidden,
					Message:    "Forbidden",
					Errors:     "Unauthorized access",
				}
				return c.JSON(http.StatusForbidden, errorResponse)
			}

			return next(c)
		}
	}
}

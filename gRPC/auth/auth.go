package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	Email   string `json:"email"`
	User_Id int    `json:"user_id"`
	jwt.RegisteredClaims
}

type UserLog struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	secretKey = []byte("venomsnake")
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized: Missing token"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Unauthorized: Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
			c.Set("email", claims.Email)
		} else {
			c.JSON(401, gin.H{"error": "Unauthorized: Invalid claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func SetTokenCookie(c *gin.Context, token string) {
	c.SetCookie("token", token, 86400, "/", "", false, true)
}
func CreateToken(user UserLog) (string, error) {
	claims := TokenClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "GinApp",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// AuthMiddleware returns a Gin middleware function for JWT authentication
func ParseToken(c *gin.Context) int {
	// 1. Get token from cookie
	tokenString, _ := c.Cookie("token")

	// 2. Parse token (skip error checking)
	token, _ := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil // Use your real secret
	})

	// 3. Extract and return UserID
	return token.Claims.(*TokenClaims).User_Id
}

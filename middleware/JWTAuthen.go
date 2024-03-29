package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWTMiddleware(c *gin.Context) {
	hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
	header := c.Request.Header.Get("Authorization")
	tokenString := strings.Replace(header, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexception  signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if Claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("userId", Claims["userId"])
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "forbidden", "message": err.Error()})
	}

	c.Next()
}

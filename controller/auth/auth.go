package auth

import (
	"fmt"
	"net/http"
	"os"
	"projectApi/interview-api/orm"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte

type RegisterBody struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Avatar    string `json:"avatar"`
}
type LoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var json RegisterBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	user := orm.User{Username: json.Username, Password: string(encryptedPassword), Email: json.Email, Avatar: json.Avatar, FirstName: json.FirstName, LastName: json.LastName}
	//vailedte user exists
	var userExists orm.User
	orm.Db.Where("username = ? OR email = ?", json.Username, json.Email).First(&userExists)
	if userExists.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": "username is duplicate or email is already",
		})
		return
	}

	orm.Db.Create(&user)
	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"massage": "User Create Success",
			"userId":  user.ID,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": "User Create Failed",
		})
	}
}

func Login(c *gin.Context) {
	var json LoginBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check User Exists
	var userExists orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExists)
	if userExists.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "User Does Not Exists"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(userExists.Password), []byte(json.Password))
	if err == nil {
		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": userExists.ID,
			"exp":    time.Now().Add(time.Hour * 1).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		fmt.Println(tokenString, err)
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Login Success", "token": tokenString})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Login Failed"})

	}
}

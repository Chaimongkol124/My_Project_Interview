package user

import (
	"net/http"
	"projectApi/interview-api/orm"

	"github.com/gin-gonic/gin"
)

type ResponeBody struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	Avatar    string
}

func ReadAll(c *gin.Context) {
	var users []orm.User
	orm.Db.Find(&users)
	var responeUser []ResponeBody
	for _, user := range users {
		responeUser = append(responeUser, MapRespone(user.Username, user.FirstName, user.LastName, user.Email, user.Avatar))
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Get Users", "Users": responeUser})
}

func Profile(c *gin.Context) {
	userID := c.MustGet("userId").(float64)
	var user orm.User
	orm.Db.First(&user, userID)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Get Users", "User": MapRespone(user.Username, user.FirstName, user.LastName, user.Email, user.Avatar)})
}

func MapRespone(username, firstName, lastName, email, avatar string) ResponeBody {
	return ResponeBody{Username: username, FirstName: firstName, LastName: lastName, Email: email, Avatar: avatar}
}

package comment

import (
	"net/http"
	"projectApi/interview-api/orm"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Massage      string  `json:"massage" binding:"required"`
	Interview_id float64 `json:"Interview_id" binding:"required"`
}
type RequestUpdateBody struct {
	Id      float64 `json:"id" binding:"required"`
	Massage string  `json:"massage" binding:"required"`
}
type ResponseComment struct {
	Id        float64
	Massage   string
	Username  string
	CreatedAt time.Time
}

func CreateComment(c *gin.Context) {
	var json RequestBody
	userId := c.MustGet("userId").(float64)
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var interViewExists orm.Interview
	orm.Db.Where("id = ?", json.Interview_id).First(&interViewExists)
	if interViewExists.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": "invaild is interview_id",
		})
		return
	}
	comment := orm.Comment{Massage: json.Massage, InterviewID: json.Interview_id, UserID: userId}
	orm.Db.Create(&comment)
	if comment.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"massage":   "Comment Create Success",
			"commentId": comment.ID,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": "Comment Create Failed",
		})
	}
}
func UpdateComent(c *gin.Context) {
	var json RequestUpdateBody
	userId := c.MustGet("userId").(float64)
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var comment orm.Comment
	orm.Db.Where("id = ?", json.Id).First(&comment)
	if comment.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": "invaild is Comment",
		})
		return
	}
	if userId != comment.UserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": "Can't Update Comment",
		})
		return
	}
	comment.Massage = json.Massage
	comment.UpdatedAt = time.Now()
	orm.Db.Save(&comment)
	if comment.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"massage":     "Comment Update Success",
			"interviewId": comment.ID,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": "Comment Update Failed",
		})
	}
}
func DeleteComment(c *gin.Context) {
	var id = c.Param("id")

	if newId, err := strconv.ParseFloat(id, 64); err == nil {
		var comment orm.Comment
		orm.Db.Where("id = ?", newId).First(&comment)
		if comment.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"massage": "Delete Failed",
			})
			return
		}
		userID := c.MustGet("userId").(float64)
		if float64(comment.UserID) != userID {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"massage": "this user can't Delete this crad Comment",
			})
			return
		}
		orm.Db.Unscoped().Delete(&comment)

		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Delete Success", "id": newId})

	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
	}
}
func ReadAllComment(c *gin.Context) {
	var idInterview = c.Param("interview_id")
	if newId, err := strconv.ParseFloat(idInterview, 64); err == nil {
		var result []ResponseComment
		orm.Db.Model(&orm.Comment{}).Select("users.username ,users.email,comments.created_at,comments.massage").Joins("JOIN users ON users.id = comments.user_id").Where("interview_id =?", newId).Find(&result)
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"massage": "ReadAll Comment Success",
			"data":    result,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	}
}

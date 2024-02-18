package interview

import (
	"net/http"
	"projectApi/interview-api/orm"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Description string `json:"description" binding:"required"`
}
type RequestUpdateBody struct {
	Id          float64 `json:"id" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Status      string  `json:"status" binding:"required"`
}
type ResponseInterview struct {
	Id          float64
	Name        string
	Description string
	Username    string
	Email       string
	CreatedAt   time.Time
	Status      string
}

func CreateIneterview(c *gin.Context) {
	var json RequestBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := c.MustGet("userId").(float64)
	interview := orm.Interview{Description: json.Description, Status: orm.StausType("To Do"), UserID: userId, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	//vaileate Fk_user_id interview

	orm.Db.Create(&interview)

	if interview.ID > 0 {
		logHistory := orm.HistoryInterview{Name: "นัดสัมภาษณ์งาน " + strconv.FormatUint(uint64(interview.ID), 10), Description: interview.Description, Status: interview.Status, InterviewID: float64(interview.ID)}
		orm.Db.Create(&logHistory)
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"massage":     "Interview Create Success",
			"interviewId": interview.ID,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": "Interview Create Failed",
		})
	}
}

func UpdateIneterview(c *gin.Context) {
	var json RequestUpdateBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var interview orm.Interview
	orm.Db.Where("id = ?", json.Id).First(&interview)
	if interview.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": "invaild is interview",
		})
		return
	}
	userID := c.MustGet("userId").(float64)

	if float64(interview.UserID) != userID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": "this user can't update this crad interview",
		})
		return
	}
	interview.Description = json.Description
	interview.Status = orm.StausType(json.Status)
	interview.UpdatedAt = time.Now()

	orm.Db.Save(&interview)
	if interview.ID > 0 {
		logHistory := orm.HistoryInterview{Name: "นัดสัมภาษณ์งาน " + strconv.FormatUint(uint64(interview.ID), 10), Description: interview.Description, Status: interview.Status, InterviewID: float64(interview.ID)}
		orm.Db.Create(&logHistory)
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"massage":     "Interview Update Success",
			"interviewId": interview.ID,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": "Interview Update Failed",
		})
	}
}

func DeleteIneterview(c *gin.Context) {
	var id = c.Param("id")

	if newId, err := strconv.ParseFloat(id, 64); err == nil {
		var interview orm.Interview
		orm.Db.Where("id = ?", newId).First(&interview)
		if interview.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"massage": "Delete Failed",
			})
			return
		}
		userID := c.MustGet("userId").(float64)
		if interview.UserID != userID {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"massage": "this user can't Delete this crad interview",
			})
			return
		}
		orm.Db.Unscoped().Delete(&interview)
		//orm.Db.Model(orm.Comment{}).Delete(&interview)
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Delete Success", "id": newId})

	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
	}
}

func GetInterView(c *gin.Context) {
	var id = c.Param("id")
	if newId, err := strconv.ParseFloat(id, 64); err == nil {
		var interview orm.Interview
		orm.Db.Where("id = ?", newId).First(&interview)
		if interview.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"massage": "Interview is not found",
			})
			return
		}
		var user orm.User
		orm.Db.Where("id = ?", interview.UserID).First(&user)

		c.JSON(http.StatusOK, gin.H{
			"status":       "ok",
			"massage":      "Get Interview  Success",
			"data":         ResponseInterview{Id: float64(interview.ID), Name: "นัดสัมภาษณ์งาน " + strconv.FormatUint(uint64(interview.ID), 10), Description: interview.Description, Username: user.Username, Email: user.Email, CreatedAt: interview.CreatedAt, Status: string(interview.Status)},
			"optionStatus": []string{"To Do", "In Progress", "Done"},
		})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	}
}

func ReadAllInterViews(c *gin.Context) {
	var result []ResponseInterview
	orm.Db.Model(&orm.Interview{}).Select("CONCAT (\"นัดสัมภาษณ์งาน \", interviews.id) AS Name,interviews.id AS ID,interviews.Description,interviews.status,users.username ,users.email,interviews.created_at").Joins("JOIN users ON users.id = interviews.user_id").Find(&result)
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"massage": "ReadAll Interview Success",
		"Data":    result,
	})
}

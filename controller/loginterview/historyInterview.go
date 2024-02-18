package loginterview

import (
	"net/http"
	"projectApi/interview-api/orm"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	Description string
	Name        string
	Status      string
}

func ReadAllHistoryInteriews(c *gin.Context) {
	var result []ResponseBody
	var idInterview = c.Param("interview_id")
	if newId, err := strconv.ParseFloat(idInterview, 64); err == nil {
		orm.Db.Model(&orm.HistoryInterview{}).Select("history_interviews.description,history_interviews.name,history_interviews.status").Where("interview_id =?", newId).Find(&result)
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"massage": "ReadAll HistoryInterview Success",
			"Data":    result,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	}
}

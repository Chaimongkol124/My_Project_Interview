package main

import (
	"fmt"
	"log"
	AuthController "projectApi/interview-api/controller/auth"
	CommentController "projectApi/interview-api/controller/comment"
	InterviewController "projectApi/interview-api/controller/interview"
	HsitoryInterviewController "projectApi/interview-api/controller/loginterview"
	UserController "projectApi/interview-api/controller/user"
	"projectApi/interview-api/middleware"
	"projectApi/interview-api/orm"
	"time"

	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/ratelimit"
)

var limit ratelimit.Limiter

func leakBucket() gin.HandlerFunc {
	prev := time.Now()
	return func(ctx *gin.Context) {
		now := limit.Take()
		log.Print(color.CyanString("%v", now.Sub(prev)))
		prev = now
	}
}
func main() {
	limit = ratelimit.New(100)
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Print("Error loading .env file")
	}
	orm.InitDb()
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(leakBucket())
	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)
	authorized := r.Group("/api", middleware.JWTMiddleware)
	authorized.GET("/user/getAll", UserController.ReadAll)
	authorized.GET("/user/profile", UserController.Profile)
	authorized.POST("/interview/create", InterviewController.CreateIneterview)
	authorized.PUT("/interview/update", InterviewController.UpdateIneterview)
	authorized.DELETE("/interview/delete/:id", InterviewController.DeleteIneterview)
	authorized.GET("/interview/detail/:id", InterviewController.GetInterView)
	authorized.GET("/interview/readall", InterviewController.ReadAllInterViews)
	authorized.GET("/comment/readall/:interview_id", CommentController.ReadAllComment)
	authorized.POST("/comment/create", CommentController.CreateComment)
	authorized.PUT("/comment/update", CommentController.UpdateComent)
	authorized.DELETE("/comment/delete/:id", CommentController.DeleteComment)
	authorized.GET("/historyInterview/readall/:interview_id", HsitoryInterviewController.ReadAllHistoryInteriews)
	r.Run("localhost:8080")
}

package main

import (
	"testDans/config"
	"testDans/handler"
	"testDans/jwt"

	"github.com/gin-gonic/gin"
)

func main() {
	// init db and framework
	db := config.InitGorm()
	r := gin.Default()

	// available endpoint
	handler := handler.NewHandler(db)
	r.GET("/", handler.Check)
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	withJwt := r.Group("/job")
	withJwt.Use(jwt.JwtAuthMiddleware())
	withJwt.GET("/list", handler.GetJobList)
	withJwt.GET("/detail", handler.GetJobDetail)

	// serve on 8080
	r.Run(":8080")
}

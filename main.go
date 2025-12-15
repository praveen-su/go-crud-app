package main

import (
	"context"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

func main() {
	InitMySQL()
	InitMongo()
	InitRedis()

	StartGrpcServer()

	r := gin.Default()
	r.POST("/users", CreateUser)
	r.GET("/users/:id", GetUser)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)

	r.Run(":8080")
}

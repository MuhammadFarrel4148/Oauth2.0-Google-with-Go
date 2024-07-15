package main

import (
	"testingoauth/database"
	"testingoauth/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	database.DatabaseConnect()

	router.GET("/login", handler.Login)
	router.GET("/callback", handler.GoogleCallBack)

	router.Run()
}
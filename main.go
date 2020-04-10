package main

import (
	handler "ginhello/handlers"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	router = gin.Default()

	router.LoadHTMLGlob("templates/*")

	initializeRoutes()

	router.Run()
}

func initializeRoutes() {
	router.GET("/", handler.ShowIndexPage)
	router.GET("article/view/:article_id", handler.GetArticle)
}

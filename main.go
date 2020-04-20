package main

import (
	handler "ginhello/handlers"
	model "ginhello/models"

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
	var err error
	db, err := model.NewDB("./models/gorm.db")
	if err != nil {
		panic("failed to connect database")
	}
	env := &handler.Env{DB: db}

	router.GET("/", env.ShowIndexPage)
	router.GET("article/view/:article_id", env.GetArticle)
}

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
	db.AutoMigrate(&model.User{})

	env := &handler.Env{DB: db}

	router.GET("/", env.ShowIndexPage)
	router.GET("/article/view/:article_id", env.GetArticle)
	router.GET("/article/new", env.CreateArticle)
	router.POST("/article/new", env.CreateArticle)

	router.GET("/signin", env.Signin)
	router.POST("/signin", env.Signin)

	router.GET("/signup", env.Signup)
	router.POST("/signup", env.Signup)

}

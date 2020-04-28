package main

import (
	handler "ginhello/handlers"
	model "ginhello/models"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var router *gin.Engine

func main() {
	port := os.Getenv("PORT")

	router = gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/css", "./templates/css")
	router.Static("/img", "./templates/img")
	initializeRoutes()

	router.Run(":" + port)
}

func initializeRoutes() {
	var err error
	db, err := model.NewDB("./models/gorm.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Article{}, &model.User{})

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

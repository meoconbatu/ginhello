package main

import (
	"encoding/gob"
	handler "ginhello/handlers"
	model "ginhello/models"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var router *gin.Engine

func main() {
	port := os.Getenv("PORT")

	router = gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	gob.Register(model.User{})

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

	router.GET("/signin", env.Signin)
	router.POST("/signin", env.Signin)

	router.GET("/signup", env.Signup)
	router.POST("/signup", env.Signup)

	router.Use(authRequired())
	{
		router.GET("/article/view", env.GetArticles)
		router.GET("/article/view/:article_id", env.GetArticle)
		router.GET("/article/new", env.CreateArticle)
		router.POST("/article/new", env.CreateArticle)
	}

}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := isAuthenticated(c)
		if currentUser == nil {
			c.Redirect(http.StatusSeeOther, "/signin")
			c.Abort()
			return
		}
		c.Next()
	}
}
func isAuthenticated(c *gin.Context) *model.User {
	session := sessions.Default(c)
	sUser := session.Get("user")
	if sUser == nil {
		return nil
	}
	user, ok := sUser.(model.User)
	if !ok {
		return nil
	}
	return &user
}

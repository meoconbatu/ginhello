package handler

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"strconv"

	model "ginhello/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Env struct
type Env struct {
	DB model.DataStore
}

// DB variable

// ShowIndexPage func
func (env *Env) ShowIndexPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Home Page"}, "index.html")
}

// GetArticles func
func (env *Env) GetArticles(c *gin.Context) {
	articles := env.DB.GetAllArticles()
	render(c, gin.H{
		"title":   "Home Page",
		"payload": articles}, "articles.html")
}

// GetArticle func
func (env *Env) GetArticle(c *gin.Context) {
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		if article, err := env.DB.GetArticleByID(articleID); err == nil {
			render(c, gin.H{
				"title":   "Home Page",
				"payload": article}, "article.html")
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
}

// CreateArticle func
func (env *Env) CreateArticle(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		render(c, gin.H{"title": "Home Page"}, "new.html")
	case "POST":
		var article model.Article
		if c.ShouldBind(&article) != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		err := env.DB.CreateArticle(&article)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}

// Signin func
func (env *Env) Signin(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		render(c, gin.H{"title": "Home Page"}, "signin.html")
	case "POST":
		username := c.PostForm("username")
		password := c.PostForm("password")
		err := env.DB.AuthenticateUser(username, password)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		session := sessions.Default(c)
		session.Set("user", model.User{Username: username})
		err = session.Save()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}

// Signup func
func (env *Env) Signup(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		render(c, gin.H{"title": "Home Page"}, "signup.html")
	case "POST":
		var user model.User
		if c.ShouldBind(&user) != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if env.DB.CreateUser(&user) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}

// Signout func
func (env *Env) Signout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	err := session.Save()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// SetupRouter func
func SetupRouter(env *Env) *gin.Engine {
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	gob.Register(model.User{})
	var accessPath string
	if gin.Mode() == gin.TestMode {
		accessPath = "../"
	} else {
		accessPath = "./"
	}
	router.LoadHTMLGlob(accessPath + "templates/*.html")
	router.Static("/css", accessPath+"templates/css")
	router.Static("/img", accessPath+"templates/img")

	router.GET("/", env.ShowIndexPage)

	router.GET("/signin", env.Signin)
	router.POST("/signin", env.Signin)

	router.GET("/signup", env.Signup)
	router.POST("/signup", env.Signup)

	router.GET("signout", env.Signout)

	router.Use(authRequired())
	{
		router.GET("/article/view", env.GetArticles)
		router.GET("/article/view/:article_id", env.GetArticle)
		router.GET("/article/new", env.CreateArticle)
		router.POST("/article/new", env.CreateArticle)
	}
	return router
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

func render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		c.XML(http.StatusOK, data["payload"])
	default:
		session := sessions.Default(c)
		user := session.Get("user")
		if user != nil {
			u, ok := user.(model.User)
			if ok {
				data["user"] = u
			}
		}
		c.HTML(http.StatusOK, templateName, data)
	}
}

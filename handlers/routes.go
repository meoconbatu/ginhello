package handler

import (
	"net/http"
	"strconv"

	model "ginhello/models"

	"github.com/gin-gonic/gin"
)

// Env struct
type Env struct {
	DB model.DataStore
}

// DB variable

// ShowIndexPage func
func (env *Env) ShowIndexPage(c *gin.Context) {
	articles := env.DB.GetAllArticles()
	render(c, gin.H{
		"title":   "Home Page",
		"payload": articles}, "index.html")
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
		c.Redirect(http.StatusSeeOther, "/")
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
		c.Redirect(http.StatusSeeOther, "/")
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
		c.Redirect(http.StatusSeeOther, "/")
	}
}
func render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		c.XML(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}

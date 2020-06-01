package handler

import (
	"encoding/gob"
	"net/http"

	model "ginhello/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Env struct
type Env struct {
	DB model.DataStore
}

// ShowIndexPage func
func (env *Env) ShowIndexPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Home Page"}, "index.html")
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

	router.GET("/auth/google/signin", env.SigninWithSocial(confGoogle))
	router.GET("/auth/google/callback", env.SigninWithSocialCallback(confGoogle, GOOGLE_USERINFO_ENDPOINT))

	router.GET("/auth/github/signin", env.SigninWithSocial(confGithub))
	router.GET("/auth/github/callback", env.SigninWithSocialCallback(confGithub, GITHUB_USERINFO_ENDPOINT))

	router.GET("/signinfail", env.SigninFail)

	router.GET("/signup", env.Signup)
	router.POST("/signup", env.Signup)
	router.GET("/signupsuccess", env.SignupSuccess)

	router.GET("/verify", env.Verify)
	router.GET("/verificationsuccess", env.VerificationSuccess)

	router.GET("/signout", env.Signout)

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

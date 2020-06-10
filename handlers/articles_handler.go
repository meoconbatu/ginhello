package handler

import (
	model "ginhello/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	case http.MethodGet:
		render(c, gin.H{"title": "Home Page"}, "new.html")
	case http.MethodPost:
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

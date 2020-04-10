package handler

import (
	"net/http"
	"strconv"

	model "ginhello/models"

	"github.com/gin-gonic/gin"
)

// ShowIndexPage func
func ShowIndexPage(c *gin.Context) {
	render(c, gin.H{
		"title":   "Home Page",
		"payload": model.GetAllArticles()}, "index.html")
}

// GetArticle func
func GetArticle(c *gin.Context) {
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		if article, err := model.GetArticleByID(articleID); err == nil {
			render(c, gin.H{
				"title":   "Home Page",
				"payload": article}, "article.html")
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
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

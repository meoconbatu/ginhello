package model

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Article type
type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var articleList = []Article{
	{ID: 1, Title: "Article 1", Content: "Article 1 body"},
	{ID: 2, Title: "Article 2", Content: "Article 2 body"},
}
var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "./gorm.db")
	// defer db.Close()
	if err != nil {
		panic("failed to connect database")
	}
	// db.AutoMigrate(&Article{})
}

// GetAllArticles Return a list of all the articles
func GetAllArticles() []Article {
	var articles []Article
	db.Find(&articles)
	return articles
}

// GetArticleByID func
func GetArticleByID(id int) (*Article, error) {
	var article Article

	db.First(&article, id)

	if article.ID == 0 {
		return nil, errors.New("Article not found")
	}
	return &article, nil
}

// CreateArticle func
func CreateArticle(article *Article) (int, error) {
	db.Create(&article)
	if (*article).ID == 0 {
		return 0, errors.New("Error when create article")
	}
	return (*article).ID, nil
}

// DeleteArticleByID func
func DeleteArticleByID(id int) {
	db.Where("id = ?", id).Delete(&Article{})
}

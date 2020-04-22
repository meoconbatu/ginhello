package model

import (
	"errors"
)

// Article type
type Article struct {
	ID      int
	Title   string `form:"title"`
	Content string `form:"content"`
}

// GetAllArticles Return a list of all the articles
func (db *DB) GetAllArticles() []Article {
	var articles []Article
	db.Find(&articles)
	return articles
}

// GetArticleByID func
func (db *DB) GetArticleByID(id int) (*Article, error) {
	var article Article

	db.First(&article, id)

	if article.ID == 0 {
		return nil, errors.New("Article not found")
	}
	return &article, nil
}

// CreateArticle func
func (db *DB) CreateArticle(article *Article) error {
	db.Create(&article)
	if (*article).ID == 0 {
		return errors.New("Error when create article")
	}
	return nil
}

// DeleteArticleByID func
func (db *DB) DeleteArticleByID(id int) {
	db.Where("id = ?", id).Delete(&Article{})
}

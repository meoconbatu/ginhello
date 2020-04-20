package model

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DataStore type
type DataStore interface {
	GetAllArticles() []Article
	GetArticleByID(id int) (*Article, error)
	CreateArticle(article *Article) (int, error)
	DeleteArticleByID(id int)
}

// DB type
type DB struct {
	*gorm.DB
}

// NewDB func
func NewDB(dataSourceName string) (*DB, error) {
	db, err := gorm.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, errors.New("failed to connect database")
	}
	return &DB{db}, nil
}

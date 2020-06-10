package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var db *DB

func init() {
	var err error
	db, err = NewDB("../models/gorm.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Article{}, &User{})
}
func createArticleTest(t *testing.T) *Article {
	articleTest := &Article{Title: "test", Content: "test"}

	err := db.CreateArticle(articleTest)

	assert.NoError(t, err)

	t.Cleanup(func() {
		db.DeleteArticleByID(articleTest.ID)
	})
	return articleTest
}
func TestCreateAndGetArticle(t *testing.T) {
	articleTest := createArticleTest(t)

	article, err := db.GetArticleByID(articleTest.ID)

	assert.NoError(t, err)
	assert.Equal(t, articleTest.ID, article.ID)
	assert.Equal(t, articleTest, article)
}

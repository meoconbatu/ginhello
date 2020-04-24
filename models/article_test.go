package model

import (
	"reflect"
	"testing"
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
	if err != nil {
		t.Fatalf("Failed to create article: %s\n", err.Error())
	}
	t.Cleanup(func() {
		db.DeleteArticleByID(articleTest.ID)
	})
	return articleTest
}
func TestGetArticle(t *testing.T) {
	articleTest := createArticleTest(t)
	article, err := db.GetArticleByID(articleTest.ID)
	if err != nil {
		t.Fatalf("Failed to get article: %s\n", err.Error())
	}
	if article.ID != articleTest.ID {
		t.Fatalf("Article Ids do not match: %d\n vs %d\n", article.ID, articleTest.ID)
	}
	if reflect.DeepEqual(article, articleTest) != true {
		t.Fatalf("Articles do not match: %+v\n vs %+v\n", article, articleTest)
	}
}

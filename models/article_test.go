package model

import (
	"reflect"
	"testing"
)

func createArticleTest(t *testing.T) *Article {
	articleTest := &Article{Title: "test", Content: "test"}
	id, err := CreateArticle(articleTest)
	if err != nil {
		t.Errorf("Failed to create article: %s\n", err.Error())
	}
	articleTest.ID = id
	t.Cleanup(func() {
		DeleteArticleByID(articleTest.ID)
	})
	return articleTest
}
func TestGetArticle(t *testing.T) {
	articleTest := createArticleTest(t)
	article, err := GetArticleByID(articleTest.ID)
	if err != nil {
		t.Errorf("Failed to get article: %s\n", err.Error())
	}
	if article.ID != articleTest.ID {
		t.Errorf("Article Ids do not match: %d\n vs %d\n", article.ID, articleTest.ID)
	}
	if reflect.DeepEqual(article, articleTest) != true {
		t.Errorf("Articles do not match: %+v\n vs %+v\n", article, articleTest)
	}
}

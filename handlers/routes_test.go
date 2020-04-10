package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShowIndexPageUnauthenticated(t *testing.T) {
	router := getRouter(true)
	router.GET("/", ShowIndexPage)

	req, _ := http.NewRequest("GET", "/", nil)

	testHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Home Page</title>") > 0
		return statusOK && pageOK
	})
}
func TestArticleUnauthenticated(t *testing.T) {
	router := getRouter(true)
	router.GET("article/view/:article_id", GetArticle)

	req, _ := http.NewRequest("GET", "article/view/1", nil)

	testHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Article 1</title>") > 0
		return statusOK && pageOK
	})
}

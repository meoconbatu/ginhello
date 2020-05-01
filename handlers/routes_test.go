package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var env Env

func init() {
	env = Env{&mockDB{}}
}
func TestShowIndexPageUnauthenticated(t *testing.T) {
	router := SetupRouter(&env)

	req, _ := http.NewRequest("GET", "/", nil)

	testHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Home Page</title>") > 0
		return statusOK && pageOK
	})
}
func TestArticleUnauthenticated(t *testing.T) {
	router := SetupRouter(&env)

	req, _ := http.NewRequest("GET", "/article/view/1", nil)

	testHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusSeeOther
		response := w.Result()
		location := response.Header.Get("Location")
		pageOK := location == "/signin"
		return statusOK && pageOK
	})
}

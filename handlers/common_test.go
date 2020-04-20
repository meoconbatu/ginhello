package handler

import (
	model "ginhello/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"os"

	"github.com/gin-gonic/gin"
)

type mockDB struct{}

func (mdb *mockDB) GetAllArticles() []model.Article {
	return []model.Article{
		{ID: 1, Title: "test", Content: "test"},
		{ID: 2, Title: "test", Content: "test"},
	}
}

func (mdb *mockDB) GetArticleByID(id int) (*model.Article, error) {
	return &model.Article{ID: 1, Title: "test", Content: "test"}, nil
}

func (mdb *mockDB) CreateArticle(article *model.Article) (int, error) {
	return 1, nil
}

func (mdb *mockDB) DeleteArticleByID(id int) {

}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
func getRouter(withTemplate bool) *gin.Engine {
	router := gin.Default()
	if withTemplate {
		router.LoadHTMLGlob("../templates/*")
	}
	return router

}
func testHTTPResponse(t *testing.T, router *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if !f(w) {
		t.Fail()
	}
}

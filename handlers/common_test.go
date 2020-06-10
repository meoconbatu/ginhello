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

var (
	env    Env
	router *gin.Engine
)

func (mdb *mockDB) GetAllArticles() []model.Article {
	return []model.Article{
		{ID: 1, Title: "test", Content: "test"},
		{ID: 2, Title: "test", Content: "test"},
	}
}

func (mdb *mockDB) GetArticleByID(id int) (*model.Article, error) {
	return &model.Article{ID: 1, Title: "test", Content: "test"}, nil
}

func (mdb *mockDB) CreateArticle(article *model.Article) error {
	return nil
}

func (mdb *mockDB) DeleteArticleByID(id int) {

}
func (mdb *mockDB) AuthenticateUser(username, password string) error {
	return nil
}

func (mdb *mockDB) CreateUser(user *model.User) error {
	return nil
}
func (mdb *mockDB) VerifyEmail(userid int) error {
	return nil
}

func (mdb *mockDB) GetVerificationToken(token string) *model.VerificationToken {
	return nil
}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	env = Env{&mockDB{}}
	router = SetupRouter(&env)
	os.Exit(m.Run())
}
func testHTTPResponse(t *testing.T, router *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if !f(w) {
		t.Fail()
	}
}
func performRequest(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"os"

	"github.com/gin-gonic/gin"
)

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

package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowIndexPage(t *testing.T) {
	w := performRequest(router, http.MethodGet, "/")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Home Page")
}
func TestRoutesUnauthenticated(t *testing.T) {
	testRoutes := []struct {
		name  string
		route string
	}{
		{"view list articles", "/article/view"},
		{"view an article", "/article/view/1"},
		{"create an article", "/article/new"},
	}
	for _, tr := range testRoutes {
		t.Run(tr.name, func(t *testing.T) {
			w := performRequest(router, http.MethodGet, tr.route)
			assert.Equal(t, http.StatusSeeOther, w.Code)
			assert.Equal(t, "/signin", w.Header().Get("Location"))
		})
	}
}

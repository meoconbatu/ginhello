package handler

import (
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Credentials type
type Credentials struct {
	Cid     string
	Csecret string
}

var conf *oauth2.Config

func init() {
	var c Credentials
	c.Cid = os.Getenv("CLIENT_ID")
	c.Csecret = os.Getenv("CLIENT_SECRET")

	var server string
	if gin.Mode() == gin.ReleaseMode {
		server = os.Getenv("HOST")
	} else {
		server = os.Getenv("HOST") + ":" + os.Getenv("PORT")
	}
	conf = &oauth2.Config{
		ClientID:     c.Cid,
		ClientSecret: c.Csecret,
		RedirectURL:  server + "/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

package handler

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"golang.org/x/oauth2/github"
)

// Credentials type
type Credentials struct {
	Cid     string
	Csecret string
}

var (
	confGoogle *oauth2.Config
	confGithub *oauth2.Config
)

func init() {

	var server string
	if gin.Mode() == gin.ReleaseMode {
		server = os.Getenv("HOST")
	} else {
		server = os.Getenv("HOST") + ":" + os.Getenv("PORT")
	}
	cGoogle := Credentials{Cid: os.Getenv("GOOGLE_CLIENT_ID"),
		Csecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}
	confGoogle = &oauth2.Config{
		ClientID:     cGoogle.Cid,
		ClientSecret: cGoogle.Csecret,
		RedirectURL:  server + "/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	cGithub := Credentials{Cid: os.Getenv("GITHUB_CLIENT_ID"),
		Csecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}
	confGithub = &oauth2.Config{
		ClientID:     cGithub.Cid,
		ClientSecret: cGithub.Csecret,
		RedirectURL:  server + "/auth/github/callback",
		Endpoint:     github.Endpoint,
	}
}
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

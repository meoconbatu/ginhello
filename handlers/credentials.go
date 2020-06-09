package handler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	model "ginhello/models"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"

	"golang.org/x/oauth2/github"
)

// Credentials type
type Credentials struct {
	Cid     string
	Csecret string
}

var (
	confGoogle   *oauth2.Config
	confGithub   *oauth2.Config
	confFacebook *oauth2.Config
)

const (
	GITHUB_USERINFO_ENDPOINT   = "https://api.github.com/user"
	GOOGLE_USERINFO_ENDPOINT   = "https://www.googleapis.com/oauth2/v3/userinfo"
	FACEBOOK_USERINFO_ENDPOINT = "https://graph.facebook.com/me"
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
	cFacebook := Credentials{Cid: os.Getenv("FACEBOOK_CLIENT_ID"),
		Csecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
	}
	confFacebook = &oauth2.Config{
		ClientID:     cFacebook.Cid,
		ClientSecret: cFacebook.Csecret,
		RedirectURL:  server + "/auth/facebook/callback",
		Endpoint:     facebook.Endpoint,
	}
}
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// SigninWithSocial func
func (env *Env) SigninWithSocial(conf *oauth2.Config) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		session := sessions.Default(c)
		state := session.Get("state")
		c.Redirect(http.StatusSeeOther, conf.AuthCodeURL(fmt.Sprintf("%s", state)))
	})

}

// SigninWithSocialCallback func
func (env *Env) SigninWithSocialCallback(conf *oauth2.Config, userinfoURL string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		session := sessions.Default(c)
		state := session.Get("state")
		if state != c.Query("state") {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Invalid session state: %s", state))
			return
		}

		tok, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		client := conf.Client(oauth2.NoContext, tok)
		resp, err := client.Get(userinfoURL)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer resp.Body.Close()

		user := model.User{}
		if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		err = env.DB.CreateUser(&user)
		if err != nil && err != model.ErrUserAlreadyExists {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		session.Set("user", model.User{Username: user.Username, VerificationTokens: []model.VerificationToken{}})
		session.Options(sessions.Options{MaxAge: 0, Path: "/", HttpOnly: true})
		err = session.Save()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusSeeOther, "/")
	})

}

package handler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	model "ginhello/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/oauth2"
)

// UserLogin struct
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Signin func
func (env *Env) Signin(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		session := sessions.Default(c)
		state := randToken()
		session.Set("state", state)
		err := session.Save()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		render(c, gin.H{"title": "Home Page"}, "signin.html")
	case "POST":
		userLogin := UserLogin{}
		err := c.BindJSON(&userLogin)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = env.DB.AuthenticateUser(userLogin.Username, userLogin.Password)
		if err == model.ErrEmailNotVerified {
			// c.Redirect(http.StatusSeeOther, "signinfail")
			c.JSON(http.StatusOK, gin.H{"status": "seeother", "redirect": "/signinfail"})
			return
		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "unauthorized", "redirect": "/signinfail"})
			// c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		session := sessions.Default(c)
		session.Set("user", model.User{Username: userLogin.Username, VerificationTokens: []model.VerificationToken{}})
		err = session.Save()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// c.Redirect(http.StatusSeeOther, "/")
		c.JSON(http.StatusOK, gin.H{"status": "ok", "redirect": "/"})
	}
}

// SigninWithGoogle func
func (env *Env) SigninWithGoogle(c *gin.Context) {
	session := sessions.Default(c)
	state := session.Get("state")
	fmt.Println(state)
	c.Redirect(http.StatusSeeOther, conf.AuthCodeURL(fmt.Sprintf("%s", state)))
}
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// SigninWithGoogleCallback func
func (env *Env) SigninWithGoogleCallback(c *gin.Context) {
	session := sessions.Default(c)
	state := session.Get("state")
	fmt.Println(state)
	if state != c.Query("state") {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Invalid session state: %s", state))
		return
	}

	tok, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	fmt.Println(tok)
	client := conf.Client(oauth2.NoContext, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
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
	err = session.Save()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

// SigninFail func
func (env *Env) SigninFail(c *gin.Context) {
	render(c, gin.H{"title": "Home Page"}, "signin_fail.html")
}

// Signup func
func (env *Env) Signup(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		render(c, gin.H{"title": "Home Page"}, "signup.html")
	case "POST":
		var user model.User
		if c.ShouldBind(&user) != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		token := uuid.New()
		verificationToken := model.VerificationToken{Token: token.String(), ExpiryDate: time.Now().Local().Add(24 * time.Hour)}
		user.VerificationTokens = append(user.VerificationTokens, verificationToken)
		if env.DB.CreateUser(&user) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		sendMail(user, token.String())
		c.Redirect(http.StatusSeeOther, "/signupsuccess")
	}
}

// SignupSuccess func
func (env *Env) SignupSuccess(c *gin.Context) {
	render(c, gin.H{"title": "Home Page"}, "signup_success.html")
}

//Verify func
func (env *Env) Verify(c *gin.Context) {
	token := c.Query("token")

	verificationToken := env.DB.GetVerificationToken(token)
	if verificationToken == nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if verificationToken.ExpiryDate.Before(time.Now()) {
		c.AbortWithStatus(http.StatusBadRequest) // should redirect to resend another email verification
		return
	}
	if err := env.DB.EnableUser(verificationToken.UserID); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusSeeOther, "/verificationsuccess")
}

// VerificationSuccess func
func (env *Env) VerificationSuccess(c *gin.Context) {
	render(c, gin.H{"title": "Home Page"}, "verification_success.html")
}

// Signout func
func (env *Env) Signout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	err := session.Save()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func sendMail(user model.User, token string) {
	from := mail.NewEmail("gin-hello", "meo.con.batu1111@gmail.com")

	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.SetTemplateID("d-d2a069ef0d344ac39d27877770d1585d")
	message.Subject = "Please confirm your email"
	p := mail.NewPersonalization()

	p.Subject = "Please confirm your email"

	tos := []*mail.Email{
		mail.NewEmail(user.Username, user.Email),
	}
	p.AddTos(tos...)

	p.SetDynamicTemplateData("Username", user.Username)
	var server string
	if gin.Mode() == gin.ReleaseMode {
		server = os.Getenv("HOST")
	} else {
		server = os.Getenv("HOST") + ":" + os.Getenv("PORT")
	}
	p.SetDynamicTemplateData("URL", server+"/verify?token="+token)

	message.AddPersonalizations(p)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

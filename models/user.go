package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// User type
type User struct {
	ID                 int
	Username           string `form:"username" json:"name"`
	Password           string `form:"password"`
	Email              string `form:"email" json:"email"`
	Enabled            bool   `json:"email_verified"`
	VerificationTokens []VerificationToken
}

var (
	// ErrUserNotFound var
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailNotVerified var
	ErrEmailNotVerified = errors.New("email not verified")
	// ErrUserAlreadyExists var
	ErrUserAlreadyExists = errors.New("username already exists")
	// ErrWrongPassword var
	ErrWrongPassword = errors.New("wrong password")
)

// AuthenticateUser func
func (db *DB) AuthenticateUser(username, password string) error {
	user, err := db.getUserByUsername(username)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return ErrWrongPassword
	}
	if !user.Enabled {
		return ErrEmailNotVerified
	}
	return nil
}

// GetUserByUsername func
func (db *DB) getUserByUsername(username string) (*User, error) {
	var user User
	db.First(&user, User{Username: username})
	if user.ID == 0 {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

// CreateUser func
func (db *DB) CreateUser(user *User) error {
	if db.existUsername(user.Username) {
		return ErrUserAlreadyExists
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Error when create user")
	}
	user.Password = string(hashedPassword)
	db.Create(&user)
	if user.ID == 0 {
		return errors.New("Error when create user")
	}
	return nil
}

// EnableUser func
func (db *DB) EnableUser(userid int) error {
	user := User{ID: userid}
	db.Model(&user).Update("Enabled", true)
	return nil
}

// ExistUsername func
func (db *DB) existUsername(username string) bool {
	var user User
	db.First(&user, User{Username: username})
	if user.ID != 0 {
		return true
	}
	return false
}

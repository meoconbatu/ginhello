package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// User type
type User struct {
	ID       int
	Username string `form:"username"`
	Password string `form:"password"`
}

// AuthenticateUser func
func (db *DB) AuthenticateUser(username, password string) error {
	user, err := db.getUserByUsername(username)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

// GetUserByUsername func
func (db *DB) getUserByUsername(username string) (*User, error) {
	var user User
	db.First(&user, User{Username: username})
	if user.ID == 0 {
		return nil, errors.New("User not found")
	}
	return &user, nil
}

// CreateUser func
func (db *DB) CreateUser(user *User) error {
	err := db.existUsername(user.Username)
	if err != nil {
		return err
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

// ExistUsername func
func (db *DB) existUsername(username string) error {
	var user User
	db.First(&user, User{Username: username})
	if user.ID != 0 {
		return errors.New("Username already exists")
	}
	return nil
}
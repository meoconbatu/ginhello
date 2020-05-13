package model

import (
	"time"
)

// VerificationToken type
type VerificationToken struct {
	ID         int
	Token      string
	ExpiryDate time.Time
	UserID     int
}

// GetVerificationToken func
func (db *DB) GetVerificationToken(token string) *VerificationToken {
	var verificationToken VerificationToken
	db.First(&verificationToken, VerificationToken{Token: token})
	if verificationToken.ID == 0 {
		return nil
	}
	return &verificationToken
}

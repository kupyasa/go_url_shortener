package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func (c *JWTCustomClaims) Valid() error {
	// Check if the expiration time has passed
	if time.Now().After(time.Unix(c.ExpiresAt.Unix(), 0)) {
		return jwt.ErrTokenExpired
	}

	// Add any additional custom validation logic here

	return nil
}

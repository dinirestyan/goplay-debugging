package models

import (
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type Login struct {
	Email    string `json:"email" validate:"required" binding:"required"`
	Password string `json:"password" validate:"required" binding:"required"`
}

type Token struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	*jwt.StandardClaims
}

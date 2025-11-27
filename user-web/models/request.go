package models

import (
	"github.com/dgrijalva/jwt-go"
)

// payload
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}

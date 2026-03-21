package models

import (
	"encoding/json"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func ToStringLog(v interface{}) {
	b, _ := json.Marshal(&v)
	fmt.Println(string(b))
}

// payload
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}

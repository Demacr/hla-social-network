package models

import "github.com/dgrijalva/jwt-go"

type Profile struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Age       int    `json:"age"`
	Sex       string `json:"sex"`
	Interests string `json:"interests"`
	City      string `json:"city"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	UserId uint
	jwt.StandardClaims
}

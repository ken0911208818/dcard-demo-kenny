package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	userId string
	jwt.StandardClaims
}

func main() {
	secretkey := []byte("secret")
	claims := Claims{
		userId: "123sss",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(secretkey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(signedString)
}

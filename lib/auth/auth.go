package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	secretKey     []byte
	tokenLifeTime time.Duration
)

type Claims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func Init(secret []byte, lifeTime time.Duration) {
	secretKey = secret
	tokenLifeTime = lifeTime
}

func Sign(userId string) (string, error) {
	fmt.Println(secretKey, tokenLifeTime)
	claims := Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenLifeTime).Unix(),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return tokenClaims.SignedString([]byte("secret"))
}

package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
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
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte("secret"))
}

func Verify(token string) (userId string, err error) {
	token = strings.Split(token, "Bearer ")[1]
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// detail for jwt token err message
	if err != nil {
		var message string
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				message = "token is malformed"
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				message = "token could not be verified because of signing problems"
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				message = "signature validation failed"
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				message = "token is expired"
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				message = "token is not yet valid before sometime"
			} else {
				message = "can not handle this token"
			}
		}
		err = errors.New(message)
		return "", err
	}

	claims, _ := tokenClaims.Claims.(*Claims)
	if claims.UserId == "" {
		return "", errors.New("token is improper")
	}
	return claims.UserId, nil
}

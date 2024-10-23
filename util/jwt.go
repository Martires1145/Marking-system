package util

import (
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"time"
)

var jwtKey = []byte(viper.GetString("jwt.token"))

type Claims struct {
	UID  uint
	Role int
	jwt.StandardClaims
}

func CreatToken(uid uint, role int) (string, error) {
	expirationTime := time.Now().Add(31 * 24 * time.Hour)
	claims := &Claims{
		Role: role,
		UID:  uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    viper.GetString("jwt.issuer"),
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i any, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}

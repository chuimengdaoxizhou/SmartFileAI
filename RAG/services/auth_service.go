package services

import (
	"RAG/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(userID string) (string, error) {
	expireHours := utils.Config.GetInt("jwt.expire_hours")
	secret := utils.Config.GetString("jwt.secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * time.Duration(expireHours)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (string, error) {
	secret := utils.Config.GetString("jwt.secret")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["userID"].(string), nil
	}
	return "", jwt.ErrSignatureInvalid
}

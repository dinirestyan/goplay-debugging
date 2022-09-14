package utils

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetClaims(ctx *gin.Context) (string, error) {
	bearerToken := ctx.Request.Header.Get("Authorization")
	tokenString := strings.ReplaceAll(bearerToken, "Bearer ", "")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		hmacSecretString := "secret"
		hmacSecret := []byte(hmacSecretString)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_id"].(string), nil
	}
	return fmt.Sprintf("%x", "00000000-0000-0000-0000-000000000000"), err
}

func IsPasswordValid(hashedPassword, loginPassword string) bool {
	hashPass := []byte(hashedPassword)
	pass := []byte(loginPassword)

	if err := bcrypt.CompareHashAndPassword(hashPass, pass); err != nil {
		return false
	}
	return true
}

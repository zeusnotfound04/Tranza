package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))


func GenerateJWT(userId , email , username string) (string , error) {
	claims := jwt.MapClaims{
		"user_id" : userId,
		"email" : email ,
		"username" : username,
		"exp" : time.Now().Add(48 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256 , claims)
	return token.SignedString(jwtSecret)
} 
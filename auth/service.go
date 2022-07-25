package auth

import (
	"log"

	"github.com/golang-jwt/jwt/v4"
)

type Service interface {
	GenerateToken(userID int) (string, error)
}

type jwtService struct{}

func NewService() *jwtService {
	return &jwtService{}
}

var secretKey = []byte("#ThinkB4YouType!")

func (s *jwtService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		log.Panic("Error while generating JWT. Error: ", err.Error())
		return "", err
	}
	return signedToken, nil
}

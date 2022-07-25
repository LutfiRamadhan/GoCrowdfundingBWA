package auth

import (
	"errors"
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

func (s *jwtService) ValidateToken(userToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid Token")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return token, err
	}
	return token, nil
}

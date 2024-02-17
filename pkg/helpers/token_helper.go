package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct{
	Id primitive.ObjectID `json:"id"`
	jwt.StandardClaims
	UserType string `json:"userType"`
}

func CreateToken(id primitive.ObjectID, userType string) (tokenString string, err error) {
	claims := &Claims{
		Id:       id,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))); err != nil {
		return "", err
	} else {
		return signedToken, nil
	}
}

func VerifyToken(token string) (string, string, error) {
	if token == "" {
		return "", "", errors.New("token is empty")
	}

	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "", errors.New("signature is invalid")
		}
		return "", "", errors.New("token is invalid")
	}
	// check the type of parsed token
	if !parsedToken.Valid {
		return "", "", errors.New("parsed token is invalid")
	}
	if claims.UserType == ""{
		return "", "", errors.New("token claims are nil")
	}
	return claims.Id.Hex(), claims.UserType, nil
}

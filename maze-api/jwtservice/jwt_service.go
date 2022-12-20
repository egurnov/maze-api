package jwtservice

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/egurnov/maze-api/maze-api/model"
)

var _ model.JWTService = &JWTService{}

type JWTService struct {
	secretKey []byte
	name      string
	ttl       time.Duration
}

type Claims struct {
	model.CustomClaims
	jwt.StandardClaims
}

func New(secretKey []byte, ttl time.Duration) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		ttl:       ttl,
	}
}

func (jwtSrc *JWTService) GenerateToken(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		CustomClaims: model.CustomClaims{
			UserID: id,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtSrc.ttl).Unix(),
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    jwtSrc.name,
			Audience:  jwtSrc.name,
			Subject:   strconv.Itoa(int(id)),
		},
	})

	t, err := token.SignedString(jwtSrc.secretKey)
	if err != nil {
		return "", err
	}
	return t, nil
}

func (jwtSrc *JWTService) ValidateToken(tokenString string) (*model.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// TODO: (ae) Switch to using RSA
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSrc.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("bad claims")
	}
	if err := claims.Valid(); err != nil {
		return nil, err
	}

	return &claims.CustomClaims, nil
}

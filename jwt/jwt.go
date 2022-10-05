package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	key        []byte
	expireTime time.Duration
}

type Claim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (j *JWT) Generate(username string) (string, error) {
	claim := Claim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.expireTime).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) Validate(t string) (*Claim, error) {
	claim := new(Claim)
	token, err := jwt.ParseWithClaims(
		t,
		claim,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.key), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claim)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("expired token")
	}
	return claim, nil
}

func New(key string, expireTime time.Duration) *JWT {
	return &JWT{
		key:        []byte(key),
		expireTime: expireTime,
	}
}

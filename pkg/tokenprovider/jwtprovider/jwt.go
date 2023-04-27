package jwtprovider

import (
	"restapi/pkg/tokenprovider"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtProvider struct {
	secret string
}

func New(secret string) *jwtProvider {
	return &jwtProvider{secret: secret}
}

type myClaims struct {
	Payload tokenprovider.TokenPayload `json:"payload"`
	jwt.RegisteredClaims
}

func (j *jwtProvider) Generate(payload tokenprovider.TokenPayload, expiryInSeconds int64) (*tokenprovider.Token, error) {
	// generate the JWT
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		payload,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Second * time.Duration(expiryInSeconds))),
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),
		},
	})

	myToken, err := t.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}

	// return the token
	return &tokenprovider.Token{
		Token:   myToken,
		Expiry:  expiryInSeconds,
		Created: time.Now(),
	}, nil
}

func (j *jwtProvider) Validate(myToken string) (*tokenprovider.TokenPayload, error) {
	res, err := jwt.ParseWithClaims(myToken, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, tokenprovider.ErrInvalidToken
	}

	// validate the token
	if !res.Valid {
		return nil, tokenprovider.ErrInvalidToken
	}

	claims, ok := res.Claims.(*myClaims)
	if !ok {
		return nil, tokenprovider.ErrInvalidToken
	}

	// return the token
	return &claims.Payload, nil
}

func (j *jwtProvider) String() string {
	return "JWT implement Provider"
}

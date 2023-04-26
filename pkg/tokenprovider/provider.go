package tokenprovider

import (
	"errors"
	"time"
)

var (
	ErrInvalidToken = errors.New("token invalid")
)

type Token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int64     `json:"expiry"`
}

type TokenPayload struct {
	UserId int64 `json:"user_id"`
}

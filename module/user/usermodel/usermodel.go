package usermodel

import (
	"errors"
	"strings"
)

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r RegisterReq) Validate() error {
	if len(strings.Trim(r.Username, " ")) == 0 {
		return errors.New("username is empty")
	}

	if len(strings.Trim(r.Password, " ")) == 0 {
		return errors.New("password is empty")
	}

	return nil
}

type RegisterRes struct {
}

type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r LoginReq) Validate() error {
	if len(strings.Trim(r.Username, " ")) == 0 {
		return errors.New("username is empty")
	}

	if len(strings.Trim(r.Password, " ")) == 0 {
		return errors.New("password is empty")
	}

	return nil
}

type LoginRes struct {
	Token string `json:"token"`
}

type User struct {
	Id             int64
	Username       string
	HashedPassword string
	Salt           string
}

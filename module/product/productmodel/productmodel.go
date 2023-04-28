package productmodel

import (
	"errors"
	"strings"
)

type CreateReq struct {
	Name string `json:"name"`
}

func (r CreateReq) Validate() error {
	if len(strings.Trim(r.Name, " ")) == 0 {
		return errors.New("name is empty")
	}

	return nil
}

type CreateRes struct {
}

type ListReq struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListRes struct {
	Products []*Product `json:"products"`
}

type UpdateReq struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (r UpdateReq) Validate() error {
	if len(strings.Trim(r.Name, " ")) == 0 {
		return errors.New("name is empty")
	}

	return nil
}

type UpdateRes struct {
}

type DetailReq struct {
	Id int64 `json:"id"`
}

type DetailRes struct {
	Product *Product `json:"product"`
}

type DeleteReq struct {
	Id int64 `json:"id"`
}

type DeleteRes struct {
}

type Product struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

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

type Product struct {
	Id   int64
	Name string
}

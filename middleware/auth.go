package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"restapi/pkg/appconst"
	"restapi/pkg/apperr"
	"restapi/pkg/response"
	"restapi/pkg/tokenprovider"
	"strings"
)

type TokenProvider interface {
	Validate(token string) (*tokenprovider.TokenPayload, error)
}

func Auth(tokenProvider TokenProvider) func(c *gin.Context) {
	return func(c *gin.Context) {
		token, err := extractTokenFromHeaderString(c.GetHeader("Authorization"))
		if err != nil {
			response.Error(c, apperr.Wrap(err, appconst.CodeUnauthorized, "", http.StatusUnauthorized))
			return
		}

		_, err = tokenProvider.Validate(token)
		if err != nil {
			response.Error(c, apperr.Wrap(err, appconst.CodeUnauthorized, "", http.StatusUnauthorized))
			return
		}

		c.Next()
	}
}

func extractTokenFromHeaderString(s string) (string, error) {
	//"Authorization" : "Bearer {token}"
	parts := strings.Split(s, " ")
	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("wrong header format")
	}

	return parts[1], nil
}

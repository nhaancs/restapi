package userbusiness

import (
	"github.com/stretchr/testify/suite"
	uBMock "restapi/module/user/userbusiness/mocks"
	"testing"
)

type businessSuite struct {
	suite.Suite
	userStore     *uBMock.UserStore
	tokenProvider *uBMock.TokenProvider
	hasher        *uBMock.Hasher
	userBusiness  *business
}

func (s *businessSuite) SetupTest() {
	// Init suite
	s.userStore = new(uBMock.UserStore)
	s.tokenProvider = new(uBMock.TokenProvider)
	s.hasher = new(uBMock.Hasher)
	s.userBusiness = New(s.userStore, s.tokenProvider, s.hasher)
}

func TestBusinessSuite(t *testing.T) {
	suite.Run(t, new(businessSuite))
}

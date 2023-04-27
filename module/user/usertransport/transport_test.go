package usertransport

import (
	"testing"

	"github.com/stretchr/testify/suite"
	uTMock "restapi/module/user/usertransport/mocks"
)

type transportSuite struct {
	suite.Suite
	userBusiness  *uTMock.UserBusiness
	userTransport *transport
}

func (s *transportSuite) SetupTest() {
	// Init suite
	s.userBusiness = new(uTMock.UserBusiness)
	s.userTransport = New(s.userBusiness)
}

func TestTransportSuite(t *testing.T) {
	suite.Run(t, new(transportSuite))
}

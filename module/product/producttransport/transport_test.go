package producttransport

import (
	"testing"

	"github.com/stretchr/testify/suite"
	pTMock "restapi/module/product/producttransport/mocks"
)

type transportSuite struct {
	suite.Suite
	tokenProvider    *pTMock.TokenProvider
	productBusiness  *pTMock.ProductBusiness
	productTransport *transport
}

func (s *transportSuite) SetupTest() {
	// Init suite
	s.productBusiness = new(pTMock.ProductBusiness)
	s.tokenProvider = new(pTMock.TokenProvider)
	s.productTransport = New(s.productBusiness, s.tokenProvider)
}

func TestTransportSuite(t *testing.T) {
	suite.Run(t, new(transportSuite))
}

package producttransport

import (
	"testing"

	"github.com/stretchr/testify/suite"
	pTMock "restapi/module/product/producttransport/mocks"
)

type transportSuite struct {
	suite.Suite
	productBusiness  *pTMock.ProductBusiness
	productTransport *transport
}

func (s *transportSuite) SetupTest() {
	// Init suite
	s.productBusiness = new(pTMock.ProductBusiness)
	s.productTransport = New(s.productBusiness)
}

func TestTransportSuite(t *testing.T) {
	suite.Run(t, new(transportSuite))
}

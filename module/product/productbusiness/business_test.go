package productbusiness

import (
	"github.com/stretchr/testify/suite"
	pBMock "restapi/module/product/productbusiness/mocks"
	"testing"
)

type businessSuite struct {
	suite.Suite
	productStore    *pBMock.ProductStore
	productBusiness *business
}

func (s *businessSuite) SetupTest() {
	// Init suite
	s.productStore = new(pBMock.ProductStore)
	s.productBusiness = New(s.productStore)
}

func TestBusinessSuite(t *testing.T) {
	suite.Run(t, new(businessSuite))
}

package productstore

import (
	"context"
	"errors"
	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"restapi/module/product/productmodel"
	"testing"
)

type storeSuite struct {
	suite.Suite
	store  *store
	dbMock sqlMock.Sqlmock
}

func (s *storeSuite) SetupTest() {
	t := s.T()
	db, dbMock := NewDBMock(t)
	{
		s.dbMock = dbMock
		s.store = New(db)
	}
}

func (s *storeSuite) TearDownTestSuite() {
	s.store.db.Close()
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(storeSuite))
}

func NewDBMock(t *testing.T) (*sqlx.DB, sqlMock.Sqlmock) {
	matcherOpt := sqlMock.QueryMatcherOption(sqlMock.QueryMatcherRegexp)
	db, mock, err := sqlMock.New(matcherOpt)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return sqlx.NewDb(db, "sqlmock"), mock
}

// ================================================================================================
func (s *storeSuite) TestStoreSuite_Insert() {
	cases := []struct {
		name    string
		product productmodel.Product
		err     error
		mock    func()
	}{
		{
			name:    "success",
			product: productmodel.Product{Name: "product name"},
			err:     nil,
			mock: func() {
				s.dbMock.ExpectExec(`INSERT INTO products`).
					WillReturnResult(sqlMock.NewResult(1, 1))
			},
		},
		{
			name:    "database error",
			product: productmodel.Product{Name: "product name"},
			err:     errors.New("dummy error"),
			mock: func() {
				s.dbMock.ExpectExec(`INSERT INTO products`).
					WillReturnError(errors.New("dummy error"))
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			err := s.store.Insert(ctx, &c.product)
			assert.Equal(t, c.err, err)
		})
	}
}

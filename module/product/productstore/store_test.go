package productstore

import (
	"context"
	"database/sql"
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

func (s *storeSuite) TestStoreSuite_Update() {
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
				s.dbMock.ExpectExec(`UPDATE products SET`).
					WillReturnResult(sqlMock.NewResult(1, 1))
			},
		},
		{
			name:    "database error",
			product: productmodel.Product{Name: "product name"},
			err:     errors.New("dummy error"),
			mock: func() {
				s.dbMock.ExpectExec(`UPDATE products SET`).
					WillReturnError(errors.New("dummy error"))
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			err := s.store.Update(ctx, &c.product)
			assert.Equal(t, c.err, err)
		})
	}
}

func (s *storeSuite) TestStoreSuite_FindById() {
	cases := []struct {
		name    string
		id      int64
		err     error
		product *productmodel.Product
		mock    func()
	}{
		{
			name:    "success",
			id:      1,
			err:     nil,
			product: &productmodel.Product{Id: 1, Name: "product name"},
			mock: func() {
				s.dbMock.ExpectQuery(`SELECT id, name FROM products WHERE id =`).
					WillReturnRows(sqlMock.NewRows([]string{"id", "name"}).
						AddRow(1, "product name"))
			},
		},
		{
			name:    "not found",
			id:      1,
			err:     ErrNotFound,
			product: nil,
			mock: func() {
				s.dbMock.ExpectQuery(`SELECT id, name FROM products WHERE id =`).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:    "internal error",
			id:      1,
			err:     errors.New("dummy error"),
			product: nil,
			mock: func() {
				s.dbMock.ExpectQuery(`SELECT id, name FROM products WHERE id =`).
					WillReturnError(errors.New("dummy error"))
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			product, err := s.store.FindById(ctx, c.id)
			assert.Equal(t, c.err, err)
			assert.Equal(t, c.product, product)
		})
	}
}

func (s *storeSuite) TestStoreSuite_Delete() {
	cases := []struct {
		name string
		id   int64
		err  error
		mock func()
	}{
		{
			name: "success",
			id:   1,
			err:  nil,
			mock: func() {
				s.dbMock.ExpectExec(`DELETE FROM products WHERE id =`).
					WillReturnResult(sqlMock.NewResult(1, 1))
			},
		},
		{
			name: "internal error",
			id:   1,
			err:  errors.New("dummy error"),
			mock: func() {
				s.dbMock.ExpectExec(`DELETE FROM products WHERE id =`).
					WillReturnError(errors.New("dummy error"))
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			err := s.store.Delete(ctx, c.id)
			assert.Equal(t, c.err, err)
		})
	}
}

func (s *storeSuite) TestStoreSuite_Select() {
	cases := []struct {
		name     string
		err      error
		products []*productmodel.Product
		mock     func()
	}{
		{
			name:     "success",
			products: []*productmodel.Product{{Id: 1, Name: "product name"}},
			err:      nil,
			mock: func() {
				s.dbMock.ExpectQuery(`SELECT id, name FROM products`).
					WillReturnRows(sqlMock.NewRows([]string{"id", "name"}).
						AddRow(1, "product name"))
			},
		},
		{
			name:     "internal error",
			products: nil,
			err:      errors.New("dummy error"),
			mock: func() {
				s.dbMock.ExpectQuery(`SELECT id, name FROM products`).
					WillReturnError(errors.New("dummy error"))
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			products, err := s.store.Select(ctx, 0, 10)
			assert.Equal(t, c.err, err)
			assert.Equal(t, c.products, products)
		})
	}
}

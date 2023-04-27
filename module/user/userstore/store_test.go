package userstore

import (
	"context"
	"database/sql"
	"errors"
	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"restapi/module/user/usermodel"
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

func (s *storeSuite) TestStoreSuite_FindByUsername() {
	cases := []struct {
		name     string
		username string
		err      error
		user     *usermodel.User
		mock     func()
	}{
		{
			name:     "success",
			username: "username",
			err:      nil,
			user:     &usermodel.User{Id: 1, Username: "user", HashedPassword: "pass", Salt: "salt"},
			mock: func() {
				s.dbMock.ExpectQuery(`SELECT username, hashed_password, salt FROM users WHERE username =`).
					WillReturnRows(sqlMock.NewRows([]string{"id", "username", "hashed_password", "salt"}).
						AddRow(1, "user", "pass", "salt"))
			},
		},
		{
			name:     "not found",
			username: "username",
			err:      ErrNotFound,
			user:     nil,
			mock: func() {
				s.dbMock.ExpectQuery(`SELECT username, hashed_password, salt FROM users WHERE username =`).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:     "internal error",
			username: "username",
			err:      errors.New("dummy error"),
			user:     nil,
			mock: func() {
				s.dbMock.ExpectQuery(`SELECT username, hashed_password, salt FROM users WHERE username =`).
					WillReturnError(errors.New("dummy error"))
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			user, err := s.store.FindByUsername(ctx, c.username)
			assert.Equal(t, c.err, err)
			assert.Equal(t, c.user, user)
		})
	}
}

func (s *storeSuite) TestStoreSuite_Insert() {
	cases := []struct {
		name string
		user usermodel.User
		err  error
		mock func()
	}{
		{
			name: "success",
			user: usermodel.User{Username: "user", HashedPassword: "pass", Salt: "salt"},
			err:  nil,
			mock: func() {
				s.dbMock.ExpectExec(`INSERT INTO users`).
					WillReturnResult(sqlMock.NewResult(1, 1))
			},
		},
		{
			name: "database error",
			user: usermodel.User{Username: "user", HashedPassword: "pass", Salt: "salt"},
			err:  errors.New("dummy error"),
			mock: func() {
				s.dbMock.ExpectExec(`INSERT INTO users`).
					WillReturnError(errors.New("dummy error"))
			},
		},
	}

	t := s.T()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			ctx := context.Background()
			err := s.store.Insert(ctx, &c.user)
			assert.Equal(t, c.err, err)
		})
	}
}

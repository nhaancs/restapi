package userstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"restapi/module/user/usermodel"
)

var ErrNotFound = errors.New("user not found")

type userRow struct {
	Id             int64  `db:"id,primarykey"`
	Username       string `db:"username"`
	HashedPassword string `db:"hashed_password"`
	Salt           string `db:"salt"`
}

func toUserRow(m *usermodel.User) *userRow {
	return &userRow{
		Id:             m.Id,
		Username:       m.Username,
		HashedPassword: m.HashedPassword,
		Salt:           m.Salt,
	}
}

func (r *userRow) toUserModel() *usermodel.User {
	return &usermodel.User{
		Id:             r.Id,
		Username:       r.Username,
		HashedPassword: r.HashedPassword,
		Salt:           r.Salt,
	}
}

type store struct {
	db  *sqlx.DB
	tbl string
}

func New(db *sqlx.DB) *store {
	return &store{
		db:  db,
		tbl: "users",
	}
}

func (s *store) FindByUsername(ctx context.Context, username string) (*usermodel.User, error) {
	cmd := fmt.Sprintf(`SELECT username, hashed_password, salt FROM %s WHERE username = ?`, s.tbl)
	var row userRow
	if err := s.db.GetContext(ctx, &row, cmd, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return row.toUserModel(), nil
}

func (s *store) Insert(ctx context.Context, user *usermodel.User) error {
	cmd := fmt.Sprintf(`INSERT INTO %s (username, hashed_password, salt) VALUES (:username, :hashed_password, :salt)`, s.tbl)
	row := toUserRow(user)
	if _, err := s.db.NamedQueryContext(ctx, cmd, row); err != nil {
		return err
	}
	return nil
}

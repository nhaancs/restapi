package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Addr, Name, User, Pass string
}

func Connect(c Config) (*sqlx.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8", c.User, c.Pass, c.Addr, c.Name)
	db, err := sqlx.Connect("mysql", connString)
	if err != nil {
		return nil, fmt.Errorf("cannot connect database: %v", err)
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(4096)

	return db, nil
}

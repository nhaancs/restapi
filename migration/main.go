package main

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"restapi/pkg/config"
	"restapi/pkg/database"
)

type appConfig struct {
	DB struct {
		Addr string
		Name string
		User string
		Pass string
	}
}

func main() {
	var cfg appConfig
	if err := config.Load("config.yaml", &cfg); err != nil {
		log.Fatalln("error loading config: ", err)
	}

	db, err := database.Connect(database.Config{
		Addr: cfg.DB.Addr,
		Name: cfg.DB.Name,
		User: cfg.DB.User,
		Pass: cfg.DB.Pass,
	})
	if err != nil {
		log.Fatalln("failed to connect database: ", err)
	}
	defer db.Close()

	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		log.Fatalln("failed create driver: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration/migrations", cfg.DB.Name, driver)
	if err != nil {
		log.Fatalln("failed get migrate instance: ", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("migrate completed with no change")
			return
		}
		log.Fatalln("failed to migrate: ", err)
	}

	log.Println("migrate completed")
}

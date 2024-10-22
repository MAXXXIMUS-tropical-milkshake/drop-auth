package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath string

	flag.StringVar(&storagePath, "db_url", "", "url to db")
	flag.StringVar(&migrationsPath, "migrations_path", "", "path to migrations")
	flag.Parse()

	if storagePath == "" {
		panic("db_url is required")
	}

	if migrationsPath == "" {
		panic("migrations_path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		storagePath+"?sslmode=disable",
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations applied successfully")
}

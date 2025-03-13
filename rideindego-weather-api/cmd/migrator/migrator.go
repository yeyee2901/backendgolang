package main

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/utils"
)

func main() {
	migrationSource := "file://sql/migrations"
	migrationTarget := utils.BuildDatasourceName(utils.DataSource{
		User:     os.Getenv("POSTGRESQL_USERNAME"),
		Password: os.Getenv("POSTGRESQL_PASSWORD"),
		Host:     os.Getenv("POSTGRESQL_HOST"),
		Database: os.Getenv("POSTGRESQL_DATABASE"),
	})
	fmt.Println(migrationTarget)

	db, err := sqlx.Connect("postgres", migrationTarget)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	mig, err := migrate.NewWithDatabaseInstance(
		migrationSource,
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}

	mig.Up()

	// standby
}

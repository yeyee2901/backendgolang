package openweather

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/utils"
)

func TestFetchData(t *testing.T) {
	err := godotenv.Load("../../build.env")
	if err != nil {
		t.Fatal(err)
	}
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	baseURL := os.Getenv("OPENWEATHER_URL")
	db, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := cleanDatabase(db)
		if err != nil {
			t.Log("[WARNING] failed to clean database")
		}
	})

	weather := NewOpenWeather(apiKey, baseURL, db)
	resp, err := weather.fetchData(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// just check for the array availability
	if len(resp.Weather) == 0 {
		t.Log("weather:", resp)
		t.Fatal("empty weather")
	}
}

func TestRefreshData(t *testing.T) {
	err := godotenv.Load("../../build.env")
	if err != nil {
		t.Fatal(err)
	}
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	baseURL := os.Getenv("OPENWEATHER_URL")
	db, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := cleanDatabase(db)
		if err != nil {
			t.Log("[WARNING] failed to clean database")
		}
	})

	weather := NewOpenWeather(apiKey, baseURL, db)
	err = weather.RefreshData(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func connectDB() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", utils.BuildDatasourceName(utils.DataSource{
		User:     os.Getenv("POSTGRESQL_USERNAME"),
		Password: os.Getenv("POSTGRESQL_PASSWORD"),
		Host:     os.Getenv("POSTGRESQL_HOST"),
		Database: os.Getenv("POSTGRESQL_DATABASE"),
	}))
}

func cleanDatabase(conn *sqlx.DB) error {
	queries := []string{
		` DELETE FROM openweather_master `,
		` DELETE FROM openweather_details `,
	}

	tx, err := conn.Beginx()
	if err != nil {
		return err
	}

	for _, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

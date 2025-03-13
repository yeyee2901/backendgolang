package rideindego

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/utils"
)

func TestFetchData(t *testing.T) {
	err := godotenv.Load("../../build.env")
	if err != nil {
		t.Fatal(err)
	}

	dbConn, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	ride := NewRideIndeGoService(os.Getenv("RIDE_INDEGO_URL"), dbConn)
	resp, err := ride.fetchData(context.Background())
	if err != nil {
		t.Fatal("failed to fetch:", err)
	}

	if resp == nil {
		t.Fatal("resp should not be nil")
	}

	if len(resp.Features) == 0 {
		t.Fatal("resp.features should not be empty")
	}

	if resp.LastUpdated.IsZero() {
		t.Fatal("resp.lastUpdated should not be zero value")
	}

	if resp.Type == "" {
		t.Fatal("resp.type should not be empty")
	}

	t.Log(resp.LastUpdated, resp.Type)
	t.Logf("response: %+v", resp.Features[1])
}

func TestRefreshData(t *testing.T) {
	err := godotenv.Load("../../build.env")
	if err != nil {
		t.Fatal(err)
	}

	dbConn, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := cleanDatabase(dbConn)
		if err != nil {
			t.Log("[WARNING] failed to clean database:", err)
		}
	})

	ride := NewRideIndeGoService(os.Getenv("RIDE_INDEGO_URL"), dbConn)
	freshData, err := ride.RefreshData(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Log("last updated: ", freshData.LastUpdated)
	t.Log("type: ", freshData.Type)
	t.Logf("features length %d, idx [0]: %+v", len(freshData.Features), freshData.Features[0])
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
		` DELETE FROM rideindego_master `,
		` DELETE FROM rideindego_features `,
		` DELETE FROM rideindego_properties `,
		` DELETE FROM rideindego_bikes `,
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

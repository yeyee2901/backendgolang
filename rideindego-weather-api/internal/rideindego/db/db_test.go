package db

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestSaveRideIndegoMaster(t *testing.T) {
	var (
		testDataFetchID     = uuid.NewString()
		testDataLastUpdated = time.Now()
		testDataFeatureType = "test"
	)
	conn, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	cleanDatabase(conn)

	db := &dbRideIndego{conn}
	tx, err := conn.Beginx()
	if err != nil {
		t.Fatal("failed to begin tx:", err)
	}

	err = db.saveTableMaster(tx, &TableRideIndegoMaster{
		FetchID:     testDataFetchID,
		LastUpdated: testDataLastUpdated,
		FeatureType: testDataFeatureType,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	// check the data
	q := `SELECT * FROM rideindego_master WHERE fetch_id = $1 AND feature_type = $2`
	res := new(TableRideIndegoMaster)
	err = conn.Get(res, q, testDataFetchID, testDataFeatureType)
	if err != nil {
		t.Fatal("could not load test data:", err)
	}
}

func connectDB() (*sqlx.DB, error) {
	err := godotenv.Load("../../../build.env")
	if err != nil {
		return nil, err
	}

	datasource := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=UTC",
		os.Getenv("POSTGRESQL_USERNAME"),
		os.Getenv("POSTGRESQL_PASSWORD"),
		os.Getenv("POSTGRESQL_HOST"),
		os.Getenv("POSTGRESQL_PORT"),
		os.Getenv("POSTGRESQL_DATABASE"),
	)

	return sqlx.Connect("postgres", datasource)
}

func cleanDatabase(conn *sqlx.DB) error {
	q1 := `
        DELETE FROM rideindego_master
    `

	tx, err := conn.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(q1)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

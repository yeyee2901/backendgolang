package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/utils"
)

func TestStoreData(t *testing.T) {
	var (
		testFetchID = uuid.NewString()
		testMaster  = &TableWeatherMaster{
			FetchID:       testFetchID,
			Coord:         fmt.Sprintf("POINT(%f %f)", 0.123, 0.345),
			Base:          "test",
			MainTemp:      100,
			MainFeelsLike: 100,
			MainTempMin:   100,
			MainTempMax:   100,
			MainPressure:  100,
			MainHumidity:  100,
			MainSeaLevel:  100,
			MainGrndLevel: 100,
			Visibility:    100,
			WindSpeed:     100,
			WindDeg:       100,
			WindGust:      100,
			CloudsAll:     100,
			Dt:            100,
			SysType:       100,
			SysID:         100,
			SysCountry:    "test",
			SysSunrise:    100,
			SysSunset:     100,
		}
		testDetails = []*TableWeatherDetails{
			{
				FetchID:     testFetchID,
				IDx:         100,
				ID:          100,
				Main:        "test",
				Description: "test",
				Icon:        "test",
			},
		}
	)

	err := godotenv.Load("../../../build.env")
	if err != nil {
		t.Fatal(err)
	}

	db, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := cleanDatabase(db)
		if err != nil {
			t.Log("[WARNING] failed to clean DB")
		}
	})

	weather := NewDBWeather(db)
	err = weather.StoreDatas(context.Background(), testMaster, testDetails)
	if err != nil {
		t.Fatal(err)
	}

	// check for sample data
	q := `SELECT * FROM openweather_master WHERE fetch_id = $1 LIMIT 1`
	res := new(TableWeatherMaster)
	err = db.Get(res, q, testFetchID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("openweather_master: ", res)
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

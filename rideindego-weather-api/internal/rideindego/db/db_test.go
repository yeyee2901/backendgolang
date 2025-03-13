package db

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/utils"

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

	err = cleanDatabase(conn)
	if err != nil {
		t.Fatal("failed to clean database:", err)
	}

	db := &dbRideIndego{conn}
	tx, err := conn.Beginx()
	if err != nil {
		t.Fatal("failed to begin tx:", err)
	}

	err = db.saveMaster(tx, &TableRideIndegoMaster{
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

func TestSaveFeatures(t *testing.T) {
	var (
		testDataFetchID     = uuid.NewString()
		testDataFeatureType = "test"
	)
	conn, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	err = cleanDatabase(conn)
	if err != nil {
		t.Fatal("failed to clean database:", err)
	}

	db := &dbRideIndego{conn}
	tx, err := conn.Beginx()
	if err != nil {
		t.Fatal("failed to begin tx:", err)
	}

	err = db.bulkSaveFeatures(tx, []*TableRideIndegoFeatures{
		{
			FetchID:     testDataFetchID,
			FeatureID:   1,
			FeatureType: testDataFeatureType,
			GeoType:     "test",
			GeoCoord:    fmt.Sprintf("POINT(%f %f)", 0.123, 0.345),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func TestStoreDatas(t *testing.T) {
	var (
		randUUID = uuid.NewString()
		now      = time.Now()

		testMaster = &TableRideIndegoMaster{
			FetchID:     randUUID,
			LastUpdated: now,
			FeatureType: "test",
		}
		testFeatures = []*TableRideIndegoFeatures{
			{
				FetchID:     randUUID,
				FeatureID:   1,
				FeatureType: "test",
				GeoType:     "test",
				GeoCoord:    fmt.Sprintf("POINT(%f %f)", 0.123, 0.345),
			},
		}
		testProperties = []*TableRideIndegoProperties{
			{
				FetchID:               randUUID,
				FeatureID:             1,
				PropertiesID:          1,
				Coordinates:           fmt.Sprintf("POINT(%f %f)", 0.123, 0.345),
				Name:                  "",
				TotalDocks:            0,
				DocksAvailable:        0,
				BikesAvailable:        0,
				ClassicBikesAvailable: 0,
				SmartBikesAvailable:   0,
				EletricBikesAvailable: 0,
				RewardBikesAvailable:  0,
				RewardDocksAvailable:  0,
				KioskStatus:           "",
				KioskPublicStatus:     "",
				KioskConnectionStatus: "",
				KioskType:             0,
				AddressStreet:         "",
				AddressCity:           "",
				AddressState:          "",
				AddressZipCode:        "",
				CloseTime:             "",
				EventEnd:              "",
				EventStart:            "",
				IsEventBased:          false,
				IsVirtual:             false,
				KioskID:               0,
				Notes:                 "",
				OpenTime:              "",
				PublicText:            "",
				Timezone:              "",
				TrikesAvailable:       0,
			},
		}
	)

	conn, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	err = cleanDatabase(conn)
	if err != nil {
		t.Fatal("failed to clean database:", err)
	}

	db := &dbRideIndego{conn}
	err = db.storeDatas(testMaster, testFeatures, testProperties)
	if err != nil {
		t.Fatal(err)
	}
}

func connectDB() (*sqlx.DB, error) {
	err := godotenv.Load("../../../build.env")
	if err != nil {
		return nil, err
	}

	return sqlx.Connect("postgres", utils.BuildDatasourceName(utils.DataSource{
		User:     os.Getenv("POSTGRESQL_USERNAME"),
		Password: os.Getenv("POSTGRESQL_PASSWORD"),
		Host:     os.Getenv("POSTGRESQL_HOST"),
		Database: os.Getenv("POSTGRESQL_DATABASE"),
	}))
}

func cleanDatabase(conn *sqlx.DB) error {
	q1 := ` DELETE FROM rideindego_master `

	q2 := ` DELETE FROM rideindego_features `

	q3 := ` DELETE FROM rideindego_properties `

	tx, err := conn.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(q1)
	if err != nil {
		return err
	}

	_, err = tx.Exec(q2)
	if err != nil {
		return err
	}

	_, err = tx.Exec(q3)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

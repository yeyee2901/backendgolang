package db

import (
	"context"
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
		testDataDataType    = "test"
	)
	conn, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := cleanDatabase(conn)
		if err != nil {
			t.Log("[WARNING] failed to clean database:", err)
		}
	})

	db := &dbRideIndego{conn}
	tx, err := conn.Beginx()
	if err != nil {
		t.Fatal("failed to begin tx:", err)
	}

	err = db.saveMaster(tx, &TableRideIndegoMaster{
		FetchID:     testDataFetchID,
		LastUpdated: testDataLastUpdated,
		DataType:    testDataDataType,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	// check the data
	q := `SELECT * FROM rideindego_master WHERE fetch_id = $1 AND data_type = $2`
	res := new(TableRideIndegoMaster)
	err = conn.Get(res, q, testDataFetchID, testDataDataType)
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

	t.Cleanup(func() {
		err := cleanDatabase(conn)
		if err != nil {
			t.Log("[WARNING] failed to clean database:", err)
		}
	})

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
			DataType:    "test",
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

		testBikes = []*TableRideIndegoBikes{
			{
				FetchID:      randUUID,
				FeatureID:    1,
				PropertiesID: 1,
				ID:           1,
				DockNumber:   1,
				IsElectric:   false,
				IsAvailable:  false,
				Battery:      30,
			},
		}
	)

	conn, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := cleanDatabase(conn)
		if err != nil {
			t.Log("[WARNING] failed to clean database:", err)
		}
	})

	db := &dbRideIndego{conn}
	err = db.StoreDatas(context.Background(), testMaster, testFeatures, testProperties, testBikes)
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

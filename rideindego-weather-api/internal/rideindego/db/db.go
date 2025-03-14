package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// DBRideIndegoProvider is the interface used to interact with
// ride indego data.
type DBRideIndegoProvider interface {
	// StoreDatas stores the data to the database
	StoreDatas(context.Context, *TableRideIndegoMaster, []*TableRideIndegoFeatures, []*TableRideIndegoProperties, []*TableRideIndegoBikes) error

	// SearchData searches the data based on parameter.
	// Time is mandatory, others are optional
	SearchData(c context.Context, at time.Time, kioskID string) (*SearchDataResult, error)
}

// DBRideIndego is used to interact with Ride Indego data
type dbRideIndego struct {
	conn *sqlx.DB
}

var ErrRideIndegoDB = fmt.Errorf("rideindego db error")

func NewDBRideIndego(db *sqlx.DB) DBRideIndegoProvider {
	return &dbRideIndego{
		conn: db,
	}
}

func (db *dbRideIndego) StoreDatas(
	ctx context.Context,
	master *TableRideIndegoMaster,
	features []*TableRideIndegoFeatures,
	properties []*TableRideIndegoProperties,
	bikes []*TableRideIndegoBikes,
) error {
	// begin transaction
	tx, err := db.conn.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Join(ErrRideIndegoDB, err)
	}

	// store master
	err = db.saveMaster(tx, master)
	if err != nil {
		tx.Rollback()
		return errors.Join(ErrRideIndegoDB, err)
	}

	// store features
	err = db.bulkSaveFeatures(tx, features)
	if err != nil {
		tx.Rollback()
		return errors.Join(ErrRideIndegoDB, err)
	}

	// store properties
	err = db.bulkSaveProperties(tx, properties)
	if err != nil {
		tx.Rollback()
		return errors.Join(ErrRideIndegoDB, err)
	}

	// store bikes
	err = db.bulkSaveBikes(tx, bikes)
	if err != nil {
		tx.Rollback()
		return errors.Join(ErrRideIndegoDB, err)
	}

	// commit
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Join(ErrRideIndegoDB, err)
	}

	return nil
}

func (db *dbRideIndego) saveMaster(tx *sqlx.Tx, param *TableRideIndegoMaster) error {
	q := `
        INSERT INTO rideindego_master
        (
            fetch_id,
            data_type,
            last_updated
        )
        VALUES
        (
            :fetch_id,
            :data_type,
            :last_updated
        )
    `

	_, err := tx.NamedExec(q, param)
	if err != nil {
		return fmt.Errorf("failed to insert rideindego master: %w", err)
	}

	return nil
}

func (db *dbRideIndego) bulkSaveFeatures(tx *sqlx.Tx, param []*TableRideIndegoFeatures) error {
	q := `
        INSERT INTO rideindego_features
        (
            fetch_id,
            feature_id,
            feature_type,
            geo_type,
            geo_coord
        )
        VALUES
        (
            :fetch_id,
            :feature_id,
            :feature_type,
            :geo_type,
            :geo_coord
        )
    `

	_, err := tx.NamedExec(q, param)
	if err != nil {
		return fmt.Errorf("failed to insert rideindego features: %w", err)
	}

	return nil
}

func (db *dbRideIndego) bulkSaveProperties(tx *sqlx.Tx, param []*TableRideIndegoProperties) error {
	q := `
        INSERT INTO rideindego_properties
        (
            fetch_id,
            feature_id,
            properties_id,
            coordinates,
            name,
            total_docks,
            docks_available,
            bikes_available,
            classic_bikes_available,
            smart_bikes_available,
            eletric_bikes_available,
            reward_bikes_available,
            reward_docks_available,
            kiosk_status,
            kiosk_public_status,
            kiosk_connection_status,
            kiosk_type,
            address_street,
            address_city,
            address_state,
            address_zip_code,
            close_time,
            event_end,
            event_start,
            is_event_based,
            is_virtual,
            kiosk_id,
            notes,
            open_time,
            public_text,
            timezone,
            trikes_available
        )
        VALUES 
        (
            :fetch_id,
            :feature_id,
            :properties_id,
            :coordinates,
            :name,
            :total_docks,
            :docks_available,
            :bikes_available,
            :classic_bikes_available,
            :smart_bikes_available,
            :eletric_bikes_available,
            :reward_bikes_available,
            :reward_docks_available,
            :kiosk_status,
            :kiosk_public_status,
            :kiosk_connection_status,
            :kiosk_type,
            :address_street,
            :address_city,
            :address_state,
            :address_zip_code,
            :close_time,
            :event_end,
            :event_start,
            :is_event_based,
            :is_virtual,
            :kiosk_id,
            :notes,
            :open_time,
            :public_text,
            :timezone,
            :trikes_available
        )
    `

	_, err := tx.NamedExec(q, param)
	if err != nil {
		return fmt.Errorf("failed to insert rideindego properties: %w", err)
	}

	return nil
}

func (db *dbRideIndego) bulkSaveBikes(tx *sqlx.Tx, param []*TableRideIndegoBikes) error {
	q := `
        INSERT INTO rideindego_bikes
        (
            fetch_id,
            feature_id,
            properties_id,
            id,
            dock_number,
            is_electric,
            is_available,
            battery
        )
        VALUES
        (
            :fetch_id,
            :feature_id,
            :properties_id,
            :id,
            :dock_number,
            :is_electric,
            :is_available,
            :battery
        )
    `

	_, err := tx.NamedExec(q, param)
	if err != nil {
		return fmt.Errorf("failed to insert rideindego properties: %w", err)
	}

	return nil
}

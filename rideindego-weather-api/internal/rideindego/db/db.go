package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// DBRideIndegoProvider is the interface used to interact with
// ride indego data.
type DBRideIndegoProvider interface{}

// DBRideIndego is used to interact with Ride Indego data
type dbRideIndego struct {
	conn *sqlx.DB
}

func NewDBRideIndego(db *sqlx.DB) DBRideIndegoProvider {
	return &dbRideIndego{
		conn: db,
	}
}

func (db *dbRideIndego) saveTableMaster(tx *sqlx.Tx, param *TableRideIndegoMaster) error {
	q := `
        INSERT INTO rideindego_master
        (
            fetch_id,
            feature_type,
            last_updated
        )
        VALUES
        (
            :fetch_id,
            :feature_type,
            :last_updated
        )
    `

	_, err := tx.NamedExec(q, param)
	if err != nil {
		return fmt.Errorf("failed to insert rideindego master: %w", err)
	}

	return nil
}

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DBWeatherProvider interface {
	StoreDatas(context.Context, *TableWeatherMaster, []*TableWeatherDetails) error
}

type dbWeather struct {
	conn *sqlx.DB
}

var ErrWeatherDB = fmt.Errorf("weather db error")

func NewDBWeather(db *sqlx.DB) DBWeatherProvider {
	return &dbWeather{
		conn: db,
	}
}

// StoreDatas implements DBWeatherProvider.
func (d *dbWeather) StoreDatas(ctx context.Context, dataMaster *TableWeatherMaster, dataDetails []*TableWeatherDetails) error {
	tx, err := d.conn.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Join(ErrWeatherDB, err)
	}

	// save master
	err = d.saveMaster(tx, dataMaster)
	if err != nil {
		tx.Rollback()
		return errors.Join(ErrWeatherDB, err)
	}

	// save details
	err = d.saveDetails(tx, dataDetails)
	if err != nil {
		tx.Rollback()
		return errors.Join(ErrWeatherDB, err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Join(ErrWeatherDB, err)
	}

	return nil
}

func (d *dbWeather) saveMaster(tx *sqlx.Tx, data *TableWeatherMaster) error {
	q := `
        INSERT INTO openweather_master
        (
            fetch_id,
            coord,
            base,
            main_temp,
            main_feels_like,
            main_temp_min,
            main_temp_max,
            main_pressure,
            main_humidity,
            main_sea_level,
            main_grnd_level,
            visibility,
            wind_speed,
            wind_deg,
            wind_gust,
            clouds_all,
            dt,
            sys_type,
            sys_id,
            sys_country,
            sys_sunrise,
            sys_sunset
        )
        VALUES
        (
            :fetch_id,
            :coord,
            :base,
            :main_temp,
            :main_feels_like,
            :main_temp_min,
            :main_temp_max,
            :main_pressure,
            :main_humidity,
            :main_sea_level,
            :main_grnd_level,
            :visibility,
            :wind_speed,
            :wind_deg,
            :wind_gust,
            :clouds_all,
            :dt,
            :sys_type,
            :sys_id,
            :sys_country,
            :sys_sunrise,
            :sys_sunset
        )
    `

	_, err := tx.NamedExec(q, data)
	if err != nil {
		return fmt.Errorf("failed to insert openweather_master: %w", err)
	}

	return nil
}

func (db *dbWeather) saveDetails(tx *sqlx.Tx, data []*TableWeatherDetails) error {
	q := `
        INSERT INTO openweather_details
        (
            fetch_id,
            idx,
            id,
            main,
            description,
            icon
        )
        VALUES
        (
            :fetch_id,
            :idx,
            :id,
            :main,
            :description,
            :icon
        )
    `

	_, err := tx.NamedExec(q, data)
	if err != nil {
		return fmt.Errorf("failed to insert openweather_details: %w", err)
	}

	return nil
}

package db

import (
	"context"
	"fmt"
	"time"
)

type SearchWeatherResult struct {
	Master  *TableWeatherMaster
	Details []*TableWeatherDetails
}

// SearchWeather implements DBWeatherProvider.
func (d *dbWeather) SearchWeather(ctx context.Context, atTime time.Time) (*SearchWeatherResult, error) {
	dataMaster, err := d.searchMaster(ctx, atTime)
	if err != nil {
		return nil, err
	}

	dataDetail, err := d.searchDetail(ctx, dataMaster.FetchID)
	if err != nil {
		return nil, err
	}

	return &SearchWeatherResult{dataMaster, dataDetail}, nil
}

func (d *dbWeather) searchMaster(ctx context.Context, at time.Time) (*TableWeatherMaster, error) {
	q := `
        SELECT
            fetch_id,
            id,
            ST_AsText(coord) AS coord,
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
            sys_sunset,
            timezone,
            name,
            cod
        FROM openweather_master
        WHERE
            dt >= $1
        LIMIT 1
    `

	res := new(TableWeatherMaster)
	err := d.conn.GetContext(ctx, res, q, at.Unix())
	if err != nil {
		return nil, fmt.Errorf("failed to query weather master: %w", err)
	}

	return res, nil
}

func (d *dbWeather) searchDetail(ctx context.Context, fetchID string) ([]*TableWeatherDetails, error) {
	q := `
        SELECT  
            fetch_id,
            idx,
            id,
            main,
            description,
            icon
        FROM openweather_details
        WHERE 
            fetch_id = $1
    `

	res := []*TableWeatherDetails{}
	err := d.conn.SelectContext(ctx, &res, q, fetchID)
	if err != nil {
		return nil, fmt.Errorf("failed to query weather details: %w", err)
	}

	return res, nil
}

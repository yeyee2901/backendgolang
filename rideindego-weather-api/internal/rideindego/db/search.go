package db

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// SearchDataResult is the search result
//
// NOTE: For 'bike' arrays & 'properties' field, we need to keep track of which feature does this belongs to,
// since each feature can have their own set of bike.
// For this, we can use:
// - fetch ID: to identify which data group does this bike belong to
// - feature ID: to identify which feature does this bike belong to
//
// This can be accomplished with a hash map, with feature ID (integer) as the key.
type SearchDataResult struct {
	Master     *TableRideIndegoMaster
	Features   []*TableRideIndegoFeatures
	Properties map[int]*TableRideIndegoProperties
	Bike       map[int][]*TableRideIndegoBikes
}

// SearchData implements DBRideIndegoProvider.
func (db *dbRideIndego) SearchData(c context.Context, at time.Time, kioskID string) (*SearchDataResult, error) {
	// kioskID can be used to optimize DB query, since feature ID
	// can be obtained by joining master with properties
	kioskIDint := -1
	featureID := -1
	searchResult := &SearchDataResult{}

	if len(kioskID) > 0 {
		i, err := strconv.Atoi(kioskID)
		if err != nil {
			return nil, fmt.Errorf("invalid kiosk ID")
		}
		kioskIDint = i
	}

	// search master
	dataMaster, err := db.searchMasterWithFeatureID(c, at, kioskIDint)
	if err != nil {
		return nil, err
	}

	searchResult.Master = &TableRideIndegoMaster{
		FetchID:     dataMaster.FetchID,
		LastUpdated: dataMaster.LastUpdated.UTC(),
		DataType:    dataMaster.DataType,
	}

	// OPTIMIZE: we can obtain the feature ID early if we use kiosk ID. Update the value if so
	if dataMaster.FeatureID > 0 {
		featureID = dataMaster.FeatureID
	}

	// search features
	dataFeatures, err := db.searchFeatures(c, dataMaster.FetchID, featureID)
	if err != nil {
		return nil, err
	}
	searchResult.Features = dataFeatures

	// search properties
	dataProperties, err := db.searchProperties(c, dataMaster.FetchID, featureID)
	if err != nil {
		return nil, err
	}
	searchResult.Properties = dataProperties

	// search bike
	dataBikes, err := db.searchBikes(c, dataMaster.FetchID, featureID)
	if err != nil {
		return nil, err
	}
	searchResult.Bike = dataBikes

	return searchResult, nil
}

type tblMasterWithFeatureID struct {
	TableRideIndegoMaster
	FeatureID int `db:"feature_id"`
}

func (db *dbRideIndego) searchMasterWithFeatureID(c context.Context, at time.Time, kioskID int) (*tblMasterWithFeatureID, error) {
	var (
		q    string
		args = []any{at}
	)

	if kioskID > 0 {
		q = `
            SELECT
                m.fetch_id,
                m.last_updated,
                m.data_type,
                p.feature_id
            FROM rideindego_master m
            LEFT OUTER JOIN rideindego_properties p ON p.fetch_id = m.fetch_id
            WHERE
                m.last_updated >= $1 
                AND p.kiosk_id = $2
            ORDER BY m.last_updated ASC
            LIMIT 1
        `
		args = append(args, kioskID)
	} else {
		q = `
            SELECT
                fetch_id,
                last_updated,
                data_type
            FROM rideindego_master
            WHERE
                last_updated >= $1
            ORDER BY last_updated ASC
            LIMIT 1
        `
	}
	res := new(tblMasterWithFeatureID)
	err := db.conn.GetContext(c, res, q, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db *dbRideIndego) searchFeatures(c context.Context, fetchID string, featureID int) ([]*TableRideIndegoFeatures, error) {
	args := []any{fetchID}
	q := `
        SELECT 
            fetch_id,
            feature_id,
            feature_type,
            geo_type,
            ST_AsText(geo_coord) AS geo_coord
        FROM rideindego_features
        WHERE
            fetch_id = $1
    `

	// using feature ID
	if featureID > 0 {
		args = append(args, featureID)
		q += " AND feature_id = $2 "
	}

	res := []*TableRideIndegoFeatures{}
	err := db.conn.SelectContext(c, &res, q, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db *dbRideIndego) searchBikes(c context.Context, fetchID string, featureID int) (map[int][]*TableRideIndegoBikes, error) {
	results := map[int][]*TableRideIndegoBikes{}
	args := []any{fetchID}
	q := `
        SELECT
            fetch_id,
            feature_id,
            properties_id,
            id,
            dock_number,
            is_electric,
            is_available,
            battery
        FROM rideindego_bikes
        WHERE
            fetch_id = $1
    `

	if featureID > 0 {
		args = append(args, featureID)
		q += " AND feature_id = $2 "
	}

	rows, err := db.conn.QueryxContext(c, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		bike := new(TableRideIndegoBikes)
		err := rows.StructScan(bike)
		if err != nil {
			return nil, err
		}

		// append correctly according to the feature ID
		results[bike.FeatureID] = append(results[bike.FeatureID], bike)
	}

	return results, nil
}

func (db *dbRideIndego) searchProperties(c context.Context, fetchID string, featureID int) (map[int]*TableRideIndegoProperties, error) {
	results := map[int]*TableRideIndegoProperties{}
	args := []any{fetchID}
	q := `
        SELECT
            fetch_id,
            feature_id,
            properties_id,
            ST_AsText(coordinates) as coordinates,
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
        FROM rideindego_properties
        WHERE 
            fetch_id = $1
    `

	if featureID > 0 {
		args = append(args, featureID)
		q += " AND feature_id = $2 "
	}

	rows, err := db.conn.QueryxContext(c, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		prop := new(TableRideIndegoProperties)
		err := rows.StructScan(prop)
		if err != nil {
			return nil, err
		}

		results[prop.FeatureID] = prop
	}

	return results, nil
}

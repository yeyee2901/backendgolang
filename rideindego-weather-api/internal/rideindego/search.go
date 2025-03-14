package rideindego

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/rideindego/db"
)

type SearchParam struct {
	At      time.Time
	KioskID string
}

func (ri *RideIndeGoService) Search(ctx context.Context, param SearchParam) (*APIResponse, error) {
	// search the data
	searchResult, err := ri.store.SearchData(ctx, param.At, param.KioskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Join(err, ErrDataNotFound)
		}

		return nil, err
	}

	apiResp := &APIResponse{}

	// parse data master
	apiResp.LastUpdated = searchResult.Master.LastUpdated
	apiResp.Type = searchResult.Master.DataType

	// parse data features. Take notes of feature ID
	// since it is used to correlate the properties
	features := make([]FeatureElement, len(searchResult.Features))
	for idx, f := range searchResult.Features {
		featureID := f.FeatureID
		coord, err := parseCoordinates(f.GeoCoord)
		if err != nil {
			return nil, err
		}

		// so we don't lost who owns what
		properties, err := parsePropertiesAndBikes(searchResult.Properties[featureID], searchResult.Bike[featureID])
		if err != nil {
			return nil, err
		}

		features[idx] = FeatureElement{
			Geometry: Geometry{
				Coordinates: []float64{coord[0], coord[1]},
				Type:        f.GeoType,
			},
			Type:       f.FeatureType,
			Properties: properties,
		}
	}
	apiResp.Features = features

	return apiResp, nil
}

func parsePropertiesAndBikes(properties *db.TableRideIndegoProperties, bikes []*db.TableRideIndegoBikes) (FeatureProperties, error) {
	longLat, err := parseCoordinates(properties.Coordinates)
	if err != nil {
		return FeatureProperties{}, fmt.Errorf("failed to parse properties coordinate")
	}

	res := FeatureProperties{
		PropertiesID:          properties.PropertiesID,
		Name:                  properties.Name,
		TotalDocks:            properties.TotalDocks,
		DocksAvailable:        properties.DocksAvailable,
		BikesAvailable:        properties.BikesAvailable,
		ClassicBikesAvailable: properties.ClassicBikesAvailable,
		SmartBikesAvailable:   properties.SmartBikesAvailable,
		EletricBikesAvailable: properties.EletricBikesAvailable,
		RewardBikesAvailable:  properties.RewardBikesAvailable,
		RewardDocksAvailable:  properties.RewardDocksAvailable,
		KioskStatus:           properties.KioskStatus,
		KioskPublicStatus:     properties.KioskPublicStatus,
		KioskConnectionStatus: properties.KioskConnectionStatus,
		KioskType:             properties.KioskType,
		AddressStreet:         properties.AddressStreet,
		AddressCity:           properties.AddressCity,
		AddressState:          properties.AddressState,
		AddressZipCode:        properties.AddressZipCode,
		CloseTime:             properties.CloseTime,
		EventEnd:              properties.EventEnd,
		EventStart:            properties.EventStart,
		IsEventBased:          properties.IsEventBased,
		IsVirtual:             properties.IsVirtual,
		KioskID:               properties.KioskID,
		Notes:                 properties.Notes,
		OpenTime:              properties.OpenTime,
		PublicText:            properties.PublicText,
		Timezone:              properties.Timezone,
		TrikesAvailable:       properties.TrikesAvailable,
		Latitude:              longLat[1],
		Longitude:             longLat[0],
		Coordinates:           []float64{longLat[0], longLat[1]},

		Bikes: []Bikes{},
	}

	for i := range bikes {
		res.Bikes = append(res.Bikes, Bikes{
			DockNumber:  bikes[i].DockNumber,
			IsElectric:  bikes[i].IsElectric,
			IsAvailable: bikes[i].IsAvailable,
			Battery:     bikes[i].Battery,
		})
	}

	return res, nil
}

func parseCoordinates(coordString string) ([2]float64, error) {
	lat := 0.0
	long := 0.0
	_, err := fmt.Sscanf(coordString, "POINT(%f %f)", &long, &lat)
	if err != nil {
		return [2]float64{}, fmt.Errorf("failed to parse rideindego coordinate: %w", err)
	}

	return [2]float64{long, lat}, nil
}

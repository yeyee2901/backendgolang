package rideindego

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/httpclient"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/rideindego/db"
)

type RideIndeGoService struct {
	baseURL string
	store   db.DBRideIndegoProvider
}

func NewRideIndeGoService(baseURL string, dbConn *sqlx.DB) *RideIndeGoService {
	return &RideIndeGoService{
		baseURL: baseURL,
		store:   db.NewDBRideIndego(dbConn),
	}
}

// RefreshData fetches a fresh data from the external API & stores it in the database
func (ri *RideIndeGoService) RefreshData() (*APIResponse, error) {
	apiResp, err := ri.fetchData()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	// save to DB
	err = ri.saveData(apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to save data: %w", err)
	}

	return apiResp, nil
}

// fetchData fetches the data from the external API
func (ri *RideIndeGoService) fetchData() (*APIResponse, error) {
	resp := new(APIResponse)
	status, err := httpclient.HTTPRequest(http.MethodGet, nil, ri.baseURL, nil, resp)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %d", status)
	}

	return resp, nil
}

// saveData saves the data to DB based on the API Response
func (ri *RideIndeGoService) saveData(dataToSave *APIResponse) error {
	fetchID := uuid.NewString()
	baseFeatureID := int(dataToSave.LastUpdated.Unix() - 1000)
	dataMaster := &db.TableRideIndegoMaster{
		FetchID:     fetchID,
		LastUpdated: dataToSave.LastUpdated,
		DataType:    dataToSave.Type,
	}

	// OPTIMIZE:
	// prepare sized array wherever possible,
	// since we know it without traversing
	dataFeatures := make([]*db.TableRideIndegoFeatures, len(dataToSave.Features))
	dataProperties := []*db.TableRideIndegoProperties{}
	dataBikes := []*db.TableRideIndegoBikes{}

	for idx, feature := range dataToSave.Features {
		featureID := baseFeatureID + idx

		// OPTIMIZE: use the prepared array
		dataFeatures[idx] = &db.TableRideIndegoFeatures{
			FetchID:     fetchID,
			FeatureID:   featureID,
			FeatureType: feature.Type,
			GeoType:     feature.Geometry.Type,
			GeoCoord:    fmt.Sprintf("POINT(%f %f)", feature.Geometry.Coordinates[0], feature.Geometry.Coordinates[1]),
		}

		dataProperties = append(dataProperties, &db.TableRideIndegoProperties{
			FetchID:               fetchID,
			FeatureID:             featureID,
			PropertiesID:          feature.Properties.PropertiesID,
			Coordinates:           fmt.Sprintf("POINT(%f %f)", feature.Properties.Coordinates[0], feature.Properties.Coordinates[1]),
			Name:                  feature.Properties.Name,
			TotalDocks:            feature.Properties.TotalDocks,
			DocksAvailable:        feature.Properties.DocksAvailable,
			BikesAvailable:        feature.Properties.BikesAvailable,
			ClassicBikesAvailable: feature.Properties.ClassicBikesAvailable,
			SmartBikesAvailable:   feature.Properties.SmartBikesAvailable,
			EletricBikesAvailable: feature.Properties.EletricBikesAvailable,
			RewardBikesAvailable:  feature.Properties.RewardBikesAvailable,
			RewardDocksAvailable:  feature.Properties.RewardDocksAvailable,
			KioskStatus:           feature.Properties.KioskStatus,
			KioskPublicStatus:     feature.Properties.KioskPublicStatus,
			KioskConnectionStatus: feature.Properties.KioskConnectionStatus,
			KioskType:             feature.Properties.KioskType,
			AddressStreet:         feature.Properties.AddressStreet,
			AddressCity:           feature.Properties.AddressCity,
			AddressState:          feature.Properties.AddressState,
			AddressZipCode:        feature.Properties.AddressZipCode,
			CloseTime:             feature.Properties.CloseTime,
			EventEnd:              feature.Properties.EventEnd,
			EventStart:            feature.Properties.EventStart,
			IsEventBased:          feature.Properties.IsEventBased,
			IsVirtual:             feature.Properties.IsVirtual,
			KioskID:               feature.Properties.KioskID,
			Notes:                 feature.Properties.Notes,
			OpenTime:              feature.Properties.OpenTime,
			PublicText:            feature.Properties.PublicText,
			Timezone:              feature.Properties.Timezone,
			TrikesAvailable:       feature.Properties.TrikesAvailable,
		})

		for bikeIdx, bike := range feature.Properties.Bikes {
			dataBikes = append(dataBikes, &db.TableRideIndegoBikes{
				FetchID:      fetchID,
				FeatureID:    featureID,
				PropertiesID: feature.Properties.PropertiesID,
				ID:           bikeIdx,
				DockNumber:   bike.DockNumber,
				IsElectric:   bike.IsElectric,
				IsAvailable:  bike.IsAvailable,
				Battery:      bike.Battery,
			})
		}
	}

	return ri.store.StoreDatas(dataMaster, dataFeatures, dataProperties, dataBikes)
}

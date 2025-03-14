package openweather

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrDataNotFound = fmt.Errorf("openweather: data not found")

type SearchParam struct {
	At time.Time
}

func (w *OpenWeatherService) Search(ctx context.Context, param SearchParam) (*APIResponse, error) {
	dbRes, err := w.db.SearchWeather(ctx, param.At)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Join(err, ErrDataNotFound)
		}

		return nil, err
	}

	apiResp := &APIResponse{
		Base: dbRes.Master.Base,
		Main: Main{
			Temp:      dbRes.Master.MainTemp,
			FeelsLike: dbRes.Master.MainFeelsLike,
			TempMin:   dbRes.Master.MainTempMin,
			TempMax:   dbRes.Master.MainTempMax,
			Pressure:  dbRes.Master.MainPressure,
			Humidity:  dbRes.Master.MainHumidity,
			SeaLevel:  dbRes.Master.MainSeaLevel,
			GrndLevel: dbRes.Master.MainGrndLevel,
		},
		Visibility: dbRes.Master.Visibility,
		Wind: Wind{
			Speed: dbRes.Master.WindSpeed,
			Deg:   dbRes.Master.WindDeg,
			Gust:  dbRes.Master.WindGust,
		},
		Clouds: Clouds{
			All: dbRes.Master.CloudsAll,
		},
		Sys: Sys{
			Type:    dbRes.Master.SysType,
			ID:      dbRes.Master.SysID,
			Country: dbRes.Master.SysCountry,
			Sunrise: dbRes.Master.SysSunrise,
			Sunset:  dbRes.Master.SysSunset,
		},
		Timezone: dbRes.Master.Timezone,
		Name:     dbRes.Master.Name,
		Cod:      dbRes.Master.Cod,
		Dt:       dbRes.Master.Dt,
		ID:       dbRes.Master.ID,

		// --- later filled ---

		// OPTIMIZE: prepare sized array
		Weather: make([]Weather, len(dbRes.Details)),
		Coord:   Coord{},
	}

	// parse coordinate
	coord, err := parseCoordinates(dbRes.Master.Coord)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coordinates: %w", err)
	}
	apiResp.Coord = coord

	// parse weather details
	for i := range dbRes.Details {
		apiResp.Weather[i] = Weather{
			ID:          dbRes.Details[i].ID,
			Main:        dbRes.Details[i].Main,
			Description: dbRes.Details[i].Description,
			Icon:        dbRes.Details[i].Icon,
		}
	}

	return apiResp, nil
}

func parseCoordinates(coordString string) (Coord, error) {
	lat := 0.0
	long := 0.0
	_, err := fmt.Sscanf(coordString, "POINT(%f %f)", &long, &lat)
	if err != nil {
		return Coord{}, err
	}

	return Coord{Lat: lat, Lon: long}, nil
}

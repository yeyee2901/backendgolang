package openweather

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/httpclient"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/openweather/db"
)

type OpenWeatherService struct {
	db      db.DBWeatherProvider
	apiKey  string
	baseURL string
}

func NewOpenWeather(apiKey string, baseURL string, conn *sqlx.DB) *OpenWeatherService {
	return &OpenWeatherService{
		apiKey:  apiKey,
		baseURL: baseURL,
		db:      db.NewDBWeather(conn),
	}
}

func (ow *OpenWeatherService) RefreshData(ctx context.Context) error {
	resp, err := ow.fetchData(ctx)
	if err != nil {
		return err
	}

	return ow.saveData(ctx, resp)
}

func (ow *OpenWeatherService) fetchData(ctx context.Context) (*APIResponse, error) {
	target, err := url.Parse(ow.baseURL)
	if err != nil {
		return nil, err
	}

	q := target.Query()
	q.Add("q", "Philadelphia")
	q.Add("appid", ow.apiKey)
	target.RawQuery = q.Encode()

	resp := new(APIResponse)
	status, err := httpclient.HTTPRequest(ctx, http.MethodGet, nil, target.String(), nil, resp)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, err
	}

	return resp, nil
}

func (ow *OpenWeatherService) saveData(ctx context.Context, dataToSave *APIResponse) error {
	fetchID := uuid.NewString()

	dataMaster := &db.TableWeatherMaster{
		FetchID:       fetchID,
		Coord:         fmt.Sprintf("POINT(%f %f)", dataToSave.Coord.Lat, dataToSave.Coord.Lon),
		Base:          dataToSave.Base,
		MainTemp:      dataToSave.Main.Temp,
		MainFeelsLike: dataToSave.Main.FeelsLike,
		MainTempMin:   dataToSave.Main.TempMin,
		MainTempMax:   dataToSave.Main.TempMax,
		MainPressure:  dataToSave.Main.Pressure,
		MainHumidity:  dataToSave.Main.Humidity,
		MainSeaLevel:  dataToSave.Main.SeaLevel,
		MainGrndLevel: dataToSave.Main.GrndLevel,
		Visibility:    dataToSave.Visibility,
		WindSpeed:     dataToSave.Wind.Speed,
		WindDeg:       dataToSave.Wind.Deg,
		WindGust:      dataToSave.Wind.Gust,
		CloudsAll:     dataToSave.Clouds.All,
		Dt:            dataToSave.Dt,
		SysType:       dataToSave.Sys.Type,
		SysID:         dataToSave.Sys.ID,
		SysCountry:    dataToSave.Sys.Country,
		SysSunrise:    dataToSave.Sys.Sunrise,
		SysSunset:     dataToSave.Sys.Sunset,
	}

	// OPTIMIZE: prepare sized array
	dataDetails := make([]*db.TableWeatherDetails, len(dataToSave.Weather))
	for i := range dataToSave.Weather {
		dataDetails[i] = &db.TableWeatherDetails{
			FetchID:     fetchID,
			IDx:         i,
			ID:          dataToSave.Weather[i].ID,
			Main:        dataToSave.Weather[i].Main,
			Description: dataToSave.Weather[i].Description,
			Icon:        dataToSave.Weather[i].Icon,
		}
	}

	return ow.db.StoreDatas(ctx, dataMaster, dataDetails)
}

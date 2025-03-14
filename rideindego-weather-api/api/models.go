package api

import (
	"time"

	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/openweather"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/rideindego"
)

type APIResponse struct {
	// At represents the actual time of the first snapshot of data on or after the requested time (README)
	At       time.Time               `json:"at"`
	Stations rideindego.APIResponse  `json:"stations"`
	Weather  openweather.APIResponse `json:"weather"`
}

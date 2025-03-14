package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/openweather"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/rideindego"
)

// HandleRefreshData gin handler
// @Summary Refresh RideIndego data
// @Tags API
// @Description fetches new ride-indego & weather data, then stores it in database
// @Param Authorization header string true "Bearer teehee-what-are-you-looking-for"
// @Produce json
// @Consume json
// @Router /api/v1/indego-data-fetch-and-store-it-db [post]
func (api *APIServer) HandleRefreshData(c *gin.Context) {
	logger := slog.Default().With("endpoint", "POST /api/v1/indego-data-fetch-and-store-it-db")
	refreshCtx, cancel := context.WithCancel(c)
	defer cancel()

	doneRideIndego := make(chan struct{})
	doneWeather := make(chan struct{})
	errRideIndegoChan := make(chan error)
	errWeatherChan := make(chan error)

	// refresh ride indego data (concurrent)
	go func() {
		ride := rideindego.NewRideIndeGoService(api.config.RideIndegoBaseURL, api.dbConn)
		_, err := ride.RefreshData(refreshCtx)
		if err != nil {
			errRideIndegoChan <- err
			doneRideIndego <- struct{}{}
			return
		}

		close(errRideIndegoChan)
		doneRideIndego <- struct{}{}
	}()

	// refresh weather data (concurrent)
	go func() {
		weather := openweather.NewOpenWeather(api.config.OpenWeatherAPIKey, api.config.OpenWeatherURL, api.dbConn)
		err := weather.RefreshData(refreshCtx)
		if err != nil {
			errWeatherChan <- err
			doneWeather <- struct{}{}
			return
		}

		close(errWeatherChan)
		doneWeather <- struct{}{}
	}()

	<-doneRideIndego
	<-doneWeather

	// parse & collect errors
	errRideIndego := <-errRideIndegoChan
	errWeather := <-errWeatherChan

	if err := errors.Join(errRideIndego, errWeather); err != nil {
		logger.Error("failed to save", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprint("failed to refresh data:", err),
		})

		return
	}

	logger.Info("refreshed database with fresh data")
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// HandleRefreshData gin handler
// @Summary Get data based on time
// @Tags API
// @Description fetches new ride-indego & weather data, then stores it in database
// @Param Authorization header string true "Bearer teehee-what-are-you-looking-for"
// @Param at query string true "2006-01-02T15:04:05Z"
// @Produce json
// @Consume json
// @Success 200 {object} APIResponse
// @Router /api/v1/stations [GET]
func (api *APIServer) HandleSearchByTime(c *gin.Context) {
	logger := slog.Default().With("endpoint", "POST /api/v1/v1/stations")
	searchCtx, cancel := context.WithCancel(c)
	defer cancel()

	at, err := time.Parse(time.RFC3339, c.Query("at"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid format for 'at' parameter"})
		return
	}

	// TODO: search ride indego data (should return the same JSON as the API)
	ride := rideindego.NewRideIndeGoService(api.config.RideIndegoBaseURL, api.dbConn)
	rideData, err := ride.Search(searchCtx, rideindego.SearchParam{At: at})
	if err != nil {
		if errors.Is(err, rideindego.ErrDataNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "no matching ride indego data found"})
			return
		}

		logger.Error("failed to get data", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	// search weather data (should return the same JSON as the API)
	weather := openweather.NewOpenWeather(api.config.OpenWeatherAPIKey, api.config.OpenWeatherURL, api.dbConn)
	weatherData, err := weather.Search(searchCtx, openweather.SearchParam{At: at})
	if err != nil {
		if errors.Is(err, openweather.ErrDataNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "no matching weather data found"})
			return
		}

		logger.Error("failed to get data", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusOK, APIResponse{
		At:       at,
		Stations: *rideData,
		Weather:  *weatherData,
	})
}

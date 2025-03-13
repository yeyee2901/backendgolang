package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
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

	// TODO: refresh weather data (concurrent)
	go func() {
		// ride := rideindego.NewRideIndeGoService(api.rideindegoBaseURL, api.dbConn)
		// _, err := ride.RefreshData(refreshCtx)
		// if err != nil {
		// 	errors <- err
		// }

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

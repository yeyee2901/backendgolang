package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type APIServer struct {
	gin        *gin.Engine
	httpServer *http.Server
	dbConn     *sqlx.DB

	config *APIConfig
}

type APIConfig struct {
	Listener             string
	SecretKey            string
	RideIndegoBaseURL    string
	ServerTimeoutSeconds int

	OpenWeatherURL    string
	OpenWeatherAPIKey string
}

func NewAPIServer(db *sqlx.DB, config *APIConfig) *APIServer {
	return &APIServer{
		gin:    gin.New(),
		dbConn: db,
		config: config,
	}
}

func (api *APIServer) RegisterMiddlewares() {
	api.gin.Use(gin.Recovery())
	api.gin.Use(CORSMiddleware())
	api.gin.Use(SimpleAuthorizationMiddleware(api.config.SecretKey))
}

func (api *APIServer) RegisterEndpoints() {
	api.gin.POST("/api/v1/indego-data-fetch-and-store-it-db", api.HandleRefreshData)

	// register swagger
	docs.SwaggerInfo.Host = api.config.Listener
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	api.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Run runs the server. This will return an error channel that can
// be waited. This error channel will return a non-nil error
// whenever the server is stopped.
func (api *APIServer) Run() <-chan error {
	errChan := make(chan error)
	httpServer := &http.Server{
		Addr:         api.config.Listener,
		Handler:      api.gin,
		ReadTimeout:  time.Duration(api.config.ServerTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(api.config.ServerTimeoutSeconds) * time.Second,
	}

	go func() {
		fmt.Println("Server listening at:", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	api.httpServer = httpServer

	return errChan
}

// Shutdown kills the HTTP Server entirely
func (s *APIServer) Shutdown() error {
	return s.httpServer.Shutdown(context.Background())
}

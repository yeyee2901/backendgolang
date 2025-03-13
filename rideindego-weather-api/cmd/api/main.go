package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/api"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/utils"
)

// @title		            API Gateway - Ride Indego x Open Weather
// @version		            1.0
// @BasePath	            /
// @description.markdown

func main() {
	err := godotenv.Load("build.env")
	if err != nil {
		fmt.Println("cannot load env var from build.env")
		os.Exit(1)
	}

	db, err := connectDB()
	if err != nil {
		fmt.Println("Cannot connect to database")
		os.Exit(1)
	}

	timeoutSeconds, err := strconv.Atoi(os.Getenv("SERVER_TIMEOUT_SECONDS"))
	if err != nil {
		timeoutSeconds = 30
	}

	cfg := &api.APIConfig{
		Listener:             os.Getenv("SERVER_LISTENER"),
		SecretKey:            os.Getenv("SERVER_SECRET_TOKEN"),
		RideIndegoBaseURL:    os.Getenv("RIDE_INDEGO_URL"),
		ServerTimeoutSeconds: timeoutSeconds,
	}
	fmt.Println(cfg)

	server := api.NewAPIServer(db, cfg)

	server.RegisterMiddlewares()
	server.RegisterEndpoints()

	errChan := server.Run()
	err = <-errChan
	if err != nil {
		fmt.Println("Server exited")
		os.Exit(1)
	}
}

func connectDB() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", utils.BuildDatasourceName(utils.DataSource{
		User:     os.Getenv("POSTGRESQL_USERNAME"),
		Password: os.Getenv("POSTGRESQL_PASSWORD"),
		Host:     os.Getenv("POSTGRESQL_HOST"),
		Database: os.Getenv("POSTGRESQL_DATABASE"),
	}))
}

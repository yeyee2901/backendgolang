package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/httpclient"
	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/logging"
)

func main() {
	logger := logging.NewFileLogger("log/app.log", "refresh-scheduler", slog.LevelInfo)
	slog.SetDefault(logger)
	err := godotenv.Load("build.env")
	if err != nil {
		slog.Error("failed to load .env", "error", err)
		return
	}

	secretKey := os.Getenv("SERVER_SECRET_TOKEN")
	if secretKey == "" {
		slog.Error("failed to load SERVER_SECRET_TOKEN")
		return
	}

	// fire it first time
	err = RefreshData(secretKey)
	if err != nil {
		slog.Error("Failed to refresh for the first time")
		os.Exit(1)
	}

	interval := parseInterval(os.Getenv("SCHEDULER_REFRESH_INTERVAL"))

	c := cron.New()
	c.AddFunc("@every "+interval, func() { RefreshData(secretKey) })
	c.Start()

	// wait forever
	osSig := make(chan os.Signal, 3)
	signal.Notify(osSig, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM)
	<-osSig
}

func RefreshData(secretKey string) error {
	targetURL := os.Getenv("SCHEDULER_REFRESH_URL")
	headersMap := map[string]string{
		"Authorization": "Bearer " + secretKey,
	}
	status, err := httpclient.HTTPRequest(
		context.Background(),
		http.MethodPost,
		headersMap,
		targetURL,
		bytes.NewBuffer([]byte{}),
		nil,
	)
	if err != nil {
		slog.Error("failed to refresh data", "error", err)
		return err
	}

	if status != http.StatusOK {
		errMsg := "failed to refresh data caused by non-200 Status"
		slog.Error(errMsg, "status", status)
		return fmt.Errorf("%s", errMsg)
	}

	slog.Info("refreshed data", "status", status)
	return nil
}

// default is 60 minutes
func parseInterval(in string) string {
	if len(in) < 1 {
		return "1h"
	}

	return in
}

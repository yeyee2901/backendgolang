package openweather

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestFetch(t *testing.T) {
	err := godotenv.Load("../../build.env")
	if err != nil {
		t.Fatal(err)
	}
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	baseURL := os.Getenv("OPENWEATHER_URL")

	weather := openWeatherService{apiKey, baseURL}
	resp, err := weather.fetchData(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// just check for the array availability
	if len(resp.Weather) == 0 {
		t.Log("weather:", resp)
		t.Fatal("empty weather")
	}
}

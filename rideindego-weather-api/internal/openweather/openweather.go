package openweather

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/httpclient"
)

type OpenWeatherProvider interface{}

type openWeatherService struct {
	apiKey  string
	baseURL string
}

func NewOpenWeather(apiKey string, baseURL string) OpenWeatherProvider {
	return &openWeatherService{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

func (ow *openWeatherService) fetchData(ctx context.Context) (*APIResponse, error) {
	target, err := url.Parse(ow.baseURL)
	if err != nil {
		return nil, err
	}

	q := target.Query()
	q.Add("q", "Philadelphia")
	q.Add("appid", ow.apiKey)
	target.RawQuery = q.Encode()

	fmt.Println(target)
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

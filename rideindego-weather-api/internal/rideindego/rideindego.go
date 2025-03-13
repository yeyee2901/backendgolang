package rideindego

import (
	"fmt"
	"net/http"

	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/httpclient"
)

type RideIndeGoService struct {
	baseURL string
}

func NewRideIndeGoService(baseURL string) *RideIndeGoService {
	return &RideIndeGoService{
		baseURL: baseURL,
	}
}

// fetchData fetches the data from the external API
func (ri *RideIndeGoService) fetchData() (*APIResponse, error) {
	resp := new(APIResponse)
	status, err := httpclient.HTTPRequest(http.MethodGet, nil, ri.baseURL, nil, resp)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %d", status)
	}

	return resp, nil
}

package httpclient_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/yeyee2901/backendgolang/rideindego-weather-api/internal/httpclient"
)

func TestHTTPClient(t *testing.T) {
	rideindegoURL := "https://www.rideindego.com/stations/json/"

	resp := map[string]any{}
	status, err := httpclient.HTTPRequest(context.Background(), http.MethodGet, nil, rideindegoURL, nil, &resp)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Status Code: %d", status)
}

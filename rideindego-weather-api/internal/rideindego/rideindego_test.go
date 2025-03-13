package rideindego

import (
	"testing"
)

func TestFetchData(t *testing.T) {
	baseURL := "https://www.rideindego.com/stations/json/"
	ride := NewRideIndeGoService(baseURL)
	resp, err := ride.fetchData()
	if err != nil {
		t.Fatal("failed to fetch:", err)
	}

	if resp == nil {
		t.Fatal("resp should not be nil")
	}

	if len(resp.Features) == 0 {
		t.Fatal("resp.features should not be empty")
	}

	if resp.LastUpdated.IsZero() {
		t.Fatal("resp.lastUpdated should not be zero value")
	}

	if resp.Type == "" {
		t.Fatal("resp.type should not be empty")
	}

	t.Log(resp.LastUpdated, resp.Type)
	t.Logf("response: %+v", resp.Features[1])
}

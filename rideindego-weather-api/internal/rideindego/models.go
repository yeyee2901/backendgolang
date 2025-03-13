package rideindego

import "time"

// APIResponse from fetching the data
type APIResponse struct {
	LastUpdated time.Time        `json:"last_updated" time_format:"2006-01-02T15:04:05.999Z"`
	Type        string           `json:"type"`
	Features    []FeatureElement `json:"features"`
}

// FeatureElement : part of API Response
type FeatureElement struct {
	Geometry   Geometry          `json:"geometry"`
	Properties FeatureProperties `json:"properties"`
}

// Geometry : part of API Response
type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type Bikes struct {
	DockNumber  int  `json:"dockNumber"`
	IsElectric  bool `json:"isElectric"`
	IsAvailable bool `json:"isAvailable"`
	Battery     int  `json:"battery"`
}

// FeatureProperties : part of API Response
type FeatureProperties struct {
	ID                    int       `json:"id"`
	Name                  string    `json:"name"`
	Coordinates           []float64 `json:"coordinates"`
	TotalDocks            int       `json:"totalDocks"`
	DocksAvailable        int       `json:"docksAvailable"`
	BikesAvailable        int       `json:"bikesAvailable"`
	ClassicBikesAvailable int       `json:"classicBikesAvailable"`
	SmartBikesAvailable   int       `json:"smartBikesAvailable"`
	EletricBikesAvailable int       `json:"eletricBikesAvailable"`
	RewardBikesAvailable  int       `json:"rewardBikesAvailable"`
	RewardDocksAvailable  int       `json:"rewardDocksAvailable"`
	KioskStatus           string    `json:"kioskStatus"`
	KioskPublicStatus     string    `json:"kioskPublicStatus"`
	KioskConnectionStatus string    `json:"kioskConnectionStatus"`
	KioskType             int       `json:"kioskType"`
	AddressStreet         string    `json:"addressStreet"`
	AddressCity           string    `json:"addressCity"`
	AddressState          string    `json:"addressState"`
	AddressZipCode        string    `json:"addressZipCode"`
	Bikes                 []Bikes   `json:"bikes"`
	CloseTime             string    `json:"closeTime"`
	EventEnd              string    `json:"eventEnd"`
	EventStart            string    `json:"eventStart"`
	IsEventBased          bool      `json:"isEventBased"`
	IsVirtual             bool      `json:"isVirtual"`
	KioskID               int       `json:"kioskId"`
	Notes                 string    `json:"notes"`
	OpenTime              string    `json:"openTime"`
	PublicText            string    `json:"publicText"`
	Timezone              string    `json:"timezone"`
	TrikesAvailable       int       `json:"trikesAvailable"`
	Latitude              float64   `json:"latitude"`
	Longitude             float64   `json:"longitude"`
}

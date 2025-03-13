package db

import "time"

type TableRideIndegoMaster struct {
	// FetchID : identifies which fetch data group does this corresponds to
	FetchID     string    `db:"fetch_id"`
	LastUpdated time.Time `db:"last_updated"`
	FeatureType string    `db:"feature_type"`
}

type TableRideIndegoBikes struct {
	// FetchID : identifies which fetch data group does this corresponds to
	FetchID string `db:"fetch_id"`

	// FeatureID : identifies which feature index does this bikes belong to
	FeatureID int `db:"feature_id"`

	// PropertiesID : identifies which properties is this.
	// Corresponds to properties.id in APIResponse
	PropertiesID int `db:"properties_id"`

	DockNumber  int  `db:"dock_number"`
	IsElectric  bool `db:"is_electric"`
	IsAvailable bool `db:"is_available"`
	Battery     int  `db:"battery"`
}

type TableRideIndegoProperties struct {
	// FetchID : identifies which fetch data group does this corresponds to
	FetchID string `db:"fetch_id"`

	// FeatureID : identifies which feature index is this
	FeatureID int `db:"feature_id"`

	// PropertiesID : identifies which properties is this.
	// Corresponds to properties.id in APIResponse
	PropertiesID int `db:"id"`

	// Coordinates is POINT type in database, parse it!
	Coordinates string `db:"coordinates"`

	Name                  string `db:"name"`
	TotalDocks            int    `db:"totalDocks"`
	DocksAvailable        int    `db:"docksAvailable"`
	BikesAvailable        int    `db:"bikesAvailable"`
	ClassicBikesAvailable int    `db:"classicBikesAvailable"`
	SmartBikesAvailable   int    `db:"smartBikesAvailable"`
	EletricBikesAvailable int    `db:"eletricBikesAvailable"`
	RewardBikesAvailable  int    `db:"rewardBikesAvailable"`
	RewardDocksAvailable  int    `db:"rewardDocksAvailable"`
	KioskStatus           string `db:"kioskStatus"`
	KioskPublicStatus     string `db:"kioskPublicStatus"`
	KioskConnectionStatus string `db:"kioskConnectionStatus"`
	KioskType             int    `db:"kioskType"`
	AddressStreet         string `db:"addressStreet"`
	AddressCity           string `db:"addressCity"`
	AddressState          string `db:"addressState"`
	AddressZipCode        string `db:"addressZipCode"`
	CloseTime             string `db:"closeTime"`
	EventEnd              string `db:"eventEnd"`
	EventStart            string `db:"eventStart"`
	IsEventBased          bool   `db:"isEventBased"`
	IsVirtual             bool   `db:"isVirtual"`
	KioskID               int    `db:"kioskId"`
	Notes                 string `db:"notes"`
	OpenTime              string `db:"openTime"`
	PublicText            string `db:"publicText"`
	Timezone              string `db:"timezone"`
	TrikesAvailable       int    `db:"trikesAvailable"`
}

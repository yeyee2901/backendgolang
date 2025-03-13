package db

import "time"

type TableRideIndegoMaster struct {
	// FetchID : identifies which fetch data group does this corresponds to
	FetchID     string    `db:"fetch_id"`
	LastUpdated time.Time `db:"last_updated"`
	DataType    string    `db:"data_type"`
}

type TableRideIndegoFeatures struct {
	// FetchID : identifies which fetch data group does this corresponds to
	FetchID string `db:"fetch_id"`

	// FeatureID : identifies which feature index does this bikes belong to
	FeatureID   int    `db:"feature_id"`
	FeatureType string `db:"feature_type"`
	GeoType     string `db:"geo_type"`
	GeoCoord    string `db:"geo_coord"`
}

type TableRideIndegoProperties struct {
	// FetchID : identifies which fetch data group does this corresponds to
	FetchID string `db:"fetch_id"`

	// FeatureID : identifies which feature index is this
	FeatureID int `db:"feature_id"`

	// PropertiesID : identifies which properties is this.
	// Corresponds to properties.id in APIResponse
	PropertiesID int `db:"properties_id"`

	// Coordinates is POINT type in database, parse it!
	Coordinates string `db:"coordinates"`

	Name                  string `db:"name"`
	TotalDocks            int    `db:"total_docks"`
	DocksAvailable        int    `db:"docks_available"`
	BikesAvailable        int    `db:"bikes_available"`
	ClassicBikesAvailable int    `db:"classic_bikes_available"`
	SmartBikesAvailable   int    `db:"smart_bikes_available"`
	EletricBikesAvailable int    `db:"eletric_bikes_available"`
	RewardBikesAvailable  int    `db:"reward_bikes_available"`
	RewardDocksAvailable  int    `db:"reward_docks_available"`
	KioskStatus           string `db:"kiosk_status"`
	KioskPublicStatus     string `db:"kiosk_public_status"`
	KioskConnectionStatus string `db:"kiosk_connection_status"`
	KioskType             int    `db:"kiosk_type"`
	AddressStreet         string `db:"address_street"`
	AddressCity           string `db:"address_city"`
	AddressState          string `db:"address_state"`
	AddressZipCode        string `db:"address_zip_code"`
	CloseTime             string `db:"close_time"`
	EventEnd              string `db:"event_end"`
	EventStart            string `db:"event_start"`
	IsEventBased          bool   `db:"is_event_based"`
	IsVirtual             bool   `db:"is_virtual"`
	KioskID               int    `db:"kiosk_id"`
	Notes                 string `db:"notes"`
	OpenTime              string `db:"open_time"`
	PublicText            string `db:"public_text"`
	Timezone              string `db:"timezone"`
	TrikesAvailable       int    `db:"trikes_available"`
}

type TableRideIndegoBikes struct {
	// FetchID : identifies which fetch group does this bike belongs to
	FetchID string `db:"fetch_id"`

	// FeatureID : identifies which feature does this bike belongs to
	FeatureID int `db:"feature_id"`

	// PropertiesID : identifies which properties does this bike belongs to
	PropertiesID int `db:"properties_id"`

	// ID : bike ID. Unique per instance & time
	ID int `db:"id"`

	DockNumber  int  `db:"dock_number"`
	IsElectric  bool `db:"is_electric"`
	IsAvailable bool `db:"is_available"`
	Battery     int  `db:"battery"`
}

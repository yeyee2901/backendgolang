package db

type TableWeatherMaster struct {
	// FetchID : identifies which data group does this belongs to
	FetchID string `db:"fetch_id"`

	ID            int     `db:"id"`
	Coord         string  `db:"coord"`
	Base          string  `db:"base"`
	MainTemp      float64 `db:"main_temp"`
	MainFeelsLike float64 `db:"main_feels_like"`
	MainTempMin   float64 `db:"main_temp_min"`
	MainTempMax   float64 `db:"main_temp_max"`
	MainPressure  int     `db:"main_pressure"`
	MainHumidity  int     `db:"main_humidity"`
	MainSeaLevel  int     `db:"main_sea_level"`
	MainGrndLevel int     `db:"main_grnd_level"`
	Visibility    int     `db:"visibility"`
	WindSpeed     float64 `db:"wind_speed"`
	WindDeg       int     `db:"wind_deg"`
	WindGust      float64 `db:"wind_gust"`
	CloudsAll     int     `db:"clouds_all"`
	Dt            int     `db:"dt"`
	SysType       int     `db:"sys_type"`
	SysID         int     `db:"sys_id"`
	SysCountry    string  `db:"sys_country"`
	SysSunrise    int     `db:"sys_sunrise"`
	SysSunset     int     `db:"sys_sunset"`
	Timezone      int     `db:"timezone"`
	Name          string  `db:"name"`
	Cod           int     `db:"cod"`
}

type TableWeatherDetails struct {
	// FetchID : identifies which data group does this belongs to
	FetchID string `db:"fetch_id"`

	// IDx : identifies what index was this captured at (when fetched)
	IDx int `db:"idx"`

	ID          int    `db:"id"`
	Main        string `db:"main"`
	Description string `db:"description"`
	Icon        string `db:"icon"`
}

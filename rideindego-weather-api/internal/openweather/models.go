package openweather

type APIResponse struct {
	Coord      Coord  `json:"coord"`
	Base       string `json:"base"`
	Main       Main   `json:"main"`
	Visibility int    `json:"visibility"`
	Wind       Wind   `json:"wind"`
	Clouds     Clouds `json:"clouds"`

	Sys      Sys    `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`

	// Dt : is datetime, looks like in unix seconds
	Dt int `json:"dt"`

	// Weather : the weather data
	Weather []Weather `json:"weather"`
}

// Weather : Weather data
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Coord struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type Clouds struct {
	All int `json:"all"`
}

type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

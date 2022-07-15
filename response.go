package openweathermap

// OneCallAPIResponse maps to a JSON response from OpenWeatherMap's OneCallAPI.
type OneCallAPIResponse struct {
	Lat            float64             `json:"lat"`
	Lon            float64             `json:"lon"`
	Timezone       string              `json:"timezone"`
	TimezoneOffset int                 `json:"timezone_offset"`
	Current        PointWeatherSummary `json:"current"`
	Minutely       []struct {
		Dt            int64 `json:"dt"`
		Precipitation int   `json:"precipitation"`
	} `json:"minutely"`
	Hourly []PointWeatherSummary `json:"hourly"`
	Alerts []struct {
		SenderName  string `json:"sender_name"`
		Event       string `json:"event"`
		Start       int64  `json:"start"`
		End         int64  `json:"end"`
		Description string `json:"description"`
	}
}

// OneCallAPIFailedResponse is used when a request has failed.
type OneCallAPIFailedResponse struct {
	Cod     int    `json:"cod"`
	Message string `json:"message"`
}

// PointWeatherSummary is a subfield of OneCallAPIResponse.
type PointWeatherSummary struct {
	CommonWeatherSummary
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
}

// CommonWeatherSummary is the common part between PointWeatherSummary and
// DailyWeatherSummary.
type CommonWeatherSummary struct {
	Dt         int64   `json:"dt"`
	Sunrise    int     `json:"sunrise"`
	Sunset     int     `json:"sunset"`
	Pressure   int     `json:"pressure"`
	Humidity   int     `json:"humidity"`
	DewPoint   float64 `json:"dew_point"`
	UVI        float64 `json:"uvi"`
	Clouds     int     `json:"clouds"`
	Visibility int     `json:"visibility"`
	WindSpeed  float64 `json:"wind_speed"`
	WindDeg    int     `json:"wind_deg"`
	WindGust   float64 `json:"wind_gust"`
	Pop        float64 `json:"pop"`
	Rain       struct {
		OneHour int `json:"1h"`
	} `json:"rain"`
	Snow struct {
		OneHour int `json:"1h"`
	} `json:"snow"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
}

// DailyWeatherSummary is a subfield of OneCallAPIResponse.
type DailyWeatherSummary struct {
	CommonWeatherSummary
	Temp struct {
		Day   float64 `json:"day"`
		Min   float64 `json:"min"`
		Max   float64 `json:"max"`
		Night float64 `json:"night"`
		Eve   float64 `json:"eve"`
		Morn  float64 `json:"morn"`
	} `json:"temp"`
	FeelsLike struct {
		Day   float64 `json:"day"`
		Night float64 `json:"night"`
		Eve   float64 `json:"eve"`
		Morn  float64 `json:"morn"`
	} `json:"feels_like"`
}

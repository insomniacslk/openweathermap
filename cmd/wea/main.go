package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/insomniacslk/openweathermap"
	"github.com/spf13/pflag"
)

var (
	flagAppID   = pflag.StringP("app-id", "a", "", "App ID (a.k.a API key)")
	flagLat     = pflag.Float64P("latitude", "l", 0.0, "Latitude")
	flagLon     = pflag.Float64P("longitude", "L", 0.0, "Longitude")
	flagExclude = pflag.StringP("exclude", "e", "", "Comma-separated list of fields to exclude from the response")
	flagUnits   = pflag.StringP("units", "u", "standard", "Units to request for response")
	flagLang    = pflag.StringP("language", "g", string(openweathermap.EN), "Language to request for response")
)

func main() {
	pflag.Parse()
	var excludes []openweathermap.Exclude
	for _, e := range strings.Split(*flagExclude, ",") {
		excludes = append(excludes, openweathermap.Exclude(e))
	}
	resp, err := openweathermap.Request(
		*flagAppID,
		*flagLat,
		*flagLon,
		excludes,
		openweathermap.Units(*flagUnits),
		openweathermap.Lang(*flagLang),
	)
	if err != nil {
		log.Fatal(err)
	}

	// prepare units
	var tempUnit, speedUnit string
	switch openweathermap.Units(*flagUnits) {
	case openweathermap.Standard:
		tempUnit = "K"
		speedUnit = "m/s"
	case openweathermap.Metric:
		tempUnit = "C"
		speedUnit = "m/s"
	case openweathermap.Imperial:
		tempUnit = "F"
		speedUnit = "mph"
	}

	// print location information
	fmt.Fprintf(os.Stdout, "Latitude                    : %f\n", resp.Lat)
	fmt.Fprintf(os.Stdout, "Longitude                   : %f\n", resp.Lon)
	fmt.Fprintf(os.Stdout, "Time zone                   : %s\n", resp.Timezone)
	fmt.Fprintf(os.Stdout, "Time zone offset            : %ds\n", resp.TimezoneOffset)
	// print current weather
	if resp.Current != nil {
		fmt.Fprintf(os.Stdout, "Current\n")
		fmt.Fprintf(os.Stdout, "  Temperature               : %.02f%s\n", resp.Current.Temp, tempUnit)
		fmt.Fprintf(os.Stdout, "  Feels like                : %.02f%s\n", resp.Current.FeelsLike, tempUnit)
		printCommonWeatherSummary(os.Stdout, &resp.Current.CommonWeatherSummary, tempUnit, speedUnit)
	}
	// print minutely weather
	fmt.Fprintf(os.Stdout, "Minutely\n")
	for _, m := range resp.Minutely {
		fmt.Fprintf(os.Stdout, "  Timestamp                 : %.02f%s\n", time.Unix(m.Dt, 0))
		fmt.Fprintf(os.Stdout, "  Precipitation             : %.02f%s\n", m.Precipitation)
	}
	// print hourly weather
	fmt.Fprintf(os.Stdout, "Hourly\n")
	for _, hourly := range resp.Hourly {
		fmt.Fprintf(os.Stdout, "  Temperature               : %.02f%s\n", hourly.Temp, tempUnit)
		fmt.Fprintf(os.Stdout, "  Feels like                : %.02f%s\n", hourly.FeelsLike, tempUnit)
		printCommonWeatherSummary(os.Stdout, &hourly.CommonWeatherSummary, tempUnit, speedUnit)
	}
	// print daily weather
	fmt.Fprintf(os.Stdout, "Daily\n")
	for _, daily := range resp.Daily {
		fmt.Fprintf(os.Stdout, "  Temperature (day)         : %.02f%s\n", daily.Temp.Day, tempUnit)
		fmt.Fprintf(os.Stdout, "  Temperature (min)         : %.02f%s\n", daily.Temp.Min, tempUnit)
		fmt.Fprintf(os.Stdout, "  Temperature (max)         : %.02f%s\n", daily.Temp.Max, tempUnit)
		fmt.Fprintf(os.Stdout, "  Temperature (morning)     : %.02f%s\n", daily.Temp.Morn, tempUnit)
		fmt.Fprintf(os.Stdout, "  Temperature (evening)     : %.02f%s\n", daily.Temp.Eve, tempUnit)
		fmt.Fprintf(os.Stdout, "  Temperature (night)       : %.02f%s\n", daily.Temp.Night, tempUnit)
		fmt.Fprintf(os.Stdout, "  Feels like                : %.02f%s\n", daily.FeelsLike, tempUnit)
		printCommonWeatherSummary(os.Stdout, &daily.CommonWeatherSummary, tempUnit, speedUnit)
	}
	// print alerts
	fmt.Fprintf(os.Stdout, "Alerts\n")
	for _, alert := range resp.Alerts {
		fmt.Fprintf(os.Stdout, "  Sender name               : %s\n", alert.SenderName)
		fmt.Fprintf(os.Stdout, "  Event                     : %s\n", alert.Event)
		fmt.Fprintf(os.Stdout, "  Description               : %s\n", alert.Description)
		fmt.Fprintf(os.Stdout, "  Start                     : %s\n", time.Unix(alert.Start, 0))
		fmt.Fprintf(os.Stdout, "  End                       : %s\n", time.Unix(alert.End, 0))
	}
}

// function to print the common part of weather summaries.
func printCommonWeatherSummary(w io.Writer, s *openweathermap.CommonWeatherSummary, tempUnit, speedUnit string) {
	fmt.Fprintf(w, "  Timestamp                 : %s\n", time.Unix(s.Dt, 0))
	fmt.Fprintf(w, "  Sunrise                   : %s\n", time.Unix(int64(s.Sunrise), 0))
	fmt.Fprintf(w, "  Sunset                    : %s\n", time.Unix(int64(s.Sunset), 0))
	fmt.Fprintf(w, "  Pressure                  : %d hPa\n", s.Pressure)
	fmt.Fprintf(w, "  Humidity                  : %d%%\n", s.Humidity)
	fmt.Fprintf(w, "  Dew point                 : %.02f%s\n", s.DewPoint, tempUnit)
	fmt.Fprintf(w, "  UV index                  : %.02f\n", s.UVI)
	fmt.Fprintf(w, "  Clouds                    : %d%%\n", s.Clouds)
	fmt.Fprintf(w, "  Visibility                : %dm\n", s.Visibility)
	fmt.Fprintf(w, "  Wind speed                : %.02f %s\n", s.WindSpeed, speedUnit)
	fmt.Fprintf(w, "  Wind direction            : %.02f degrees\n", s.WindSpeed)
	if s.WindGust != nil {
		fmt.Fprintf(w, "  Wind gust                 : %.02f %s\n", *s.WindGust, speedUnit)
	}
	if s.Rain != nil {
		fmt.Fprintf(w, "  Rain (last hour)          : %.02f mm\n", s.Rain.OneHour)
	}
	if s.Snow != nil {
		fmt.Fprintf(w, "  Snow (last hour)          : %.02f mm\n", s.Snow.OneHour)
	}
	fmt.Fprintf(w, "  Precipitation probability : %.02f%%\n", s.Pop)
	for _, wea := range s.Weather {
		fmt.Fprintf(w, "  Weather condition         : %s: %s\n", wea.Main, wea.Description)
	}
}

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
	flagDebug   = pflag.BoolP("debug", "d", false, "Enable debug output")
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
		*flagDebug,
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
	tz := time.FixedZone(resp.Timezone, resp.TimezoneOffset)
	w := os.Stdout

	// print location information
	fmt.Fprintf(w, "Latitude                    : %f\n", resp.Lat)
	fmt.Fprintf(w, "Longitude                   : %f\n", resp.Lon)
	fmt.Fprintf(w, "Time zone                   : %s\n", resp.Timezone)
	fmt.Fprintf(w, "Time zone offset            : %ds\n", resp.TimezoneOffset)
	fmt.Fprintf(w, "\n")
	// print current weather
	if resp.Current != nil {
		fmt.Fprintf(w, "Current\n")
		fmt.Fprintf(w, "  Temperature               : %.02f%s\n", resp.Current.Temp, tempUnit)
		fmt.Fprintf(w, "  Feels like                : %.02f%s\n", resp.Current.FeelsLike, tempUnit)
		printCommonWeatherSummary(w, &resp.Current.CommonWeatherSummary, tempUnit, speedUnit, tz)
		if resp.Current.Rain.OneHour != nil {
			fmt.Fprintf(w, "  Rain (last hour)          : %.02f mm\n", *resp.Current.Rain.OneHour)
		}
		if resp.Current.Snow.OneHour != nil {
			fmt.Fprintf(w, "  Snow (last hour)          : %.02f mm\n", *resp.Current.Snow.OneHour)
		}
		fmt.Fprintf(w, "\n")
	}
	// print minutely weather
	fmt.Fprintf(w, "Minutely\n")
	for _, minutely := range resp.Minutely {
		fmt.Fprintf(w, "  Timestamp                 : %s\n", time.Unix(minutely.Dt, 0))
		fmt.Fprintf(w, "  Precipitation             : %d mm\n", minutely.Precipitation)
		fmt.Fprintf(w, "\n")
	}
	// print hourly weather
	fmt.Fprintf(w, "Hourly\n")
	for _, hourly := range resp.Hourly {
		fmt.Fprintf(w, "  Temperature               : %.02f%s\n", hourly.Temp, tempUnit)
		fmt.Fprintf(w, "  Feels like                : %.02f%s\n", hourly.FeelsLike, tempUnit)
		printCommonWeatherSummary(w, &hourly.CommonWeatherSummary, tempUnit, speedUnit, tz)
		if hourly.Rain.OneHour != nil {
			fmt.Fprintf(w, "  Rain (last hour)          : %.02f mm\n", *hourly.Rain.OneHour)
		}
		if hourly.Snow.OneHour != nil {
			fmt.Fprintf(w, "  Snow (last hour)          : %.02f mm\n", *hourly.Snow.OneHour)
		}
		fmt.Fprintf(w, "\n")
	}
	// print daily weather
	fmt.Fprintf(w, "Daily\n")
	for _, daily := range resp.Daily {
		fmt.Fprintf(w, "  Temperature (day)         : %.02f%s\n", daily.Temp.Day, tempUnit)
		fmt.Fprintf(w, "  Temperature (min)         : %.02f%s\n", daily.Temp.Min, tempUnit)
		fmt.Fprintf(w, "  Temperature (max)         : %.02f%s\n", daily.Temp.Max, tempUnit)
		fmt.Fprintf(w, "  Temperature (morning)     : %.02f%s\n", daily.Temp.Morn, tempUnit)
		fmt.Fprintf(w, "  Temperature (evening)     : %.02f%s\n", daily.Temp.Eve, tempUnit)
		fmt.Fprintf(w, "  Temperature (night)       : %.02f%s\n", daily.Temp.Night, tempUnit)
		fmt.Fprintf(w, "  Feels like                : %.02f%s\n", daily.FeelsLike, tempUnit)
		printCommonWeatherSummary(w, &daily.CommonWeatherSummary, tempUnit, speedUnit, tz)
		if daily.Rain != nil {
			fmt.Fprintf(w, "  Rain (last hour)          : %.02f mm\n", *daily.Rain)
		}
		if daily.Snow != nil {
			fmt.Fprintf(w, "  Snow (last hour)          : %.02f mm\n", *daily.Snow)
		}
		fmt.Fprintf(w, "\n")
	}
	// print alerts
	fmt.Fprintf(w, "Alerts\n")
	for _, alert := range resp.Alerts {
		fmt.Fprintf(w, "  Sender name               : %s\n", alert.SenderName)
		fmt.Fprintf(w, "  Event                     : %s\n", alert.Event)
		fmt.Fprintf(w, "  Description               : %s\n", alert.Description)
		fmt.Fprintf(w, "  Start                     : %s\n", time.Unix(alert.Start, 0))
		fmt.Fprintf(w, "  End                       : %s\n", time.Unix(alert.End, 0))
		fmt.Fprintf(w, "\n")
	}
}

// function to print the common part of weather summaries.
func printCommonWeatherSummary(w io.Writer, s *openweathermap.CommonWeatherSummary, tempUnit, speedUnit string, tz *time.Location) {
	fmt.Fprintf(w, "  Timestamp                 : %s\n", time.Unix(s.Dt, 0))
	fmt.Fprintf(w, "  Sunrise                   : %s\n", time.Unix(int64(s.Sunrise), 0).In(tz).Format("15:04:05"))
	fmt.Fprintf(w, "  Sunset                    : %s\n", time.Unix(int64(s.Sunset), 0).In(tz).Format("15:04:05"))
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
	fmt.Fprintf(w, "  Precipitation probability : %.02f%%\n", s.Pop)
	for _, wea := range s.Weather {
		fmt.Fprintf(w, "  Weather condition         : %s: %s\n", wea.Main, wea.Description)
	}
}

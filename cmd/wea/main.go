package main

import (
	"fmt"
	"log"
	"strings"

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
	fmt.Printf("%+v\n", resp)
}

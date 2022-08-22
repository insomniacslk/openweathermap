package main

import (
	"fmt"
	"log"
	"os"

	"github.com/insomniacslk/openweathermap/geocoding"
	"github.com/spf13/pflag"
)

var (
	flagAppID     = pflag.StringP("app-id", "a", "", "App ID (a.k.a API key)")
	flagCity      = pflag.StringP("city", "c", "", "City name (only direct geocoding request)")
	flagState     = pflag.StringP("state", "s", "", "State name (only direct geocoding request)")
	flagCountry   = pflag.StringP("country", "C", "", "Country name (only direct geocoding request)")
	flagLatitude  = pflag.Float64P("lat", "l", 0.0, "Latitude (only reverse geocoding request)")
	flagLongitude = pflag.Float64P("lon", "L", 0.0, "Longitude (only reverse geocoding request)")
	flagLimit     = pflag.IntP("limit", "m", geocoding.DefaultLimit, "Maximum number of results (between 0 and 5, 0 means default)")
	flagDebug     = pflag.BoolP("debug", "d", false, "Enable debug output")
	flagReverse   = pflag.BoolP("reverse", "r", false, "Do a reverse geocoding request instead of a direct one. Requires --lat/--lon")
)

func main() {
	pflag.Parse()
	var (
		err  error
		resp *geocoding.Response
	)
	if *flagReverse {
		resp, err = geocoding.ReverseGeocoding(
			*flagAppID,
			&geocoding.ReverseGeocodingRequest{
				Lat:   *flagLatitude,
				Lon:   *flagLongitude,
				Limit: *flagLimit,
			},
			*flagDebug,
		)
	} else {
		resp, err = geocoding.DirectGeocoding(
			*flagAppID,
			&geocoding.DirectGeocodingRequest{
				City:        *flagCity,
				State:       *flagState,
				CountryCode: *flagCountry,
				Limit:       *flagLimit,
			},
			*flagDebug,
		)
	}
	if err != nil {
		log.Fatal(err)
	}
	w := os.Stdout
	for _, item := range *resp {
		fmt.Fprintf(w, "Name:      %s\n", item.Name)
		fmt.Fprintf(w, "Country:   %s\n", item.Country)
		fmt.Fprintf(w, "Latitude:  %f\n", item.Lat)
		fmt.Fprintf(w, "Longitude: %f\n", item.Lon)
		fmt.Fprintf(w, "Got %d localized names\n", len(item.LocalNames))
	}
}

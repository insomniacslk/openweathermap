package main

import (
	"fmt"
	"log"
	"os"

	"github.com/insomniacslk/openweathermap"
	"github.com/insomniacslk/openweathermap/find"
	"github.com/spf13/pflag"
)

var (
	flagAppID = pflag.StringP("app-id", "a", "", "App ID (a.k.a API key)")
	flagQuery = pflag.StringP("query", "q", "", "Query string")
	flagUnits = pflag.StringP("units", "u", "standard", "Units to request for response")
	flagDebug = pflag.BoolP("debug", "d", false, "Enable debug output")
)

func main() {
	pflag.Parse()
	resp, err := find.Request(
		*flagAppID,
		*flagQuery,
		openweathermap.Units(*flagUnits),
		*flagDebug,
	)
	if err != nil {
		log.Fatal(err)
	}
	w := os.Stdout
	fmt.Fprintf(w, "Count:         %d\n", len((*resp).List))
	for _, item := range (*resp).List {
		fmt.Fprintf(w, "ID:        %d\n", item.ID)
		fmt.Fprintf(w, "Name:      %s\n", item.Name)
		fmt.Fprintf(w, "Latitude:  %f\n", item.Coord.Lat)
		fmt.Fprintf(w, "Longitude: %f\n", item.Coord.Lon)
		// TODO print the remaining fields
	}
}

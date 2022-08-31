package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/insomniacslk/openweathermap/airpollution"
	"github.com/spf13/pflag"
)

var (
	flagAppID = pflag.StringP("app-id", "a", "", "App ID (a.k.a API key)")
	flagLat   = pflag.Float64P("lat", "l", 0.0, "Latitude")
	flagLon   = pflag.Float64P("lon", "L", 0.0, "Longitude")
	flagStart = pflag.Int64P("start", "s", 0, "Start time (UNIX timestamp)")
	flagEnd   = pflag.Int64P("end", "e", 0, "End time (UNIX timestamp)")
	flagDebug = pflag.BoolP("debug", "d", false, "Enable debug output")
)

func main() {
	pflag.Parse()
	var (
		err  error
		resp *airpollution.Response
	)
	resp, err = airpollution.Request(
		*flagAppID,
		&airpollution.AirPollutionRequest{
			Lat:   *flagLat,
			Lon:   *flagLon,
			Start: *flagStart,
			End:   *flagEnd,
		},
		*flagDebug,
	)
	if err != nil {
		log.Fatal(err)
	}
	w := os.Stdout
	fmt.Fprintf(w, "Coord             : %f,%f\n", resp.Coord.Lat, resp.Coord.Lon)
	for idx, item := range resp.List {
		fmt.Fprintf(w, "Item #%d\n", idx+1)
		fmt.Fprintf(w, "Datetime          : %s\n", time.Unix(item.Dt, 0))
		fmt.Fprintf(w, "Air Quality Index : %d\n", item.Main.AQI)
		fmt.Fprintf(w, "CO                : %.03f\n", item.Components.CO)
		fmt.Fprintf(w, "NO                : %.03f\n", item.Components.NO)
		fmt.Fprintf(w, "NO2               : %.03f\n", item.Components.NO2)
		fmt.Fprintf(w, "O3                : %.03f\n", item.Components.O3)
		fmt.Fprintf(w, "SO2               : %.03f\n", item.Components.SO2)
		fmt.Fprintf(w, "PM2.5             : %.03f\n", item.Components.PM25)
		fmt.Fprintf(w, "PM10              : %.03f\n", item.Components.PM10)
		fmt.Fprintf(w, "NH3               : %.03f\n", item.Components.NH3)
	}
}

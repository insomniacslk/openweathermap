package openweathermap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// base URL for the One Call API 3.0, see https://openweathermap.org/api/one-call-3
// FIXME check that 3.0 works
// const baseURL = "https://api.openweathermap.org/data/3.0/onecall"
const baseURL = "https://api.openweathermap.org/data/2.5/onecall"

// Request uses OpenWeatherMap's One Call API 3.0.
func Request(appID string, lat, lon float64, exclude []Exclude, units Units, lang Lang, debug bool) (*Weather, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	q := u.Query()
	q.Set("lat", strconv.FormatFloat(lat, 'f', 3, 32))
	q.Set("lon", strconv.FormatFloat(lon, 'f', 3, 32))
	q.Set("appid", appID)
	if len(exclude) != 0 {
		var excludes string
		excludeStrings := make([]string, 0, len(exclude))
		for _, e := range exclude {
			excludeStrings = append(excludeStrings, string(e))
		}
		excludes = strings.Join(excludeStrings, ",")
		q.Set("exclude", excludes)
	}
	if units != "" {
		q.Set("units", string(units))
	}
	if lang != "" {
		q.Set("lang", string(lang))
	}
	u.RawQuery = q.Encode()

	log.Printf("%s", u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP GET failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP body: %w", err)
	}
	if resp.StatusCode != 200 {
		var fresp OneCallAPIFailedResponse
		if err := json.Unmarshal(body, &fresp); err != nil {
			return nil, fmt.Errorf("HTTP GET returned status '%s', and could not unmarshal response message: %w", resp.Status, err)
		}
		return nil, fmt.Errorf("Request failed with %d: %s", fresp.Cod, fresp.Message)
	}

	if debug {
		fmt.Fprintf(os.Stderr, "%s\n", string(body))
	}
	var apiResp Weather
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return &apiResp, nil
}

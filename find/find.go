// Package find implements OpenWeatherMap's find API which seems to be part of
// the data API, but the find method is not documented, see
// https://openweathermap.org/current.
package find

// TODO this should implement the entire data API, plus the undocumented find
// method.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/insomniacslk/openweathermap"
)

var baseURL = url.URL{
	Scheme: "https",
	Host:   "api.openweathermap.org",
	Path:   "/data/2.5/",
}

// Response represents a data find response.
type Response struct {
	Cod     string `json:"cod"`
	Message string `json:"message"`
	Count   int    `json:"count"`
	List    []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		} `json:"main"`
		Dt   int `json:"dt"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		} `json:"wind"`
		Sys struct {
			Country string `json:"country"`
		} `json:"sys"`
		// Rain undefined for now
		// Snow undefined for now
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"list"`
	Name       string            `json:"name"`
	LocalNames map[string]string `json:"local_names"`
	Lat        float64           `json:"lat"`
	Lon        float64           `json:"lon"`
	Country    string            `json:"country"`
	State      string            `json:"state"`
}

// Request executes a data find request.
func Request(appID, query string, units openweathermap.Units, debug bool) (*Response, error) {
	u := baseURL // copy
	u.Path += "find"
	q := u.Query()
	q.Set("appid", appID)
	q.Set("q", query)
	q.Set("units", string(units))
	u.RawQuery = q.Encode()
	if debug {
		fmt.Fprintf(os.Stderr, "URL: %s", u.String())
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP GET failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP body: %w", err)
	}
	if debug {
		fmt.Fprintf(os.Stderr, "Response: %s", string(body))
	}

	var apiResp Response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Request failed with %s: %s", apiResp.Cod, apiResp.Message)
	}

	if debug {
		fmt.Fprintf(os.Stderr, "%s\n", string(body))
	}
	return &apiResp, nil
}

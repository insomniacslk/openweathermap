// Package airpollution implements OpenWeatherMap's air pollution API
// described at https://openweathermap.org/api/air-pollution .
package airpollution

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var baseURL = url.URL{
	Scheme: "https",
	Host:   "api.openweathermap.org",
	Path:   "/data/2.5/air_pollution",
}

// DefaultLimit is the default maximum number of geocoding results to be
// returned.
const DefaultLimit = 5

// FailureResponse is the response structure used when an API call has failed.
type FailureResponse struct {
	Cod     int    `json:"cod"`
	Message string `json:"message"`
}

type AirPollutionRequest struct {
	Lat   float64
	Lon   float64
	Start int64
	End   int64
}

// Response represents an air pollution response.
type Response struct {
	Coord struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
	List []struct {
		Dt   int64 `json:"dt"`
		Main struct {
			AQI int `json:"aqi"`
		} `json:"main"`
		Components struct {
			CO   float64 `json:"co"`
			NO   float64 `json:"no"`
			NO2  float64 `json:"no2"`
			O3   float64 `json:"o3"`
			SO2  float64 `json:"so2"`
			PM25 float64 `json:"pm2_5"`
			PM10 float64 `json:"pm10"`
			NH3  float64 `json:"nh3"`
		} `json:"components"`
	} `json:"list"`
}

// Request executes an air pollution request.
func Request(appID string, req *AirPollutionRequest, debug bool) (*Response, error) {
	u := baseURL // copy
	q := u.Query()
	q.Set("lat", strconv.FormatFloat(req.Lat, 'f', -1, 32))
	q.Set("lon", strconv.FormatFloat(req.Lon, 'f', -1, 32))
	if req.Start != 0 {
		q.Set("start", strconv.FormatInt(req.Start, 10))
	}
	if req.End != 0 {
		q.Set("end", strconv.FormatInt(req.End, 10))
	}
	q.Set("appid", appID)
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
	// first check if the call has failed
	if resp.StatusCode != 200 {
		var fresp FailureResponse
		if err := json.Unmarshal(body, &fresp); err != nil {
			return nil, fmt.Errorf("HTTP GET returned status '%s', and could not unmarshal response message: %w", resp.Status, err)
		}
		return nil, fmt.Errorf("Request failed with %d: %s", fresp.Cod, fresp.Message)
	}

	var apiResp Response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	if debug {
		fmt.Fprintf(os.Stderr, "%s\n", string(body))
	}
	return &apiResp, nil
}

// Package geocoding implements OpenWeatherMap's direct and reverse geocoding API
// described at https://openweathermap.org/api/geocoding-api .
package geocoding

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
	Path:   "/geo/1.0/",
}

// DefaultLimit is the default maximum number of geocoding results to be
// returned.
const DefaultLimit = 5

// FailureResponse is the response structure used when an API call has failed.
type FailureResponse struct {
	Cod     int    `json:"cod"`
	Message string `json:"message"`
}

// DirectGeocodingRequest represents a direct geocoding request.
type DirectGeocodingRequest struct {
	City        string
	State       string
	CountryCode string
	Limit       int
}

// Response represents a direct geocoding response.
type Response []struct {
	Name       string            `json:"name"`
	LocalNames map[string]string `json:"local_names"`
	Lat        float64           `json:"lat"`
	Lon        float64           `json:"lon"`
	Country    string            `json:"country"`
}

// ReverseGeocodingRequest represents a reverse geocoding request.
type ReverseGeocodingRequest struct {
	Lat   float64
	Lon   float64
	Limit int
}

// DirectGeocoding executes a direct geocoding request.
func DirectGeocoding(appID string, req *DirectGeocodingRequest, debug bool) (*Response, error) {
	u := baseURL // copy
	u.Path += "direct"
	q := u.Query()
	q.Set("q", fmt.Sprintf("%s,%s,%s", req.City, req.State, req.CountryCode))
	q.Set("limit", strconv.FormatInt(int64(req.Limit), 10))
	u.RawQuery = q.Encode()

	body, err := request(appID, req.Limit, &u, debug)
	if err != nil {
		return nil, err
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

func request(appID string, limit int, u *url.URL, debug bool) ([]byte, error) {
	q := u.Query()
	q.Set("appid", appID)
	if limit == 0 {
		limit = DefaultLimit
	}
	q.Set("limit", strconv.FormatInt(int64(limit), 10))
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

	return body, nil
}

// ReverseGeocoding executes a reverse geocoding request.
func ReverseGeocoding(appID string, req *ReverseGeocodingRequest, debug bool) (*Response, error) {
	u := baseURL // copy
	u.Path += "reverse"
	q := u.Query()
	q.Set("lat", strconv.FormatFloat(req.Lat, 'g', -1, 32))
	q.Set("lon", strconv.FormatFloat(req.Lon, 'g', -1, 32))
	u.RawQuery = q.Encode()

	body, err := request(appID, req.Limit, &u, debug)
	if err != nil {
		return nil, err
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

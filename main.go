package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Image data
type Image struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

// Observation data
type Observation struct {
	StationID     string  `json:"station_id"`
	TimeString    string  `json:"observation_time"`
	Condition     string  `json:"weather"`
	ConditionURL  string  `json:"icon_url"`
	Temperature   string  `json:"temperature_string"`
	Fahrenheit    float64 `json:"temp_f"`
	Celsius       float64 `json:"temp_c"`
	Wind          string  `json:"wind_string"`
	WindDegrees   float64 `json:"wind_degrees"`
	WindSpeed     float64 `json:"wind_mph"`
	FeelsLike     string  `json:"feelslike_string"`
	Precipitation string  `json:"precip_today_in"`
}

// WeatherResponse is the response data
type WeatherResponse struct {
	Response           MetaInfo    `json:"response"`
	CurrentObservation Observation `json:"current_observation"`
}

// MetaInfo is the Metadata contained in response
type MetaInfo struct {
	Version string `json:"version"`
	Terms   string `json:"termsofService"`
}

// GetURL takes a URL and returns the result as []byte
func GetURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting url: %s", err)
	}
	defer resp.Body.Close()

	buffer := bytes.NewBuffer(nil)

	_, err = io.Copy(buffer, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error copying response body: %s", err)
	}

	return buffer.Bytes(), nil
}

func main() {
	var (
		apiKey   = "cc410e1c00a12efa"
		location = "UT/Saint_George"
	)
	query := fmt.Sprintf("http://api.wunderground.com/api/%s/conditions/q/%s.json", apiKey, location)
	data, err := GetURL(query)
	if err != nil {
		log.Fatalf("error encountered getting URL: %s", err)
	}
	var response WeatherResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Fatalf("error encountered unmarshaling json: %s", err)
	}
	fmt.Printf("Weather: %s", response.CurrentObservation.Temperature)
}

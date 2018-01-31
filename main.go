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

// Location information
type Location struct {
	Full      string `json:"full"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Elevation string `json:"elevation"`
}

// Observation data
type Observation struct {
	DisplayLocation     Location `json:"display_location"`
	ObservationLocation Location `json:"observation_location"`
	StationID           string   `json:"station_id"`
	TimeString          string   `json:"observation_time"`
	Condition           string   `json:"weather"`
	ConditionURL        string   `json:"icon_url"`
	Temperature         string   `json:"temperature_string"`
	Fahrenheit          float64  `json:"temp_f"`
	Celsius             float64  `json:"temp_c"`
	Wind                string   `json:"wind_string"`
	WindDegrees         float64  `json:"wind_degrees"`
	WindSpeed           float64  `json:"wind_mph"`
	FeelsLike           string   `json:"feelslike_string"`
	Precipitation       string   `json:"precip_today_in"`
}

// Date data
type Date struct {
	Epoch     int    `json:"epoch"`
	Pretty    string `json:"pretty"`
	Day       int    `json:"day"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
	Hour      int    `json:"hour"`
	Minute    int    `json:"min"`
	Second    int    `json:"sec"`
	MonthName string `json:"monthname"`
	Weekday   string `json:"weekday"`
}

// Temperature data
type Temperature struct {
	Fahrenheit float64 `json:"fahrenheit"`
	Celsius    float64 `json:"celsius"`
}

// Depth data
type Depth struct {
}

// ForecastDay data
type ForecastDay struct {
	Date          Date        `json:"date"`
	Period        int         `json:"period"`
	High          Temperature `json:"high"`
	Low           Temperature `json:"low"`
	Conditions    string      `json:"conditions"`
	Precipitation Depth
}

// SimpleForecast data
type SimpleForecast struct {
	ForecastDays []ForecastDay `json:"forecastday"`
}

// Forecast data
type Forecast struct {
	SimpleForecast SimpleForecast `json:"simpleforecast"`
}

// WeatherResponse is the response data
type WeatherResponse struct {
	Response           MetaInfo    `json:"response"`
	CurrentObservation Observation `json:"current_observation"`
	Forecast           Forecast    `json:"forecast"`
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

	fmt.Printf("Current conditions for %s\n", response.CurrentObservation.DisplayLocation.Full)
	fmt.Printf("Station ID: %s\n", response.CurrentObservation.StationID)
	fmt.Printf("Temperature: %s feels like %s\n", response.CurrentObservation.Temperature, response.CurrentObservation.FeelsLike)
	fmt.Printf("Sky: %s\n", response.CurrentObservation.Condition)
}

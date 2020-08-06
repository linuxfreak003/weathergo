package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kelvins/geocoder"
)

var ArrowMap = map[string]string{
	"NW": "⬁",
	"NE": "⬀",
	"SE": "⬂",
	"SW": "⬃",
	"N":  "⇧",
	"S":  "⇩",
	"E":  "⇨",
	"W":  "⇦",
}

type Response struct {
	Properties Properties `json:"properties"`
}

type Elevation struct {
	Value    float64 `json:"value"`
	UnitCode string  `json:"unitCode"`
}

type Period struct {
	Number           int       `json:"number"`
	Name             string    `json:"name"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	IsDaytime        bool      `json:"isDaytime"`
	Temperature      int       `json:"temperature"`
	TemperatureUnit  string    `json:"temperatureUnit"`
	WindSpeed        string    `json:"windSpeed"`
	WindDirection    string    `json:"windDirection"`
	Icon             string    `json:"icon"`
	ShortForecast    string    `json:"shortForecast"`
	DetailedForecast string    `json:"detailedForecast"`
}

type Properties struct {
	Updated           string    `json:"updated"`
	Units             string    `json:"units"`
	ForecastGenerator string    `json:"forecastGenerator"`
	GeneratedAt       time.Time `json:"generatedAt"`
	UpdateTime        time.Time `json:"updateTime"`
	Elevation         Elevation `json:"elevation"`
	Periods           []Period  `json:"periods"`
}

func DownloadURL(url string) ([]byte, error) {
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

func usage(version string) {
	fmt.Printf("WeatherGo %s\n", version)
	fmt.Printf("Usage: %s [flags]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	// This example will need an API Key in your project
	// to set your API KEY follow directions as explained here:
	// https://developers.google.com/maps/documentation/geocoding/get-api-key
	version := "v0.3"
	var apiKey, zip string
	var city, state, country string
	var elevation, forecast, help bool

	flag.StringVar(&apiKey, "key", "", "`API key` from Open Weather Map")
	flag.StringVar(&country, "country", "United States", "Country")
	flag.StringVar(&city, "city", "", "City")
	flag.StringVar(&state, "state", "", "State")
	flag.StringVar(&zip, "zip", "", "`postal code`")
	flag.BoolVar(&elevation, "e", false, "Show elevation")
	flag.BoolVar(&forecast, "f", false, "Show forecast")
	flag.BoolVar(&help, "h", false, "Show help")

	flag.Parse()

	if help {
		usage(version)
	}

	geocoder.ApiKey = apiKey

	// Convert address to location (latitude, longitude)
	loc, err := geocoder.Geocoding(geocoder.Address{
		City:       city,
		State:      state,
		Country:    country,
		PostalCode: zip,
	})
	if err != nil {
		log.Fatalf("There was an error getting longitude/latitude: %v", err)
	}

	addresses, err := geocoder.GeocodingReverse(loc)
	address := addresses[0]

	url := fmt.Sprintf("https://api.weather.gov/points/%f,%f/forecast", loc.Latitude, loc.Longitude)
	bs, err := DownloadURL(url)
	if err != nil {
		log.Fatalf("There was an error getting forecast: %v", err)
	}

	var res Response
	err = json.Unmarshal(bs, &res)
	if err != nil {
		log.Panicf("Could not unmarshal json: %v", err)
	}

	fmt.Printf("Forecast for %s\n", address.FormattedAddress)
	if elevation {
		fmt.Printf("Elevation: %0.2f %s\n",
			res.Properties.Elevation.Value,
			res.Properties.Elevation.UnitCode,
		)
	}
	if !forecast {
		res.Properties.Periods = []Period{res.Properties.Periods[0]}
	}
	for _, period := range res.Properties.Periods {
		fmt.Printf("%s:\n", period.Name)
		fmt.Printf("\tTemperature: %d %s\n", period.Temperature, period.TemperatureUnit)
		fmt.Printf("\tSky: %s\n", period.ShortForecast)
		fmt.Printf("\tWind: %s to the %s %s\n", period.WindSpeed, period.WindDirection, ArrowMap[period.WindDirection])
	}
}

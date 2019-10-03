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
)

type WeatherResponse struct {
	Cod     string     `json:"cod"`
	Message float64    `json:"message"`
	Cnt     int        `json:"cnt"`
	City    City       `json:"city"`
	List    []Forecast `json:"list"`
}

type City struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Coord   Coord  `json:"coord"`
	Country string `json:"country"`
}

type Coord struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Forecast struct {
	Dt      int       `json:"dt"`
	Main    Main      `json:"main"`
	Weather []Weather `json:"weather"`
	Wind    Wind      `json:"wind"`
	Sys     Sys       `json:"sys"`
	DtTxt   string    `json:"dttxt"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  float64 `json:"pressure"`
	SeaLevel  float64 `json:"sea_level"`
	GrndLevel float64 `json:"grnd_level"`
	Humidity  float64 `json:"humidity"`
	TempKf    float64 `json:"temp_kf"`
}

type Weather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

type Sys struct {
	Pod string `json:"pod"`
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

func usage(version string) {
	fmt.Printf("WeatherGo %s\n", version)
	fmt.Printf("Usage: %s [flags]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	version := "v0.2"
	var apiKey, location, zip string
	var elevation, help bool

	flag.StringVar(&apiKey, "key", "adcd7639f2766b843be9c964cdccd3e2", "`API key` from Open Weather Map")
	flag.StringVar(&location, "loc", "Pheonix", "location in form of `state/city`")
	flag.StringVar(&zip, "zip", "84770", "`zip code`")
	// flag.BoolVar(&elevation, "e", false, "Show elevation")
	flag.BoolVar(&help, "h", false, "Show help")
	flag.Parse()

	if help {
		usage(version)
	}

	query := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?q=%s&APPID=%s&units=imperial", location, apiKey)
	data, err := GetURL(query)
	if err != nil {
		log.Fatalf("error encountered getting URL: %s", err)
	}

	var response WeatherResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Fatalf("error encountered unmarshalling json: %s", err)
	}

	fmt.Printf("Current conditions for %s\n", response.City.Name)
	forecast := response.List[response.Cnt-1]
	fmt.Printf("Temperature: %0.2f\n", forecast.Main.Temp)
	fmt.Printf("Sky: %s\n", forecast.Weather[0].Main)
	fmt.Printf("Wind: %0.2f\n", forecast.Wind.Speed)
	fmt.Printf("Pressure: %0.2f\n", forecast.Main.Pressure)
	// if elevation {
	// 	fmt.Printf("Elevation: %0.2f m\n", forecast.Main.GrndLevel)
	// }
}

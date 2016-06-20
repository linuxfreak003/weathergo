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
	"strconv"
)

func download(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Download failed: %s", err)
		return ""
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)

	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	s := string(buf.Bytes())
	return s
}

func parseCurrentConditions(jstring string) map[string]map[string]string {
	var iface interface{}
	err := json.Unmarshal([]byte(jstring), &iface)
	if err != nil {
		log.Fatal(err)
	}

	m := iface.(map[string]interface{})
	current := m["current_observation"].(map[string]interface{})

	today := make(map[string]map[string]string)
	info := make(map[string]string)
	for k, v := range current {
		switch vv := v.(type) {
		case float64:
			info[k] = strconv.FormatFloat(vv, 'e', 10, 64)
		case string:
			info[k] = vv
		case map[string]interface{}:
			newMap := make(map[string]string)
			for key, val := range vv {
				switch vval := val.(type) {
				case float64:
					newMap[key] = strconv.FormatFloat(vval, 'e', 10, 64)
				case string:
					newMap[key] = vval
				}
			}
			today[k] = newMap
		default:
			log.Printf("Encountered unknown type")
		}
	}
	today["main"] = info
	return today
}

func parseForecast(jstring string) []map[string]interface{} {
	var iface interface{}
	err := json.Unmarshal([]byte(jstring), &iface)
	if err != nil {
		log.Fatal(err)
	}

	m := iface.(map[string]interface{})
	m = m["forecast"].(map[string]interface{})
	m = m["simpleforecast"].(map[string]interface{})
	all := m["forecastday"].([]interface{})

	var days []map[string]interface{}
	for _, item := range all {
		days = append(days, item.(map[string]interface{}))
	}

	return days
}

func main() {
	var zip, apikey string
	//var days int
	var humidity, help, elevation bool

	//This api key is specific to me, if you want to use this application please use your own.
	//API key is free, simply go to https://www.wunderground.com/weather/api and make an account.
	flag.StringVar(&apikey, "key", "92d518fe1c24dc58", "API key from Weather Underground")
	flag.StringVar(&zip, "zip", "84770", "Zipcode")
	//flag.StringVar(&state, "state", "UT", "State")
	//flag.StringVar(&city, "city", "SAINT_GEORGE", "City")
	//flag.IntVar(&days, "days", 1, "Days to forecast")
	flag.BoolVar(&humidity, "h", false, "Humidity")
	flag.BoolVar(&elevation, "e", false, "Elevation")
	flag.BoolVar(&help, "help", false, "Help information")

	flag.Parse()

	fmt.Printf("WeatherGo by Jared Chapman\n\n")

	if help {
		fmt.Printf("Usage: %s [options...]\n", os.Args[0])
		fmt.Printf("Options:\n")
		fmt.Println("-zip\tZipcode")
		fmt.Println("-key\tAPI key to use")
		fmt.Println("-days\tNumber of days to forecast(not yet implemented)")
		fmt.Println("-e\tShow Elevation")
		fmt.Println("-h\tShow Humidity")
		fmt.Println("-f\tShow Forecast")

		os.Exit(0)
	}

	jsonString := download("http://api.wunderground.com/api/" + apikey + "/conditions/forecast/q/" + zip + ".json")
	parsedInfo := parseCurrentConditions(jsonString)
	parsedForecast := parseForecast(jsonString)

	//fmt.Println(parsedForecast)
	for _, day := range parsedForecast {
		d := day["date"].(map[string]interface{})
		fmt.Println(d["pretty"])
	}

	fmt.Printf("Location: %s\n", parsedInfo["display_location"]["full"])
	fmt.Printf("Temperature: %s\n", parsedInfo["main"]["temperature_string"])
	fmt.Printf("%s\n", parsedInfo["main"]["observation_time"])
	if elevation {
		fmt.Printf("Elevation: %s\n", parsedInfo["observation_location"]["elevation"])
	}
	fmt.Printf("Sky: %s\n", parsedInfo["main"]["weather"])
	if humidity {
		fmt.Printf("Humidity: %s\n", parsedInfo["main"]["relative_humidity"])
	}
}

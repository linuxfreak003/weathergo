package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
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

func parseJsonString(s string) map[string]map[string]string {
	var i interface{}
	err := json.Unmarshal([]byte(s), &i)
	if err != nil {
		log.Fatal(err)
	}

	m := i.(map[string]interface{})
	current := m["current_observation"].(map[string]interface{})

	stuff := make(map[string]map[string]string)
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
			stuff[k] = newMap
		default:
			log.Printf("Encountered unknown type")
		}
	}
	stuff["main"] = info
	return stuff
}

func main() {
	var zip, state, city string
	var days int
	var humidity, help bool
	flag.StringVar(&zip, "zip", "84770", "Zipcode")
	flag.StringVar(&state, "state", "UT", "State")
	flag.StringVar(&city, "city", "SAINT_GEORGE", "City")
	flag.IntVar(&days, "days", 1, "Days to forecast")
	flag.BoolVar(&humidity, "h", false, "Humidity")
	flag.BoolVar(&help, "help", false, "Show help information")
	flag.Parse()

	fmt.Println("Weather app")
	s := download("http://api.wunderground.com/api/92d518fe1c24dc58/conditions/q/" + state + "/" + city + ".json")

	parsedInfo := parseJsonString(s)

	fmt.Printf("Weather report for %s\n", parsedInfo["display_location"]["full"])
	fmt.Printf("at %s\n", parsedInfo["observation_location"]["elevation"])
	fmt.Printf("%s\n", parsedInfo["main"]["observation_time"])
	fmt.Printf("Sky: %s\n", parsedInfo["main"]["weather"])
	fmt.Printf("Temperature: %s\n", parsedInfo["main"]["temperature_string"])
	fmt.Printf("Humidity: %s\n", parsedInfo["main"]["relative_humidity"])
}

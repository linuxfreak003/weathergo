package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func parseForecast(jstring string) []map[string]map[string]string {
	var iface interface{}
	err := json.Unmarshal([]byte(jstring), &iface)
	if err != nil {
		log.Fatal(err)
	}

	m := iface.(map[string]interface{})
	m = m["forecast"].(map[string]interface{})
	m = m["simpleforecast"].(map[string]interface{})
	all := m["forecastday"].([]interface{})

	var days []map[string]map[string]string
	for _, item := range all {
		current := item.(map[string]interface{})
		info := make(map[string]string)
		forecast := make(map[string]map[string]string)
		for k, v := range current {
			switch vv := v.(type) {
			case float64:
				if vv-float64(int(vv)) == 0 {
					info[k] = strconv.FormatFloat(vv, 'f', 0, 64)
				} else {
					info[k] = strconv.FormatFloat(vv, 'f', 2, 64)
				}
			case string:
				info[k] = vv
			case map[string]interface{}:
				newMap := make(map[string]string)
				for key, val := range vv {
					switch vval := val.(type) {
					case float64:
						if vval-float64(int(vval)) == 0 {
							newMap[key] = strconv.FormatFloat(vval, 'f', 0, 64)
						} else {
							newMap[key] = strconv.FormatFloat(vval, 'f', 2, 64)
						}
					case string:
						newMap[key] = vval
					}
				}
				forecast[k] = newMap
			default:
				log.Printf("Encountered unknown type")
			}
		}
		forecast["main"] = info
		days = append(days, forecast)
	}
	return days
}

func main() {
	var loc, apikey, config string
	var days int
	var humidity, help, elevation, forecast bool

	//This api key is specific to me, if you want to use this application please use your own.
	//API key is free, simply go to https://www.wunderground.com/weather/api and make an account.
	flag.StringVar(&apikey, "key", "", "API key from Weather Underground")
	flag.StringVar(&loc, "loc", "CA/San_Francisco", "Location")
	flag.StringVar(&config, "c", "", "Config file")
	flag.IntVar(&days, "days", 11, "Days to forecast")
	flag.BoolVar(&humidity, "h", false, "Humidity")
	flag.BoolVar(&elevation, "e", false, "Elevation")
	flag.BoolVar(&help, "help", false, "Help information")
	flag.BoolVar(&forecast, "f", false, "Forecast")

	flag.Parse()

	if help {
		fmt.Println("WeatherGo by Jared Chapman v0.1")

		fmt.Printf("Usage: %s [Options]\n", os.Args[0])
		fmt.Println("CONFIG:")
		fmt.Println("  -c=<filename>\tConfig file to use for location parameters")
		fmt.Println("  \t\tNote that parameters from file will be")
		fmt.Println("  \t\toverwritten by any flags")
		fmt.Println("  -key=<key>\tAPI key to use")
		fmt.Println("LOCATION:")
		fmt.Println("  -loc=<zip>\tZipcode")
		fmt.Println("  -loc=<SS/CITY>\tState/City")
		fmt.Println("INFORMATION:")
		fmt.Println("  -days\tNumber of days to forecast(limit of 10)")
		fmt.Println("  -e\tShow Elevation")
		fmt.Println("  -h\tShow Humidity")
		fmt.Println("  -f\tShow Forecast(4 day)")
		fmt.Println("EXAMPLES:")
		fmt.Printf("  %s -key=<api_key>\n", os.Args[0])
		fmt.Printf("  %s -loc=80432 -f -e -h\n", os.Args[0])
		fmt.Printf("  %s -loc=CA/San_Francisco -f -e -h\n", os.Args[0])
		fmt.Printf("  %s -c=config_file\n", os.Args[0])

		os.Exit(0)
	}

	if config != "" {
		dat, err := ioutil.ReadFile(config)
		if err != nil {
			log.Fatal(err)
		}
		lines := strings.Split(string(dat), "\n")
		for _, line := range lines {
			params := strings.Split(line, "=")
			if len(params) > 1 {
				switch params[0] {
				case "key":
					if apikey == "" {
						apikey = params[1]
					}
				case "loc":
					if loc == "CA/San_Francisco" {
						loc = params[1]
					}
				case "days":
					if days == 11 {
						days, err = strconv.Atoi(params[1])
						if err != nil {
							log.Fatalf("Invalid parameter for days: %s", params[1])
						}
					}
				case "h":
					if !humidity {
						if params[1] == "1" {
							humidity = true
						}
					}
				case "e":
					if !elevation {
						if params[1] == "1" {
							elevation = true
						}
					}
				case "f":
					if !forecast {
						if params[1] == "1" {
							forecast = true
						}
					}
				default:
					log.Fatalf("Invalid parameter in %s: %s", config, params[0])
				}

			}
		}

		//os.Exit(0)
	}

	var parsedForecast []map[string]map[string]string
	jsonString := download("http://api.wunderground.com/api/" + apikey + "/conditions/forecast10day/q/" + loc + ".json")
	parsedInfo := parseCurrentConditions(jsonString)
	if forecast {
		parsedForecast = parseForecast(jsonString)
	}

	//fmt.Println(parsedForecast)
	fmt.Printf("Location: %s\n", parsedInfo["display_location"]["full"])
	fmt.Printf("Temperature: %s\n", parsedInfo["main"]["temperature_string"])
	if elevation {
		fmt.Printf("Elevation: %s\n", parsedInfo["observation_location"]["elevation"])
	}
	fmt.Printf("Sky: %s\n", parsedInfo["main"]["weather"])
	if humidity {
		fmt.Printf("Humidity: %s\n", parsedInfo["main"]["relative_humidity"])
	}
	fmt.Printf("%s\n", parsedInfo["main"]["observation_time"])

	if forecast {
		fmt.Println("Forecast:")
		for i, day := range parsedForecast {
			if i < days {
				fmt.Printf("  %s %s %s, %s", day["date"]["weekday_short"], day["date"]["monthname_short"], day["date"]["day"], day["date"]["year"])
				fmt.Printf("\t%s with a high of %s °F (%s °C)", day["main"]["conditions"], day["high"]["fahrenheit"], day["high"]["celsius"])
				fmt.Printf(" and low of %s °F (%s °C)", day["low"]["fahrenheit"], day["low"]["celsius"])
				fmt.Printf(" and %s\" rainfall\n", day["qpf_allday"]["in"])
			}
		}
	}

}

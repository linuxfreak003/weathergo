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
	var days, rain int
	var humidity, help, elevation, forecast, vers, pressure, wind bool
	version := "v0.1"

	//This api key is specific to me, if you want to use this application please use your own.
	//API key is free, simply go to https://www.wunderground.com/weather/api and make an account.
	flag.StringVar(&apikey, "key", "", "API `key` from Weather Underground")
	flag.StringVar(&loc, "loc", "autoip", "Location: zipcode|state/city")
	flag.StringVar(&config, "c", "", "Config `file`")
	flag.BoolVar(&humidity, "h", true, "Show Humidity")
	flag.BoolVar(&elevation, "e", false, "Show Elevation")
	flag.BoolVar(&pressure, "p", false, "Show Pressure")
	flag.BoolVar(&wind, "w", false, "Show Wind")
	flag.BoolVar(&help, "help", false, "Show Help Information")
	flag.BoolVar(&vers, "v", false, "Show Version")
	flag.BoolVar(&forecast, "f", false, "Show Forecast")
	flag.IntVar(&days, "days", 11, "Days to forecast")
	flag.IntVar(&rain, "r", 0, "Days to total predicted rainfall")

	flag.Parse()

	if vers {
		fmt.Printf("WeatherGo %s\n", version)
		os.Exit(0)
	}

	if help {
		fmt.Printf("WeatherGo %s\n", version)

		fmt.Printf("Usage: %s [Options]\n", os.Args[0])
		fmt.Println("CONFIG:")
		fmt.Println("  -c=<filename>\tConfig file to use for location parameters")
		fmt.Println("  \t\tNote that parameters from file will be")
		fmt.Println("  \t\toverwritten by any flags")
		fmt.Println("  -key=<key>\tAPI key to use")
		fmt.Println("LOCATION:")
		fmt.Println("  -loc=<zip>\tZipcode")
		fmt.Println("  -loc=<ST/CITY>\tState/City")
		fmt.Println("INFORMATION:")
		fmt.Println("  -days\tNumber of days to forecast(limit of 10)")
		fmt.Println("  -e\tShow Elevation")
		fmt.Println("  -h\tShow Humidity")
		fmt.Println("  -f\tShow Forecast(4 day)")
		fmt.Println("  -r\tShow total predicted rainfall for next # days.")
		fmt.Println("EXAMPLES:")
		fmt.Printf("  %s -key=<api_key>\n", os.Args[0])
		fmt.Printf("  %s -loc=80432 -f -e -h\n", os.Args[0])
		fmt.Printf("  %s -loc=CA/San_Francisco -f -e -h\n", os.Args[0])
		fmt.Printf("  %s -c=config_file\n", os.Args[0])

		os.Exit(0)
	}

	if config == "" {
		homeconfig := os.Getenv("HOME") + "/.weathergo"
		if _, err := os.Stat(homeconfig); os.IsNotExist(err) {

		} else {
			config = homeconfig
		}
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
				params[0] = strings.TrimSpace(params[0])
				params[1] = strings.TrimSpace(params[1])
				switch params[0] {
				case "key":
					if apikey == "" {
						apikey = params[1]
					}
				case "loc":
					if loc == "autoip" {
						loc = params[1]
					}
				case "days":
					if days == 11 {
						days, err = strconv.Atoi(params[1])
						if err != nil {
							log.Fatalf("Invalid parameter for days: %s", params[1])
						}
					}
				case "r":
					if rain == 0 {
						rain, err = strconv.Atoi(params[1])
						if err != nil {
							log.Fatalf("Invalid parameter for rain: %s", params[1])
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
				case "p":
					if !pressure {
						if params[1] == "1" {
							pressure = true
						}
					}
				case "w":
					if !wind {
						if params[1] == "1" {
							wind = true
						}
					}
				default:
					log.Fatalf("Invalid parameter in %s: %s", config, params[0])
				}

			}
		}
	}

	if apikey == "" {
		log.Fatal("Error: No API key set!")
	}

	var parsedForecast []map[string]map[string]string
	jsonString := download("http://api.wunderground.com/api/" + apikey + "/conditions/forecast10day/q/" + loc + ".json")
	parsedInfo := parseCurrentConditions(jsonString)
	if forecast {
		parsedForecast = parseForecast(jsonString)
	}

	//fmt.Println(parsedForecast)
	fmt.Printf("Current conditions for %s\n", parsedInfo["display_location"]["full"])
	fmt.Printf("Station ID: %s\n", parsedInfo["main"]["station_id"])
	fmt.Printf("Temperature: %s\n", parsedInfo["main"]["temperature_string"])
	fmt.Printf("Sky: %s\n", parsedInfo["main"]["weather"])
	if humidity {
		fmt.Printf("Humidity: %s\n", parsedInfo["main"]["relative_humidity"])
	}
	if elevation {
		fmt.Printf("Elevation: %s\n", parsedInfo["observation_location"]["elevation"])
	}
	if pressure {
		fmt.Printf("Pressure: %sin\n", parsedInfo["main"]["pressure_in"])
	}
	if wind {
		fmt.Printf("Wind: %s\n", parsedInfo["main"]["wind_string"])
	}
	fmt.Printf("%s\n", parsedInfo["main"]["observation_time"])

	if forecast || rain > 0 {
		fmt.Println("Forecast:")
		rainfall := 0.0
		for i, day := range parsedForecast {
			if i < rain {
				r := day["qpf_allday"]["in"]
				f, err := strconv.ParseFloat(r, 64)
				if err != nil {
					log.Fatal("Invalid value for rainfall:", r)
				}
				rainfall += f
			}
			if i < days && forecast {
				fmt.Printf("  %s %s %s, %s\n", day["date"]["weekday_short"], day["date"]["monthname_short"], day["date"]["day"], day["date"]["year"])
				fmt.Printf("\tConditions:\t%s\n\tHigh:\t%s °F (%s °C)", day["main"]["conditions"], day["high"]["fahrenheit"], day["high"]["celsius"])
				fmt.Printf("\n\tLow:\t%s °F (%s °C)", day["low"]["fahrenheit"], day["low"]["celsius"])
				fmt.Printf("\n\tRainfall:\t%s\"", day["qpf_allday"]["in"])
				fmt.Printf("\n\tSnowfall:\t%s\"", day["snow_allday"]["in"])
				fmt.Printf("\n\tHumidity:\t%s%%\n", day["main"]["avehumidity"])
			}
		}
		if rain > 0 {
			fmt.Printf("%0.2f total inches of rainfall predicted in next %d days.\n", rainfall, rain)
		}
	}
}

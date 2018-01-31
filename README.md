# WeatherGo

A simple command-line weather tool written in Go

## Getting started

You must have golang installed and working on your system with a valid
$GOPATH.
To use this tool you first need to obtain an api key from
[Weather Underground](https://www.wunderground.com/weather/api)
The key is free for developers. You are only limited to 10 api calls per minute,
and 500 per day.

## Installation

To install run

`git clone https://github.com/linuxfreak003/weathergo.git`

then `cd weathergo; go install`

(For `go install` to work your $GOBIN variable must be set, optionally
  you can run `go build` instead)

## Usage: `weathergo [-c config_file | -key <apikey>] [Options...]`

```bash
CONFIG:
  -c <filename> Config file to use for location parameters
                Note that parameters from file will be
                overwritten by any command-line flags
  -key <key>    API key to use
LOCATION:
  -loc <zip>      Zipcode
  -loc <ST/CITY>  State/City
INFORMATION:
  -days <#> Number of days to forecast(limit of 10)
  -e        Show Elevation
  -h        Show Humidity
  -f        Show Forecast(10 day default)
  -v        Show Version information
  -r        Show Rainfall predictions
EXAMPLES:
  weathergo -key <api_key>
  weathergo -loc 80432 -f -e -h -key <api_key>
  weathergo -loc CA/San_Francisco -f -e -h -c <config_file>
  weathergo -c <config_file>
```

### Config

To set variables in config file on each line list `<flag>=<value>` (for boolean variables value is 1)

```txt
EXAMPLE CONFIG:
key=<api_key>
loc=CA/San_Francisco
days=7
f=1
h=1
```

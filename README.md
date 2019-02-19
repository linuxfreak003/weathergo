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

### Using `go get`

`go get -u github.com/linuxfreak003/weathergo`

## Usage

Command:
`weathergo -key <apikey> [Options...]`

```bash
CONFIG:
  -key <key>    API key to use
LOCATION:
  -loc <zip>      Zipcode
  -loc <ST/CITY>  State/City
INFORMATION:
  -e        Show Elevation
  -h        Show Humidity #Depricated for now
  -v        Show Version information #Depricated for now
  -r        Show Rainfall predictions #Depricated for now
EXAMPLES:
  weathergo -key <api_key>
  weathergo -loc 80432 -key <api_key>
  weathergo -loc CA/San_Francisco -e
```
### LICENSE

[![LICENSE](https://img.shields.io/pypi/l/Django.svg)](LICENSE)

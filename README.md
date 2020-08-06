# WeatherGo

A simple command-line weather tool written in Go

## Getting started

You must have golang installed and working on your system with a valid $GOPATH.
To use this tool you first need to obtain a Google geocode api key from
[here](https://developers.google.com/maps/documentation/geocoding/get-api-key)

You will probably have to start a free trial if you do not already have an account

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
  -city <city>        Zipcode
  -state <state>      Zipcode
  -country <country>  Country (defaults to "United States")
  -zip <ZIP>          Postal Code
INFORMATION:
  -e        Show Elevation
  -h        Show Humidity #Depricated for now
  -v        Show Version information #Depricated for now
  -r        Show Rainfall predictions #Depricated for now
EXAMPLES:
  weathergo -key <api_key>
  weathergo -zip 80432 -key <api_key>
  weathergo -state California -city "San Francisco" -e
```
### LICENSE

[![LICENSE](https://img.shields.io/pypi/l/Django.svg)](LICENSE)

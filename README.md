# WeatherGo
A simple command-line weather tool written in Go

###Getting started
You must have golang installed and working on your system with a valid
$GOPATH.
To use this tool you first need to obtain an api key from
[Weather Underground](https://www.wunderground.com/weather/api)
It is free for developers. You are just limited to 10 api calls per minute,
and 500 per day.

###Installation
To install run
`git clone https://github.com/linuxfreak003/weathergo.git`
then `cd weathergo; go install`
(For `go install` to work your $GOBIN variable must be set)

####Usage: `weathergo [-c=config_file | -key=<apikey>] [Options...]`

```
CONFIG:
  -c=<filename> Config file to use for location parameters
                Note that parameters from file will be
                overwritten by any flags
  -key=<key>    API key to use
LOCATION:
  -loc=<zip>      Zipcode
  -loc=<SS/CITY>  State/City
INFORMATION:
  -days=<#> Number of days to forecast(limit of 10)
  -e        Show Elevation
  -h        Show Humidity
  -f        Show Forecast(10 day default)
  -v        Show Version information
EXAMPLES:
  weathergo -key=<api_key>
  weathergo -loc=80432 -f -e -h -key=<api_key>
  weathergo -loc=CA/San_Francisco -f -e -h -c=<config_file>
  weathergo -c=<config_file>
```

####Example Config:
```
key=<api_key>
loc=autoip
days=7
f=1
```

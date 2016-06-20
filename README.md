# weathergo
A Simple command line weather tool written in Go

##API key from Weather Underground
To use this tool you first need to obtain an api key from
[Weather Underground](https://www.wunderground.com/weather/api)
It is free for developers. You are just limited to 10 api calls per minute,
and 500 per day.

###Usage: `weathergo [OPTION...]`

###Options:
	```
	-h		Show Humidity
  -zip  Zipcode to use
  -key  Weather Underground API key to use
  -days Number of days to forecast(not yet implemented)
  -e    Show Elevation
  -h    Show Humidity
	```

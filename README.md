# weathergo
A Simple command line weather tool written in Go

##Getting started
You must have golang installed and working on your system with a valid
GOPATH.
To use this tool you first need to obtain an api key from
[Weather Underground](https://www.wunderground.com/weather/api)
It is free for developers. You are just limited to 10 api calls per minute,
and 500 per day.

###Installation
To install run `git clone linuxfreak003/weathergo`
then `cd weathergo; go install`

If you do choose to install weathergo, understand that is it still heavily
under development (I dont even have version numbers yet).

###Usage: `weathergo [Options...]`

####Options:
```
-zip  Zipcode to use
-key  Weather Underground API key to use
-days Number of days to forecast (up to 10)
-e    Show Elevation
-h    Show Humidity
```

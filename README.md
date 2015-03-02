# geolookup - CLI in Go

## Overview 
Command line tool for looking up geo location information of IP addresses.

Requires MaxMind [GeoLite2 (free)](http://dev.maxmind.com/geoip/geoip2/geolite2/)
or [GeoIP2](http://www.maxmind.com/en/geolocation_landing) databases.

Based on [geoip2](https://github.com/oschwald/geoip2-golang) library by oschwald

## Installation

```
go install github.com/ctrlrsf/geolookup
```

## Sample usage

```
$ geolookup --countrydb ~/GeoLite2-Country.mmdb --country --ipv4Address 8.8.8.8
United States
```

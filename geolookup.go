package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/codegangsta/cli"
	"github.com/oschwald/geoip2-golang"
)

func main() {
	app := cli.NewApp()
	app.Name = "geolookup"
	app.Author = "Rene Fragoso"
	app.Email = "ctrlrsf@gmail.com"
	app.Usage = "Quick tool for querying Maxmind Geo IP database file"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "country",
			Usage: "Query country field",
		},
		cli.StringFlag{
			Name: "ipv4Address",
		},
		cli.StringFlag{
			Name:  "countrydb",
			Value: "GeoLite2-Country.mmdb",
			Usage: "Maxmind database file",
		},
	}

	app.Action = func(c *cli.Context) {
		ipv4Address := c.String("ipv4Address")
		countryDb := c.String("countrydb")
		country := c.Bool("country")

		if country {
			queryGeoIp(countryDb, ipv4Address)
		}

	}

	app.Run(os.Args)
}

func queryGeoIp(countryDb string, ipv4Address string) {
	db, err := geoip2.Open(countryDb)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(ipv4Address)

	record, err := db.Country(ip)
	if err != nil {
		log.Fatal(err)
	}

	countryName := record.Country.Names["en"]
	if countryName == "" {
		fmt.Printf("Country not found for IP: %s\n", ipv4Address)
	} else {
		fmt.Printf("%s\n", record.Country.Names["en"])
	}
}

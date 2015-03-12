package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/oschwald/geoip2-golang"
)

const countryDbURL = "http://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.mmdb.gz"
const countryDbOutputFile = "GeoLite2-Country.mmdb"

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
		cli.BoolFlag{
			Name:  "update",
			Usage: "Update Maxmind GeoLite2 database",
		},
		cli.StringFlag{
			Name: "ipv4Address",
		},
		cli.StringFlag{
			Name:  "countrydb",
			Value: countryDbOutputFile,
			Usage: "Maxmind database file",
		},
	}

	app.Action = func(c *cli.Context) {
		ipv4Address := c.String("ipv4Address")
		countryDb := c.String("countrydb")
		country := c.Bool("country")
		update := c.Bool("update")

		if update {
			updateGeoLite2DB()
		}

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

func updateGeoLite2DB() {
	resp, err := http.Get(countryDbURL)
	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("Error with HTTP request: %v\n", err)
	}

	gzipReader, err := gzip.NewReader(resp.Body)
	defer gzipReader.Close()
	if err != nil {
		fmt.Printf("Error decompressing response: %v", err)
	}

	dbBytes, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		fmt.Printf("Error reading decompressed bytes: %v", err)
	}

	err = ioutil.WriteFile(countryDbOutputFile, dbBytes, 0744)

	if err != nil {
		fmt.Printf("Error writing file: %v", err)
	}
}

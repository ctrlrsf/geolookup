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
const version = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "geolookup"
	app.Author = "Rene Fragoso"
	app.Email = "ctrlrsf@gmail.com"
	app.Version = version
	app.Usage = "Quick tool for querying Maxmind Geo IP database file"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "country",
			Usage: "Query country field",
		},
		cli.BoolFlag{
			Name:  "update",
			Usage: "Download or update Maxmind GeoLite2 database",
		},
		cli.StringFlag{
			Name:  "ipv4",
			Usage: "Specify IPv4 address to lookup",
		},
		cli.StringFlag{
			Name:  "countrydb",
			Value: countryDbOutputFile,
			Usage: "Maxmind GeoLite2 country DB file",
		},
	}

	app.Action = func(c *cli.Context) {
		ipv4 := c.String("ipv4")
		countryDb := c.String("countrydb")
		country := c.Bool("country")
		update := c.Bool("update")

		if update {
			updateGeoLite2DB()
		}

		if country {
			queryGeoIp(countryDb, ipv4)
		}

	}

	app.Run(os.Args)
}

func queryGeoIp(countryDb string, ipv4 string) {
	db, err := geoip2.Open(countryDb)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(ipv4)

	record, err := db.Country(ip)
	if err != nil {
		log.Fatal(err)
	}

	countryName := record.Country.Names["en"]
	if countryName == "" {
		fmt.Printf("Country not found for IP: %s\n", ipv4)
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

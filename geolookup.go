package main

import (
	"flag"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"os"
)

func main() {
	countryDb := flag.String("countrydb", "GeoLite2-Country.mmdb",
		"Binary MaxMind GeoLite DB file")

	countryFlag := flag.Bool("country", false, "REQUIRED: Look up country for IP")

	ipv4Address := flag.String("ipv4", "", "REQUIRED: IP address")

	flag.Parse()

	if *countryFlag == false {
		fmt.Fprintf(os.Stderr, "Error: country flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *ipv4Address == "" {
		fmt.Fprintf(os.Stderr, "Error: ipv4Address flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	db, err := geoip2.Open(*countryDb)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(*ipv4Address)

	record, err := db.Country(ip)
	if err != nil {
		log.Fatal(err)
	}

	countryName := record.Country.Names["en"]
	if countryName == "" {
		fmt.Printf("Country not found for IP: %s\n", *ipv4Address)
	} else {
		fmt.Printf("%s\n", record.Country.Names["en"])
	}
}

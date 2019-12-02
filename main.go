package main

import (
	"flag"
	"log"
	"os"

	"github.com/ae33/port-knock/config"
	"github.com/ae33/port-knock/knock"
)

func main() {
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)

	configPath := flagSet.String("config-path", "/etc/port-knock/port-knock.yml", "Location of the port-knock Systemd config file.")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("error parsing command-line flags: '%v'", err)
	}

	conf := config.ParseConfig(*configPath)

	err = knock.UdpKnocks(conf)
	if err != nil {
		log.Fatalf("error executing udp port knocking: '%v'", err)
	}

	os.Exit(0)
}

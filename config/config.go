package config

import (
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Ports     []uint16       `yaml:"ports"`
	Host      string         `yaml:"host"`
	QuitAfter *time.Duration `yaml:"quit_after"`
	WaitSleep *time.Duration `yaml:"wait_sleep"`
}

func ParseConfig(configPath string) Config {
	log.Printf("reading config from '%s'", configPath)

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("error reading configuration file: '%v'", err)
	}

	var conf Config
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		log.Fatalf("error unmarshalling configuration file to yaml: '%v'", err)
	}

	// Check for mandatory fields
	if len(conf.Ports) <= 0 {
		log.Fatalf("no remote ports configured")
	}

	if conf.Host == "" {
		log.Fatalf("no remote host configured")
	}

	// Set defaults
	if conf.QuitAfter == nil {
		conf.QuitAfter = func() *time.Duration { d := 10 * time.Second; return &d }()
	}

	if conf.WaitSleep == nil {
		conf.WaitSleep = func() *time.Duration { d := 100 * time.Microsecond; return &d }()
	}

	return conf
}

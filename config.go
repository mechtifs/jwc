package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type ClientConfig struct {
	Mode     int  `toml:"mode"`
	Delay    int  `toml:"delay"`
	Parallel int  `toml:"parallel"`
	Keep     bool `toml:"keep"`
	Verbose  bool `toml:"verbose"`
}

type SummaryConfig struct {
	Enabled  bool `toml:"enabled"`
	Interval int  `toml:"interval"`
}

type SessionConfig struct {
	BaseURI    string   `toml:"base_uri"`
	UserAgent  string   `toml:"user_agent"`
	Credential string   `toml:"credential"`
	Targets    []string `toml:"targets"`
}

type Config struct {
	Client  ClientConfig  `toml:"client"`
	Summary SummaryConfig `toml:"summary"`
	Session SessionConfig `toml:"session"`
}

func parseConfig() {
	if len(os.Args) == 1 {
		log.Fatal("[FATAL] Please specify the path of config file.")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("[FATAL] Failed to open file:", err)
	}
	defer file.Close()
	toml.NewDecoder(file).Decode(&config)
	if config.Client.Verbose {
		log.Println("[DEBUG]", config)
	}
}

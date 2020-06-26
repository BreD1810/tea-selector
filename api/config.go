package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// A Config represents the config file data.
type Config struct {
	Server struct {
		Port            string `yaml:"port"`
		RegisterEnabled bool   `yaml:"registerenabled"`
		SigningKey      string `yaml:"signingkey"`
	} `yaml:"server"`
	Database struct {
		Location string   `yaml:"location"`
		TeaTypes []string `yaml:"teaTypes"`
		Owners   []string `yaml:"owners"`
	} `yaml:"database"`
}

func getConfig() Config {
	f, err := os.Open("config.yml")

	if err != nil {
		log.Fatal("Error opening config.yml")
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal("Error parsing config.yml")
	}

	logConfig(cfg)
	return cfg
}

func logConfig(cfg Config) {
	if cfg.Server.SigningKey == "" {
		log.Fatal("Error: no signing key in config")
	} else {
		log.Println("Signing key set successfully")
	}

	if cfg.Server.Port == "" {
		log.Fatal("Error: no port specified")
	} else {
		log.Printf("Port: %v\n", cfg.Server.Port)
	}

	if cfg.Server.RegisterEnabled {
		log.Println(`POST /register endpoint enabled`)
	} else {
		log.Println(`POST /register endpoint disabled`)
	}
	log.Printf("Database Location: %v\n", cfg.Database.Location)
	log.Printf("Tea types: %q\n", cfg.Database.TeaTypes)
	log.Printf("Owners: %q\n", cfg.Database.Owners)
}

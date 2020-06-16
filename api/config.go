package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// A Config represents the config file data.
type Config struct {
	Server struct {
		Port string `yaml:"port"`
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
	log.Printf("Port: %v\n", cfg.Server.Port)
	log.Printf("Database Location: %v\n", cfg.Database.Location)
	log.Printf("Tea types: %q\n", cfg.Database.TeaTypes)
}

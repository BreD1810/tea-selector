package main

import (
	"log"
	"os"
	"gopkg.in/yaml.v2"
)

type config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func getConfig() config {
	f, err := os.Open("config.yml")
	
	if err != nil {
		log.Fatal("Error opening config.yml")	
	}
	defer f.Close()

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal("Error parsing config.yml")
	}

	logConfig(cfg)
	return cfg
}

func logConfig(cfg config) {
	log.Printf("Port: %v\n", cfg.Server.Port)
}

// Package config provides a way to load configuration from a file or environment variables
package config

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	Port             string `json:"port"`
	DefaultLanguage  string `json:"default_language"`
	LegacyEndpoint   string `json:"legacy_endpoint"`
	DatabaseType     string `json:"database_type"`
	DatabaseURL      string `json:"database_url"`
	DatabasePort     string `json:"database_port"`
	DatabasePassword string `json:"database_password"`
}

var defaultConfiguration = Configuration{
	Port:            ":8080",
	DefaultLanguage: "english",
}

func (c *Configuration) LoadFromEnv() {
	if lang := os.Getenv("DEFAULT_LANGUAGE"); lang != "" {
		c.DefaultLanguage = lang
	}
	if port := os.Getenv("PORT"); port != "" {
		c.Port = port
	}
}

func (c *Configuration) ParsePort() {
	if c.Port[0] != ':' {
		c.Port = ":" + c.Port
	}
	if _, err := strconv.Atoi(string(c.Port[1:])); err != nil {
		c.Port = defaultConfiguration.Port
	}

}

func (c *Configuration) LoadFromJSON(path string) error {
	log.Printf("loading configuration from file: %s\n", path)
	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("unable to load file: %s\n", err.Error())
		return errors.New("unable to load configuration")
	}
	if err := json.Unmarshal(b, c); err != nil {
		log.Printf("unable to parse file: %s\n", err.Error())
		return errors.New("unable to load configuration")
	}
	return nil
}

func LoadConfiguration() Configuration {
	cfgfilePtr := flag.String("config_file", "", "load configurations from a file")
	portPtr := flag.String("port", "", "set port")

	flag.Parse()

	cfg := defaultConfiguration

	if cfgfilePtr != nil && *cfgfilePtr != "" {
		if err := cfg.LoadFromJSON(*cfgfilePtr); err != nil {
			log.Printf("unable to load configuration from json: %s, using default values", *cfgfilePtr)
		}
	}

	cfg.LoadFromEnv()

	if portPtr != nil && *portPtr != "" {
		cfg.Port = *portPtr
	}

	cfg.ParsePort()
	return cfg
}

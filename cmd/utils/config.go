package utils

import (
	"encoding/json"
	"log"
	"os"
)

type Dfuse struct {
	ApiKey  string `json:"api_key"`
	Manager string `json:"manager"`
}

type CryptoLions struct {
	Server string `json:"server"`
}

type Server struct {
	ServerType  string      `json:"serverType"` // dfuse cryptolions
	Dfuse       Dfuse       `json:"dfuse"`
	CryptoLions CryptoLions `json:"cryptolions"`
}
type MongoConf struct {
	URI         string `json:"uri"`
	Database    string `json:"database"`
	MaxPoolSize uint64 `json:"maxPoolSize"`
	MinPoolSize uint64 `json:"minPoolSize"`
	AppName     string `json:"appName"`
}

type Config struct {
	Network  string    `json:"network"`
	Manager  string    `json:"manager"`
	Server   Server    `json:"server"`
	Mongo    MongoConf `json:"mongo"`
	Interval string    `json:"interval"`
	Listen   string    `json:"listen"`
}

func ParseConfig(filename string) (c Config, err error) {
	configFile, err := os.Open(filename)
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(configFile).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
	return
}

package utils

import (
	"encoding/json"
	"log"
	"os"
)

type Server struct {
	Network string `json:"network"`
	ApiKey  string `json:"api_key"`
	Manager string `json:"manager"'`
}
type Config struct {
	Server Server `json:"server"`
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

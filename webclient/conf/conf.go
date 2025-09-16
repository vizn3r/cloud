package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Conf struct {
	WebClient struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"webClient"`
	Port int `json:"port"`
}

var GlobalConf Conf

func LoadConfig(path string) error {
	raw, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open config file %s: %v", path, err)
		return fmt.Errorf("couldn't load '%s' config file", path)
	}

	parser := json.NewDecoder(raw)
	if err = parser.Decode(&GlobalConf); err != nil {
		log.Printf("Failed to decode config file %s: %v", path, err)
		return fmt.Errorf("couldn't decode '%s' config file", path)
	}

	return nil
}

func LoadFromBytes(data []byte) error {
	parser := json.NewDecoder(bytes.NewReader(data))
	if err := parser.Decode(&GlobalConf); err != nil {
		return fmt.Errorf("couldn't decode config file")
	}

	return nil
}

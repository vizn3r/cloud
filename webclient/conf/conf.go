package conf

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Conf struct {
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

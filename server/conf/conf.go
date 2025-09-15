package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type Conf struct {
	WebClientHost string `json:"WebClientHost"`
}

var GlobalConf Conf

func LoadConfig(path string) error {
	raw, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("couldn't load '%s' config file", path)
	}

	parser := json.NewDecoder(raw)
	if err = parser.Decode(&GlobalConf); err != nil {
		return fmt.Errorf("couldn't decode '%s' config file", path)
	}

	return nil
}

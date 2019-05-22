package hikaricommon

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func LoadConfig(filePath string, config interface{}) error {
	log.Printf("loading config file '%v'\n", filePath)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}

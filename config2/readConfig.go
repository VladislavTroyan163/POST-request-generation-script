package config2

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigApp struct {
	ServerURLFL            string `json:"serverURLFL"`
	ServerURLUL            string `json:"serverURLUL"`
	ServerURLAuthorization string `json:"serverURLAuthorization"`
	WorkersCount           int    `json:"workersCount"`
	Mode                   string `json:"mode"`
	WorkMonth              string `json:"workMonth"`
}

func ReadConfigApp() (ConfigApp, error) {
	var config ConfigApp
	file, err := os.OpenFile("configProgram.json", os.O_RDWR, 0644)
	if err != nil {
		return ConfigApp{}, fmt.Errorf("readConfig: ошибка чтения конфигурации: \n%w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return ConfigApp{}, fmt.Errorf("readConfig: ошибка чтения конфигурации: \n%w", err)
	}
	return config, nil
}

package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ServerURL               string `json:"serverURL"`
	ServerURLAuthorization  string `json:"serverURLAuthorization"`
	WorkersCount            int    `json:"workersCount"`
	Mode                    string `json:"mode"`
	ДатаНачала              string `json:"датаНачала"`
	ДатаКонца               string `json:"датаКонца"`
	ПоставщикиСписок        string `json:"поставщикиСписок"`
	УслугиСписокПоставщиков string `json:"услугиСписокПоставщиков"`
	WorkMonth               string `json:"workMonth"`
}

func ReadConfig(nameFileConfig string) (Config, error) {
	var config Config
	file, err := os.OpenFile(nameFileConfig, os.O_RDWR, 0644)
	if err != nil {
		return Config{}, fmt.Errorf("readConfig: ошибка чтения конфигурации: \n%w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("readConfig: ошибка чтения конфигурации: \n%w", err)
	}
	return config, nil
}

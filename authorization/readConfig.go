package authorization

import (
	"encoding/json"
	"fmt"
	"os"

	configAppPack "github.com/VladislavTroyan163/POST-request-generation-script/config"
)

type Config struct {
	ConfigApp     configAppPack.ConfigApp
	ConfigRequest ConfigRequest
}

type ConfigRequest struct {
}

func ReadConfig() (Config, error) {
	configApp, err := configAppPack.ReadConfigApp()
	if err != nil {
		return Config{}, fmt.Errorf("readConfig: ошибка чтения конфигурации: \n%w", err)
	}

	var configRequest ConfigRequest
	file, err := os.OpenFile("authorization/configRequest.json", os.O_RDWR, 0644)
	if err != nil {
		return Config{}, fmt.Errorf("readConfig: ошибка чтения конфигурации: \n%w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&configRequest)
	if err != nil {
		return Config{}, fmt.Errorf("readConfig: ошибка чтения конфигурации: \n%w", err)
	}

	Config := Config{
		ConfigApp:     configApp,
		ConfigRequest: configRequest,
	}
	return Config, nil
}

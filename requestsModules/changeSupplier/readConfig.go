package changeSupplier

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
	ДатаНачалаFL              string `json:"датаНачалаFL"`
	ДатаНачалаUL              string `json:"датаНачалаUL"`
	ДатаКонцаFL               string `json:"датаКонцаFL"`
	ДатаКонцаUL               string `json:"датаКонцаUL"`
	ПоставщикиСписокFL        string `json:"поставщикиСписокFL"`
	ПоставщикиСписокUL        string `json:"поставщикиСписокUL"`
	УслугиСписокПоставщиковFL string `json:"услугиСписокПоставщиковFL"`
	УслугиСписокПоставщиковUL string `json:"услугиСписокПоставщиковUL"`
}

func ReadConfig() (Config, error) {
	configApp, err := configAppPack.ReadConfigApp()
	if err != nil {
		return Config{}, fmt.Errorf("readConfig: ошибка чтения конфигурации: \n%w", err)
	}

	var configRequest ConfigRequest
	file, err := os.OpenFile("requestsModules/changeSupplier/configRequest.json", os.O_RDWR, 0644)
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

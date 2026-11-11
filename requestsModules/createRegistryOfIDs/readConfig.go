package createRegistryOfIDs

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
	ObjectName string `json:"objectName"`
	MethodName string `json:"methodName"`
	TextMeter  string `json:"textMeter"`
	TextLs     string `json:"textLs"`
	TextTu     string `json:"textTu"`
}

func ReadConfig() (Config, error) {
	configApp, err := configAppPack.ReadConfigApp()
	if err != nil {
		return Config{}, fmt.Errorf("readConfig: ошибка чтения конфигурации: \n%w", err)
	}

	var configRequest ConfigRequest
	file, err := os.OpenFile("requestsModules/createRegistryOfIDs/configRequest.json", os.O_RDWR, 0644)
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

func ReadFileId() ([]Row_idObject, error) {

	config, err := ReadConfig()
	if err != nil {
		return []Row_idObject{}, fmt.Errorf("ReadFileId: ошибка чтения конфигурации: \n%w", err)
	}
	fileName := ""
	switch config.ConfigApp.Mode {
	case "FL":
		fileName = "registries/registryLs/registryLs.json"
	case "UL":
		fileName = "registries/registryTu/registryTu.json"
	case "ASD":
		fileName = "registries/registryMeter/registryMeter.json"
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		return []Row_idObject{}, fmt.Errorf("ReadFileId: ошибка чтения файла: \n%w", err)
	}

	var row_idNum []Row_idObject
	err = json.Unmarshal(data, &row_idNum)
	if err != nil {
		return []Row_idObject{}, fmt.Errorf("ReadFileId: ошибка сериализации: \n%w", err)
	}
	return row_idNum, nil
}

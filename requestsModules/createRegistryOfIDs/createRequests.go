package createregistryofids

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	Tasks []Tasks `json:"tasks"`
	Info  Info    `json:"info"`
}

type Tasks struct {
	ObjectName string `json:"objectName"`
	MethodName string `json:"methodName"`
	Params     Params `json:"params"`
}

type Params struct {
	Фильтр Фильтр `json:"фильтр"`
}

type Фильтр struct {
	Номерлс         int  `json:"номерлс"`
	РазвернутьПапки int  `json:"развернутьПапки"`
	СКорня          int  `json:"сКорня"`
	Переход         bool `json:"переход"`
}

type Info struct {
	WorkMonth string `json:"workMonth"`
}

func createФильтр(номерЛС int) (Фильтр, error) {
	ФильтрData := Фильтр{
		Номерлс:         номерЛС,
		РазвернутьПапки: 1,
		СКорня:          1,
		Переход:         false,
	}
	return ФильтрData, nil
}

func createParams(лицевойСчет int) (Params, error) {
	ФильтрData, err := createФильтр(лицевойСчет)
	if err != nil {
		return Params{}, err
	}

	ParamsData := Params{
		Фильтр: ФильтрData,
	}
	return ParamsData, nil
}

func createTasks(params Params) Tasks {
	TasksData := Tasks{
		ObjectName: "ЛицевыеСчета",
		MethodName: "получить",
		Params:     params,
	}
	return TasksData
}

func createInfo(config Config) (Info, error) {
	InfoData := Info{
		WorkMonth: config.ConfigApp.WorkMonth,
	}
	return InfoData, nil
}

func CreateRequestСreateRegistryOfIDs(лицевойСчет int, config Config) (string, error) {
	ParamsData, err := createParams(лицевойСчет)
	if err != nil {
		return "", fmt.Errorf("CreateRequestСreateRegistryOfIDs: ошибка чтения конфигурации: \n%w", err)
	}
	TasksData := createTasks(ParamsData)
	InfoData, err := createInfo(config)
	if err != nil {
		return "", fmt.Errorf("CreateRequestСreateRegistryOfIDs: ошибка чтения конфигурации: \n%w", err)
	}
	Request := Request{
		Tasks: []Tasks{TasksData},
		Info:  InfoData,
	}

	jsonRequest, err := json.Marshal(Request)
	if err != nil {
		return "", fmt.Errorf("CreateRequestСreateRegistryOfIDs: ошибка сериализации: \n%w", err)
	}

	return string(jsonRequest), nil
}

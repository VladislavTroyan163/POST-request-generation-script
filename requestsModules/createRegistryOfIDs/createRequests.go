package createRegistryOfIDs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Request struct {
	Tasks []Task `json:"tasks"` // 1. Сделали срезом []Task и добавили тег "tasks"
	Info  Info   `json:"info"`
}

type Task struct {
	ObjectName string `json:"objectName"`
	MethodName string `json:"methodName"`
	Params     Params `json:"params"` // 2. Исправили тег с "tasks" на "params"
}

type Params struct {
	Text string `json:"text"`
}

type Info struct {
	WorkMonth string `json:"workMonth"`
}

func CreateRequestNumRegistryOfIDs(NumberLsTu int, номерПУ string, config Config) (string, error) {
	text := ""
	switch config.ConfigApp.Mode {
	case "FL":
		text = config.ConfigRequest.TextLs
		text = strings.ReplaceAll(text, "!!!!!", strconv.Itoa(NumberLsTu))
	case "UL":
		text = config.ConfigRequest.TextTu
		text = strings.ReplaceAll(text, "!!!!!", strconv.Itoa(NumberLsTu))
	case "ASD":
		text = config.ConfigRequest.TextMeter
		text = strings.ReplaceAll(text, "!!!!!", номерПУ)
	}
	request := Request{
		Tasks: []Task{
			{
				ObjectName: config.ConfigRequest.ObjectName,
				MethodName: config.ConfigRequest.MethodName,
				Params: Params{
					Text: text,
				},
			},
		},
		Info: Info{
			WorkMonth: config.ConfigApp.WorkMonth,
		},
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("CreateRequestСreateRegistryOfIDs: ошибка сериализации: \n%w", err)
	}

	return string(jsonRequest), nil

}

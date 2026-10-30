package changesupplier

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
	Запись   Запись `json:"запись"`
	Владелец int    `json:"владелец"`
}

type Запись struct {
	НомерЗаписи                int     `json:"$номерЗаписи"`
	ЭтоПапка                   bool    `json:"$этоПапка"`
	Папка                      *string `json:"$папка"`
	ЕстьИзменения              bool    `json:"$естьИзменения"`
	НомерУслуги                int     `json:"номер услуги"`
	СчетСписокПоставщиков      int     `json:"счет-список поставщиков"`
	Наименование               string  `json:"наименование"`
	ПоставщикиСписокНазвание   string  `json:"поставщики-список>название"`
	ДатНач                     string  `json:"датнач"`
	ДатКнц                     string  `json:"даткнц"`
	ПоставщикиДоговорыУкНомер  int     `json:"поставщики-договоры ук>номер"`
	ПоставщикиДоговорыУкДатнач *string `json:"поставщики-договоры ук>датнач"`
	ПоставщикиСписок           string  `json:"поставщики-список"`
	УслугиСписокПоставщиков    string  `json:"услуги-список поставщиков"`
	Принимаетвседолги          int     `json:"принимаетвседолги"`
	ДолгДатНач                 *string `json:"@долгдатнач"`
	ДолгДатКнц                 *string `json:"@долгдаткнц"`
	Примечание                 string  `json:"примечание"`
	ВторойНомер                string  `json:"@2номер"`
	Название                   string  `json:"$название"`
	ПолеПапка                  string  `json:"$полеПапка"`
	ТекущаяСтраница            int     `json:"$текущаяСтраница"`
	ТолькоЧтение               bool    `json:"$толькоЧтение"`
	Блокировка                 int     `json:"$блокировка"`
	БлокирующийПользователь    string  `json:"$блокирующийПользователь"`
}

type Info struct {
	WorkMonth string `json:"workMonth"`
}

func createЗапись(лицевойСчет int, config Config) (Запись, error) {
	датаНачала := ""
	датаКонца := ""
	поставщикиСписок := ""
	услугиСписокПоставщиков := ""

	switch config.ConfigApp.Mode {
	case "FL":
		датаНачала = config.ConfigRequest.датаНачалаFL
		датаКонца = config.ConfigRequest.датаКонцаFL
		поставщикиСписок = config.ConfigRequest.поставщикиСписокFL
		услугиСписокПоставщиков = config.ConfigRequest.услугиСписокПоставщиковFL
	case "UL":
		датаНачала = config.ConfigRequest.датаНачалаUL
		датаКонца = config.ConfigRequest.датаКонцаUL
		поставщикиСписок = config.ConfigRequest.поставщикиСписокUL
		услугиСписокПоставщиков = config.ConfigRequest.услугиСписокПоставщиковUL
	}

	ЗаписьData := Запись{
		НомерЗаписи:             -1,
		ЕстьИзменения:           true,
		СчетСписокПоставщиков:   лицевойСчет,
		ДатНач:                  датаНачала,
		ДатКнц:                  датаКонца,
		ПоставщикиСписок:        поставщикиСписок,
		УслугиСписокПоставщиков: услугиСписокПоставщиков,
	}
	return ЗаписьData, nil
}

func createParams(лицевойСчет int, config Config) (Params, error) {
	ЗаписьData, err := createЗапись(лицевойСчет, config)
	if err != nil {
		return Params{}, err
	}

	ParamsData := Params{
		Запись:   ЗаписьData,
		Владелец: лицевойСчет,
	}
	return ParamsData, nil
}

func createTasks(params Params) Tasks {
	TasksData := Tasks{
		ObjectName: "ЛицевыеСчета.Поставщики",
		MethodName: "вставить",
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

func CreateRequestChangeSupplier(лицевойСчет int, config Config) (string, error) {
	ParamsData, err := createParams(лицевойСчет, config)
	if err != nil {
		return "", fmt.Errorf("CreateRequestChangeSupplier: ошибка чтения конфигурации: \n%w", err)
	}
	TasksData := createTasks(ParamsData)
	InfoData, err := createInfo(config)
	if err != nil {
		return "", fmt.Errorf("CreateRequestChangeSupplier: ошибка чтения конфигурации: \n%w", err)
	}
	Request := Request{
		Tasks: []Tasks{TasksData},
		Info:  InfoData,
	}

	jsonRequest, err := json.Marshal(Request)
	if err != nil {
		return "", fmt.Errorf("CreateRequestChangeSupplier: ошибка сериализации: \n%w", err)
	}

	return string(jsonRequest), nil
}

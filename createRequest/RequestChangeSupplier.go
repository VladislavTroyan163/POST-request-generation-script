package createrequest

import (
	"encoding/json"
	"fmt"
)

const (
	датаНачала              = "2026-10-01T00:00:00.000"
	датаКонца               = "2045-05-09T00:00:00.000"
	поставщикиСписок        = "25429"
	услугиСписокПоставщиков = "139067"
	workMonth               = "2026-09-01T00:00:00.000"
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

func createЗапись(лицевойСчет int) Запись {
	ЗаписьData := Запись{
		НомерЗаписи:             -1,
		ЕстьИзменения:           true,
		СчетСписокПоставщиков:   лицевойСчет,
		ДатНач:                  датаНачала,
		ДатКнц:                  датаКонца,
		ПоставщикиСписок:        поставщикиСписок,
		УслугиСписокПоставщиков: услугиСписокПоставщиков,
	}
	return ЗаписьData
}

func createParams(лицевойСчет int) Params {
	ЗаписьData := createЗапись(лицевойСчет)

	ParamsData := Params{
		Запись:   ЗаписьData,
		Владелец: лицевойСчет,
	}
	return ParamsData
}

func createTasks(params Params) Tasks {
	TasksData := Tasks{
		ObjectName: "ЛицевыеСчета.Поставщики",
		MethodName: "вставить",
		Params:     params,
	}
	return TasksData
}

func createInfo() Info {
	InfoData := Info{
		WorkMonth: workMonth,
	}
	return InfoData
}

func CreateRequest(лицевойСчет int) (string, error) {
	ParamsData := createParams(лицевойСчет)
	TasksData := createTasks(ParamsData)
	InfoData := createInfo()

	Request := Request{
		Tasks: []Tasks{TasksData},
		Info:  InfoData,
	}

	jsonRequest, err := json.Marshal(Request)
	if err != nil {
		return "", fmt.Errorf("createRequest: ошибка сериализации: %w", err)
	}

	return string(jsonRequest), nil
}

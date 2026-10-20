package registry

import (
	"encoding/json"
	"fmt"
	"os"
)

type Rows_ids struct {
	Row_id int `json:"row_id"`
}

func ReadRegistry(mode string) ([]Rows_ids, error) {
	var ObjectsRow_id []Rows_ids
	var filePath string

	switch mode {
	case "UL":
		filePath = "registry\\registryFL.json"
	case "FL":
		filePath = "registry\\registryUL.json"
	default:
		return []Rows_ids{}, fmt.Errorf("readConfig: ошибка чтения конфигурации: \nНе верно указан mode.")
	}

	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return []Rows_ids{}, fmt.Errorf("readConfig: ошибка чтения конфигурации: \n%w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&ObjectsRow_id)
	if err != nil {
		return []Rows_ids{}, fmt.Errorf("readConfig: ошибка чтения конфигурации: \n%w", err)
	}
	return ObjectsRow_id, nil
}

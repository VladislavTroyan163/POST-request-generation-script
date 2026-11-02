package registry

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Rows_ids struct {
	Row_id int `json:"row_id"`
}

func ReadRegistry(mode string) ([]Rows_ids, error) {
	var ObjectsRow_id []Rows_ids
	filePath := "registry/registry.xlsx"

	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return []Rows_ids{}, fmt.Errorf("ReadRegistry: ошибка чтения реестра: \n%w", err)
	}
	defer file.Close()

	return ObjectsRow_id, nil
}

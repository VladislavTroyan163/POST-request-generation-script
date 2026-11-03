package registry

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Rows_ids struct {
	Row_id int `json:"row_id"`
}

func ReadRegistry(mode string, fileName string, sheetNum int, colName string) ([]Rows_ids, error) {
	var resultRow_id []Rows_ids
	filePath := ""

	switch mode {
	case "FL":
		filePath = fmt.Sprint("registry/registryLs/", fileName)
	case "UL":
		filePath = fmt.Sprint("registry/registryTu/", fileName)
	}

	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return []Rows_ids{}, fmt.Errorf("ReadRegistry: ошибка чтения реестра: \n%w", err)
	}
	defer file.Close()

	sheet := file.GetSheetName(sheetNum - 1)
	if sheet == "" {
		return []Rows_ids{}, fmt.Errorf("ReadRegistry: лист с индексом %d не найден", sheetNum)
	}

	rows, err := file.Rows(sheet)
	if err != nil {
		return []Rows_ids{}, fmt.Errorf("ReadRegistry: ошибка чтения реестра: \n%w", err)
	}
	defer rows.Close()

	colomnId := -1
	if rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return []Rows_ids{}, fmt.Errorf("ReadRegistry: ошибка чтения строки: \n%w", err)
		}
		for i, colomnName := range columns {
			if strings.TrimSpace(colomnName) == strings.TrimSpace(colName) {
				colomnId = i
				break
			}
		}
	}

	if colomnId == -1 {
		return []Rows_ids{}, fmt.Errorf("ReadRegistry: колонка %q не найдена", colName)
	}

	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return []Rows_ids{}, fmt.Errorf("ReadRegistry: ошибка чтения строки: \n%w", err)
		}

		if len(columns) > colomnId {
			row_id_str := strings.TrimSpace(columns[colomnId])
			if row_id_str == "" {
				continue
			}
			row_id_int, err := strconv.Atoi(row_id_str)
			if err != nil {
				// Заменить на логирование
				return []Rows_ids{}, fmt.Errorf("ReadRegistry: ошибка преобразования типов string -> int: \n%w", err)
			}
			resultRow_id = append(resultRow_id, Rows_ids{Row_id: row_id_int})
		}
	}

	return resultRow_id, nil
}

package registries

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Numbers struct {
	Numbers_ls    int
	Numbers_tu    int
	Numbers_meter string
}

func ReadRegistry(mode string, fileName string, sheetNum int, colName string) ([]Numbers, error) {
	var resultNumbers []Numbers
	filePath := ""

	switch mode {
	case "FL":
		filePath = fmt.Sprint("registries/registryLs/", fileName)
	case "UL":
		filePath = fmt.Sprint("registries/registryTu/", fileName)
	case "ASD":
		filePath = fmt.Sprint("registries/registryMeter/", fileName)
	}

	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return []Numbers{}, fmt.Errorf("ReadRegistryNumber: ошибка чтения реестра: \n%w", err)
	}
	defer file.Close()

	sheet := file.GetSheetName(sheetNum - 1)
	if sheet == "" {
		return []Numbers{}, fmt.Errorf("ReadRegistryNumber: лист с индексом %d не найден", sheetNum)
	}

	rows, err := file.Rows(sheet)
	if err != nil {
		return []Numbers{}, fmt.Errorf("ReadRegistryNumber: ошибка чтения реестра: \n%w", err)
	}
	defer rows.Close()

	colomnId := -1
	if rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return []Numbers{}, fmt.Errorf("ReadRegistryNumber: ошибка чтения строки: \n%w", err)
		}
		for i, colomnName := range columns {
			if strings.TrimSpace(colomnName) == strings.TrimSpace(colName) {
				colomnId = i
				break
			}
		}
	}

	if colomnId == -1 {
		return []Numbers{}, fmt.Errorf("ReadRegistryNumber: колонка %q не найдена", colName)
	}

	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return []Numbers{}, fmt.Errorf("ReadRegistryNumber: ошибка чтения строки: \n%w", err)
		}

		if len(columns) > colomnId {
			numbersStr := strings.TrimSpace(columns[colomnId])
			if numbersStr == "" {
				continue
			}
			if mode == "ASD" {
				resultNumbers = append(resultNumbers, Numbers{Numbers_meter: numbersStr})
				continue
			}

			numbersInt, err := strconv.Atoi(numbersStr)
			if err != nil {
				// Заменить на логирование
				return []Numbers{}, fmt.Errorf("ReadRegistryNumber: ошибка преобразования типов string -> int: \n%w", err)
			}
			switch mode {
			case "FL":
				resultNumbers = append(resultNumbers, Numbers{Numbers_ls: numbersInt})
			case "UL":
				resultNumbers = append(resultNumbers, Numbers{Numbers_tu: numbersInt})
			}

		}
	}
	if len(resultNumbers) == 0 {
		return []Numbers{}, fmt.Errorf("Реестр пуст")
	}

	return resultNumbers, nil
}

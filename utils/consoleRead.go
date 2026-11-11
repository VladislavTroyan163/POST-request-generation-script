package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadConsoleSettingsFile() (fileName string, sheetNumInt int, columnName string, err error) {
	count := 0
	for i := 0; i < 1; {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Введите название файла: ")
		scanner.Scan()
		fileName = strings.TrimSpace(scanner.Text())
		if fileName == "" {
			fmt.Println("Ошибка: Название файла не может быть пустым!")
			count++
			continue
		}

		fmt.Print("Введите порядковый номер листа: ")
		scanner.Scan()
		sheetNumStr := strings.TrimSpace(scanner.Text())
		sheetNumInt, err = strconv.Atoi(sheetNumStr)
		if err != nil {
			fmt.Println("Ошибка: порядковый номер листа - число!")
			count++
			continue
		}

		fmt.Print("Введите название столбца: ")
		scanner.Scan()
		columnName = strings.TrimSpace(scanner.Text())

		if columnName == "" {
			fmt.Println("Ошибка: Название столбца не может быть пустым!")
			count++
			continue
		}
		if count == 3 {
			err = fmt.Errorf("Превышено число попыток: %d", count)
			return
		}
		i++
	}
	return

}

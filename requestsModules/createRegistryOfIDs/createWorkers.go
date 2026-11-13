package createRegistryOfIDs

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	loggerpack "github.com/VladislavTroyan163/POST-request-generation-script/logger"
	registriespack "github.com/VladislavTroyan163/POST-request-generation-script/registries"
)

type Row_idObject struct {
	Row_id  int    `json:"row_id"`
	Numbers string `json:"numbers"`
}

func CreateWorkersNumRegistryOfIDs(
	registryNumber []registriespack.Numbers,
	workersCount int,
	tokenAuthorization string,
) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var fileName string
	row_idObject := []Row_idObject{}
	jobs := make(chan int)
	rowIDs := make([]int, 0, len(registryNumber))

	//Создание логирования
	logger, logFile, err := loggerpack.CreateLogger("преобразование к row_id")
	if err != nil {
		return fmt.Errorf("CreateWorkersNumRegistryOfIDs: ошибка создания логера: \n%w", err)
	}
	defer logFile.Close()
	//--------------------

	configCreateRegistryOfIDs, err := ReadConfig()
	if err != nil {
		return fmt.Errorf("CreateWorkersNumRegistryOfIDs: ошибка чтения конфигурации: \n%w", err)
	}
	startAll := time.Now()
	for w := 0; w < workersCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				лицевойСчет := 0
				номерПУ := ""
				logNumber := ""
				logObjectName := ""

				switch configCreateRegistryOfIDs.ConfigApp.Mode {
				case "FL":
					лицевойСчет = registryNumber[i].Numbers_ls
					fileName = "registries/registryLs/registryLs.json"
					logNumber = strconv.Itoa(лицевойСчет)
					logObjectName = "Лицевой счет"
				case "UL":
					лицевойСчет = registryNumber[i].Numbers_tu
					fileName = "registries/registryTu/registryTu.json"
					logNumber = strconv.Itoa(лицевойСчет)
					logObjectName = "Точка учета"
				case "ASD":
					номерПУ = registryNumber[i].Numbers_meter
					fileName = "registries/registryMeter/registryMeter.json"
					logNumber = номерПУ
					logObjectName = "Номер счетчика"
				}
				start := time.Now()
				requestNumbers, err := CreateRequestNumRegistryOfIDs(лицевойСчет, номерПУ, configCreateRegistryOfIDs)
				if err != nil {
					logger.Printf("%s %s: ошибка создания запроса: %v\n\n", logObjectName, logNumber, err)
					continue
				}

				resultRequest, row_id, err := ExecuteRequestNumRegistryOfIDs(requestNumbers, configCreateRegistryOfIDs, tokenAuthorization)
				elapsed := time.Since(start)
				if err != nil {
					logger.Printf("%s %s: ошибка HTTP-запроса: %v\n\n", logObjectName, logNumber, err)
					continue
				}

				logger.Printf("%s %s\nРезультат: %vВремя выполнения запроса: %v\n\n", logObjectName, logNumber, resultRequest, elapsed)

				mu.Lock()
				row_idObject = append(row_idObject, Row_idObject{Row_id: row_id, Numbers: logNumber})
				rowIDs = append(rowIDs, row_id)
				mu.Unlock()
			}
		}()
	}

	for i := range registryNumber {
		jobs <- i
	}
	close(jobs)

	wg.Wait()

	logger.Printf("Всего обработано заданий: %d\nОбщее время: %v\n", len(rowIDs), time.Since(startAll))

	jsonData, err := json.MarshalIndent(row_idObject, "", "    ")
	if err != nil {
		return fmt.Errorf("CreateWorkersNumRegistryOfIDs: ошибка сериализации: \n%w", err)
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("CreateWorkersNumRegistryOfIDs: ошибка чтения файла: \n%w", err)
	}
	defer file.Close()

	file.Write(jsonData)

	return nil
}

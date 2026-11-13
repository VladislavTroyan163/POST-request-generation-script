package changeSupplier

import (
	"fmt"
	"sync"
	"time"

	loggerpack "github.com/VladislavTroyan163/POST-request-generation-script/logger"
	createregistryofidspack "github.com/VladislavTroyan163/POST-request-generation-script/requestsModules/createRegistryOfIDs"
)

type Row_idObject struct {
	Row_id  int    `json:"row_id"`
	Numbers string `json:"numbers"`
}

func CreateWorkersChangeSupplier(workersCount int, tokenAuthorization string, rowIDs []createregistryofidspack.Row_idObject) error {
	//Создание логирования
	logger, logFile, err := loggerpack.CreateLogger("внесение поставщиков")
	if err != nil {
		return fmt.Errorf("CreateWorkersChangeSupplier: ошибка создания логера: \n%w", err)
	}
	defer logFile.Close()
	//--------------------

	configChangeSuppliers, err := ReadConfig()
	if err != nil {
		return fmt.Errorf("CreateWorkersChangeSupplier: ошибка чтения конфигурации: \n%w", err)
	}

	mode := configChangeSuppliers.ConfigApp.Mode
	if mode != "FL" && mode != "UL" {
		return fmt.Errorf("CreateWorkersChangeSupplier: ошибка: \nневозможно выполнить процедуру в модуле %s", mode)
	}

	var wg sync.WaitGroup
	jobs := make(chan int)

	startAll := time.Now()

	for w := 0; w < workersCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				start := time.Now()
				row_idLsTu := rowIDs[i].Row_id

				request, err := CreateRequestChangeSupplier(row_idLsTu, configChangeSuppliers)
				if err != nil {
					logger.Printf("row_id %d: ошибка создания запроса: %v\n\n", row_idLsTu, err)
					continue
				}

				resultRequest, err := ExecuteRequestChangeSupplier(request, row_idLsTu, configChangeSuppliers, tokenAuthorization)

				elapsed := time.Since(start)
				if err != nil {
					logger.Printf("row_id %d: ошибка HTTP-запроса: %v\n\n", row_idLsTu, err)
					continue
				}
				logger.Printf("Результат: %vВремя выполнения запроса: %v\n\n", resultRequest, elapsed)
			}
		}()
	}

	for i := range rowIDs {
		jobs <- i
	}
	close(jobs)

	wg.Wait()

	logger.Printf("Всего обработано заданий: %d\nОбщее время: %v\n", len(rowIDs), time.Since(startAll))
	return nil
}

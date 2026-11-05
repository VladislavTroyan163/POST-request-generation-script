package main

import (
	"fmt"
	"sync"
	"time"

	configpack "github.com/VladislavTroyan163/POST-request-generation-script/config"
	. "github.com/VladislavTroyan163/POST-request-generation-script/logger"
	. "github.com/VladislavTroyan163/POST-request-generation-script/registries"
	authorizationpack "github.com/VladislavTroyan163/POST-request-generation-script/requestsModules/authorization"
	changesupplierspack "github.com/VladislavTroyan163/POST-request-generation-script/requestsModules/changeSupplier"
	createregistryofidspack "github.com/VladislavTroyan163/POST-request-generation-script/requestsModules/createRegistryOfIDs"
)

func main() {
	//Создание логирования
	logger, logFile := CreateLogger()
	defer logFile.Close()
	//----------------------------------------------------------------------------------

	//Чтение конфига приложения
	configApp, err := configpack.ReadConfigApp()
	if err != nil {
		fmt.Println("main: ошибка чтения конфигурации приложения: \n", err)
		return
	}
	//----------------------------------------------------------------------------------

	//Авторизация: конфигурация, формирование запроса, выполнение запроса.
	configAuthorization, err := authorizationpack.ReadConfig()
	if err != nil {
		fmt.Println("main: ошибка чтения конфигурации: \n", err)
		return
	}
	requestAuthorization, err := authorizationpack.CreateRequestAuthorization()
	if err != nil {
		fmt.Println("main: ошибка при создании запроса авторизации: \n", err)
		return
	}
	tokenAuthorization, err := authorizationpack.ExecuteRequestAuthorization(requestAuthorization, configAuthorization)
	if err != nil {
		fmt.Println("main: ошибка при выполнении запроса авторизации: \n", err)
		return
	}
	//----------------------------------------------------------------------------------

	//Чтение реестра лицевых счетов
	registryNumber, err := ReadRegistry(configApp.Mode, "registryLs.xlsx", 1, "Номер ЛС")
	if err != nil {
		logger.Println(err)
		return
	}
	//----------------------------------------------------------------------------------

	//Установка количества потоков
	workersCount := configApp.WorkersCount

	if len(registryNumber) < workersCount {
		workersCount = len(registryNumber)
	}
	if workersCount == 0 {
		logger.Println("Реестр пуст")
		return
	}
	//----------------------------------------------------------------------------------

	//Перевод объектов реестра в формат списка row_id
	configCreateRegistryOfIDs, err := createregistryofidspack.ReadConfig()
	if err != nil {
		fmt.Println("main: ошибка чтения конфигурации: \n", err)
		return
	}
	rowIDs := createregistryofidspack.CreateWorkersCreateRegistryOfIDs(
		registryNumber,
		workersCount,
		configCreateRegistryOfIDs,
		tokenAuthorization,
	)
	fmt.Print(rowIDs)
	//----------------------------------------------------------------------------------

	var wg sync.WaitGroup
	jobs := make(chan int)

	configChangeSuppliers, err := changesupplierspack.ReadConfig()
	if err != nil {
		fmt.Println("main: ошибка чтения конфигурации: \n", err)
		return
	}

	startAll := time.Now()

	for w := 0; w < workersCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				start := time.Now()
				лицевойСчет := rowIDs[i]

				request, err := changesupplierspack.CreateRequestChangeSupplier(лицевойСчет, configChangeSuppliers)
				if err != nil {
					logger.Printf(
						"ЛС %d: ошибка создания запроса: %v",
						лицевойСчет, err)
					break
				}

				resultRequest, err := changesupplierspack.ExecuteRequestChangeSupplier(request, лицевойСчет, configChangeSuppliers, tokenAuthorization)

				elapsed := time.Since(start)
				if err != nil {
					logger.Printf(
						"ЛС %d: ошибка HTTP-запроса: %v \n\n",
						лицевойСчет, err)
					continue
				}
				fmt.Println("Результат запроса:", resultRequest)
				logger.Printf(
					"Результат: %vВремя выполнения запроса: %v\n\n",
					resultRequest,
					elapsed,
				)
			}
		}()
	}

	for i := range registryNumber {
		jobs <- i
	}
	close(jobs)

	wg.Wait()

	logger.Printf("Всего обработано заданий: %d\n",
		len(registryNumber))
	logger.Printf("Общее время: %v\n", time.Since(startAll))
}

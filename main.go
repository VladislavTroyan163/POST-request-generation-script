package main

import (
	"fmt"
	"sync"
	"time"

	configpack "github.com/VladislavTroyan163/POST-request-generation-script/config"
	. "github.com/VladislavTroyan163/POST-request-generation-script/logger"
	. "github.com/VladislavTroyan163/POST-request-generation-script/registry"
	authorizationpack "github.com/VladislavTroyan163/POST-request-generation-script/requestsModules/authorization"
	changesupplierspack "github.com/VladislavTroyan163/POST-request-generation-script/requestsModules/changesupplier"
	createregistryofidspack "github.com/VladislavTroyan163/POST-request-generation-script/requestsModules/createregistryofids"
)

func main() {
	//Создание логирования
	logger, logFile := CreateLogger()
	defer logFile.Close()

	//Чтение конфига приложения
	configApp, err := configpack.ReadConfigApp()
	if err != nil {
		fmt.Println("main: ошибка чтения конфигурации приложения: \n", err)
		return
	}

	//Чтение реестра лицевых счетов
	sliceЛицевыхСчетов, err := ReadRegistry(configApp.Mode)
	if err != nil {
		logger.Println(err)
		return
	}

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

	//Установка количества потоков
	workersCount := configApp.WorkersCount

	if len(sliceЛицевыхСчетов) < workersCount {
		workersCount = len(sliceЛицевыхСчетов)
	}
	if workersCount == 0 {
		logger.Println("Реестр пуст")
		return
	}

	//----------------------------

	configChangeSuppliers, err := changesupplierspack.ReadConfig()
	if err != nil {
		fmt.Println("main: ошибка чтения конфигурации: \n", err)
		return
	}
	//-------------

	var wg sync.WaitGroup
	jobs := make(chan int)

	startAll := time.Now()

	for w := 0; w < workersCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				start := time.Now()
				лицевойСчет := sliceЛицевыхСчетов[i].Row_id

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

	for i := range sliceЛицевыхСчетов {
		jobs <- i
	}
	close(jobs)

	wg.Wait()

	logger.Printf("Всего обработано заданий: %d\n",
		len(sliceЛицевыхСчетов))
	logger.Printf("Общее время: %v\n", time.Since(startAll))

	configCreateRegistryOfIDs, err := createregistryofidspack.ReadConfig()
	if err != nil {
		fmt.Println("main: ошибка чтения конфигурации: \n", err)
		return
	}
	Request, err := createregistryofidspack.CreateRequestСreateRegistryOfIDs(3128055, configCreateRegistryOfIDs)
	fmt.Println(Request)
}

package main

import (
	"fmt"
	"sync"
	"time"

	. "github.com/VladislavTroyan163/POST-request-generation-script/config"
	. "github.com/VladislavTroyan163/POST-request-generation-script/createRequests"
	createregistryofidspack "github.com/VladislavTroyan163/POST-request-generation-script/createregistryofids"
	. "github.com/VladislavTroyan163/POST-request-generation-script/executeRequests"
	. "github.com/VladislavTroyan163/POST-request-generation-script/logger"
	. "github.com/VladislavTroyan163/POST-request-generation-script/registry"
)

func main() {
	logger, logFile := CreateLogger()
	defer logFile.Close()

	config, err := ReadConfig("config\\configUL.json")
	if err != nil {
		fmt.Println("main: ошибка чтения конфигурации: \n", err)
		return
	}

	sliceЛицевыхСчетов, err := ReadRegistry(config.Mode)
	if err != nil {
		logger.Println(err)
		return
	}

	requestAuthorization, err := CreateRequestAuthorization()
	if err != nil {
		fmt.Println("main: ошибка при создании запроса авторизации: \n", err)
		return
	}
	tokenAuthorization, err := ExecuteRequestAuthorization(requestAuthorization, config)
	if err != nil {
		fmt.Println("main: ошибка при выполнении запроса авторизации: \n", err)
		return
	}

	workersCount := config.WorkersCount

	if len(sliceЛицевыхСчетов) < workersCount {
		workersCount = len(sliceЛицевыхСчетов)
	}
	if workersCount == 0 {
		logger.Println("Реестр пуст")
		return
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
				лицевойСчет := sliceЛицевыхСчетов[i].Row_id

				request, err := CreateRequestChangeSupplier(лицевойСчет, config)
				if err != nil {
					logger.Printf(
						"ЛС %d: ошибка создания запроса: %v",
						лицевойСчет, err)
					break
				}

				resultRequest, err := ExecuteRequestChangeSupplier(request, лицевойСчет, config, tokenAuthorization)

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

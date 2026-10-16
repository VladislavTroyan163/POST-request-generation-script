package main

import (
	"fmt"
	"sync"
	"time"

	. "github.com/VladislavTroyan163/POST-request-generation-script/createRequests"
	. "github.com/VladislavTroyan163/POST-request-generation-script/executeRequests"
	. "github.com/VladislavTroyan163/POST-request-generation-script/logger"
	. "github.com/VladislavTroyan163/POST-request-generation-script/registry"
)

func main() {
	logger, logFile := CreateLogger()
	defer logFile.Close()

	sliceЛицевыхСчетов, err := ReadRegistry()
	if err != nil {
		logger.Println(err)
		return
	}

	workersCount := 30

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

				request, err := CreateRequestChangeSupplier(лицевойСчет)
				if err != nil {
					logger.Printf(
						"ЛС %d: ошибка создания запроса: %v",
						лицевойСчет, err)
					break
				}

				resultRequest, err := ExecuteRequestChangeSupplier(request, лицевойСчет)

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
}

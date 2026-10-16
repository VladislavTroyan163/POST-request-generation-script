package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	. "github.com/VladislavTroyan163/POST-request-generation-script/createRequests"
	. "github.com/VladislavTroyan163/POST-request-generation-script/executeRequests"
	. "github.com/VladislavTroyan163/POST-request-generation-script/registry"
)

func main() {
	sliceЛицевыхСчетов, err := ReadRegistry()
	if err != nil {
		log.Fatal(err)
	}

	workersCount := 30

	if len(sliceЛицевыхСчетов) < workersCount {
		workersCount = len(sliceЛицевыхСчетов)
	}
	if workersCount == 0 {
		fmt.Println("Реестр пуст")
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
					log.Printf(
						"ЛС %d: ошибка создания запроса: %v",
						лицевойСчет, err)
					continue
				}

				resultRequest, err := ExecuteRequestChangeSupplier(request, лицевойСчет)

				elapsed := time.Since(start)
				if err != nil {
					log.Printf(
						"ЛС %d: ошибка HTTP-запроса: %v",
						лицевойСчет, err)
					continue
				}
				fmt.Printf(
					"ЛС: %d | Результат: %v | Время: %v\n",
					лицевойСчет,
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

	fmt.Printf("Всего обработано заданий: %d\n",
		len(sliceЛицевыхСчетов))
	fmt.Printf("Общее время: %v\n", time.Since(startAll))
}

package createRegistryOfIDs

import (
	"sync"

	registriespack "github.com/VladislavTroyan163/POST-request-generation-script/registries"
)

func CreateWorkersCreateRegistryOfIDs(
	registryNumber []registriespack.Numbers,
	workersCount int,
	configCreateRegistryOfIDs Config,
	tokenAuthorization string,
) []int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	jobs := make(chan int)
	rowIDs := make([]int, 0, len(registryNumber))

	for w := 0; w < workersCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				лицевойСчет := 0
				// start := time.Now()
				switch configCreateRegistryOfIDs.ConfigApp.Mode {
				case "FL":
					лицевойСчет = registryNumber[i].Numbers_ls
				case "UL":
					лицевойСчет = registryNumber[i].Numbers_tu
				}

				requestNumbers, err := CreateRequestСreateRegistryOfIDs(лицевойСчет, configCreateRegistryOfIDs)
				if err != nil {
					break
				}

				_, row_id, err := ExecuteRequestcreateRegistryOfIDs(requestNumbers, лицевойСчет, configCreateRegistryOfIDs, tokenAuthorization)

				// elapsed := time.Since(start)
				// if err != nil {
				// 	continue
				// }

				mu.Lock()
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
	return rowIDs
}

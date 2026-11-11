package createRegistryOfIDs

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

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

	configCreateRegistryOfIDs, err := ReadConfig()
	if err != nil {
		return fmt.Errorf("CreateRequestСreateRegistryOfIDs: ошибка чтения конфигурации: \n%w", err)
	}

	for w := 0; w < workersCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				лицевойСчет := 0
				номерПУ := ""
				logNumber := ""

				switch configCreateRegistryOfIDs.ConfigApp.Mode {
				case "FL":
					лицевойСчет = registryNumber[i].Numbers_ls
					fileName = "registries/registryLs/registryLs.json"
					logNumber = strconv.Itoa(лицевойСчет)
				case "UL":
					лицевойСчет = registryNumber[i].Numbers_tu
					fileName = "registries/registryTu/registryTu.json"
					logNumber = strconv.Itoa(лицевойСчет)
				case "ASD":
					номерПУ = registryNumber[i].Numbers_meter
					fileName = "registries/registryMeter/registryMeter.json"
					logNumber = номерПУ
				}
				start := time.Now()
				requestNumbers, err := CreateRequestNumRegistryOfIDs(лицевойСчет, номерПУ, configCreateRegistryOfIDs)
				if err != nil {
					break
				}

				_, row_id, err := ExecuteRequestNumRegistryOfIDs(requestNumbers, configCreateRegistryOfIDs, tokenAuthorization)

				elapsed := time.Since(start)
				fmt.Printf("Лицевому счету %d соответствует row id: %d\nВремя выполнения запроса: %v\n\n", лицевойСчет, row_id, elapsed)

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

	jsonData, err := json.MarshalIndent(row_idObject, "", "    ")
	if err != nil {
		return fmt.Errorf("CreateRequestСreateRegistryOfIDs: ошибка сериализации: \n%w", err)
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return fmt.Errorf("CreateWorkersNumRegistryOfIDs: ошибка чтения файла: \n%w", err)
	}
	defer file.Close()

	file.Write(jsonData)

	return nil
}

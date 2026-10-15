package main

import (
	"fmt"
	"log"
	"time"

	. "github.com/VladislavTroyan163/POST-request-generation-script/createRequests"
	. "github.com/VladislavTroyan163/POST-request-generation-script/executeRequests"
	. "github.com/VladislavTroyan163/POST-request-generation-script/registry"
)

func main() {
	лицевойСчет := 0

	sliceЛицевыхСчетов, err := ReadRegistry()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(sliceЛицевыхСчетов); i++ {
		start := time.Now() // Начало итерации
		лицевойСчет = sliceЛицевыхСчетов[i].Row_id

		request, err := CreateRequestChangeSupplier(лицевойСчет)
		if err != nil {
			log.Fatal(err)
		}

		resultRequest, err := ExecuteRequestChangeSupplier(request, лицевойСчет)
		if err != nil {
			log.Fatal(err)
		}
		elapsed := time.Since(start) // Время выполнения
		fmt.Println(request, resultRequest, elapsed)
	}
}

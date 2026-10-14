package main

import (
	"fmt"
	"log"

	. "github.com/VladislavTroyan163/POST-request-generation-script/createRequests"
	. "github.com/VladislavTroyan163/POST-request-generation-script/executeRequests"
)

func main() {
	лицевойСчет := 924126

	request, err := CreateRequestChangeSupplier(лицевойСчет)
	if err != nil {
		log.Fatal(err)
	}
	resultRequest, err := ExecuteRequestChangeSupplier(request, лицевойСчет)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(request, resultRequest)
}

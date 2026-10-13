package main

import (
	"fmt"
	"log"

	. "github.com/VladislavTroyan163/POST-request-generation-script/createrequest"
)

func main() {
	лицевойСчет := 924126

	request, err := CreateRequest(лицевойСчет)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(request)
}

package createRegistryOfIDs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	Tasks []ResponseTask `json:"tasks"`
}

type ResponseTask struct {
	Result ResponseResult `json:"result"`
}

type ResponseResult struct {
	Записи [][]any `json:"записи"`
}

func ExecuteRequestcreateRegistryOfIDs(request string, лицевойСчет int, config Config, token string) (resultRequest string, resultRow_id int, err error) {
	serverURL := ""
	switch config.ConfigApp.Mode {
	case "FL":
		serverURL = config.ConfigApp.ServerURLFL
	case "UL":
		serverURL = config.ConfigApp.ServerURLUL
	}

	body := strings.NewReader(request)
	req, err := http.NewRequest(
		http.MethodPost,
		serverURL,
		body,
	)
	if err != nil {
		return "", 0, fmt.Errorf("ExecuteRequestcreateRegistryOfIDs: ошибка: \n%w", err)
	}

	req.Header.Set("s-access-token", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("ExecuteRequestcreateRegistryOfIDs: ошибка: \n%w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("ExecuteRequestcreateRegistryOfIDs: ошибка: \n%w", err)
	}

	var response Response
	if json.Unmarshal(responseBody, &response) != nil {
		return "", 0, fmt.Errorf("ExecuteRequestcreateRegistryOfIDs: ошибка: \n%w", err)
	}
	if len(response.Tasks) == 0 ||
		len(response.Tasks[0].Result.Записи) == 0 ||
		len(response.Tasks[0].Result.Записи[0]) == 0 {
		return "", 0, fmt.Errorf("ExecuteRequestcreateRegistryOfIDs: запись не найдена")
	}
	number, ok := response.Tasks[0].Result.Записи[0][0].(float64)
	if !ok {
		return "", 0, fmt.Errorf("ExecuteRequestcreateRegistryOfIDs: $номерЗаписи имеет неожиданный тип")
	}

	resultRow_id = int(number)
	resultRequest = fmt.Sprintln(
		"\nHTTP-запрос завершился успешно:", resp.Status,
		"\nrow_id лицевого счета:", resultRow_id,
	)

	return resultRequest, resultRow_id, nil
}

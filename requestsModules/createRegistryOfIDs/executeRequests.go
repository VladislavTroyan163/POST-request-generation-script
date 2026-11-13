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
	Result int `json:"result"`
}

func ExecuteRequestNumRegistryOfIDs(request string, config Config, token string) (resultRequest string, resultRow_id int, err error) {
	serverURL := ""
	switch config.ConfigApp.Mode {
	case "FL":
		serverURL = config.ConfigApp.ServerURLFL
	case "UL":
		serverURL = config.ConfigApp.ServerURLUL
	case "ASD":
		serverURL = config.ConfigApp.ServerURLASD
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
	if len(response.Tasks) == 0 {
		return "", 0, fmt.Errorf("ExecuteRequestcreateRegistryOfIDs: запись не найдена")
	}

	resultRow_id = response.Tasks[0].Result
	resultRequest = fmt.Sprintln(
		"\nHTTP-запрос завершился успешно:", resp.Status,
		"\nrow_id:", resultRow_id,
	)

	return resultRequest, resultRow_id, nil
}
